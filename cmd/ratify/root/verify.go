// Copyright The Ratify Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package root

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/ratify-project/ratify-cli/v2/internal/verifier"
	"github.com/ratify-project/ratify-go"
	"github.com/ratify-project/ratify-verifier-go/notation"
	"github.com/spf13/cobra"
	"oras.land/oras-go/v2/registry"
)

type verifyOptions struct {
	subject        string
	configFilePath string
	storePath      string
}

type validationResult struct {
	Succeeded       bool               `json:"succeeded"`
	ArtifactReports []validationReport `json:"artifactReports"`
}

type validationReport struct {
	Subject         string               `json:"subject"`
	Artifact        string               `json:"artifact"`
	ArtifactType    string               `json:"artifactType"`
	Results         []verificationResult `json:"results"`
	ArtifactReports []validationReport   `json:"artifactReports"`
}

type verificationResult struct {
	Succeeded    bool   `json:"succeeded"`
	Description  string `json:"description,omitempty"`
	VerifierName string `json:"verifierName"`
	VerifierType string `json:"verifierType"`
	Detail       any    `json:"detail,omitempty"`
}

func verifyCommand(opts *verifyOptions) *cobra.Command {
	if opts == nil {
		opts = &verifyOptions{}
	}
	longMessage := `Verify an artifact
Prerequisite: added a trust store for notation verifier and an OCI store saving artifacts.

Example - Verify an artifact:
  ratify verify --subject <subject> --config <config file> --store <store path>
`

	cmd := &cobra.Command{
		Use:   "verify",
		Short: "Verify the provided artifact",
		Long:  longMessage,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runVerify(cmd, opts)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.subject, "subject", "s", "", "subject to verify")
	flags.StringVarP(&opts.configFilePath, "config", "c", "", "path to the verifier config file")
	flags.StringVarP(&opts.storePath, "store", "t", "", "path to the store")
	return cmd
}

func runVerify(cmd *cobra.Command, opts *verifyOptions) error {
	if opts.subject == "" {
		return errors.New("subject is required")
	}
	ref, err := registry.ParseReference(opts.subject)
	if err != nil {
		return fmt.Errorf("Invalid subject reference: %s, err: %w", opts.subject, err)
	}
	repo := ref.Registry + "/" + ref.Repository

	store, err := createStore(cmd.Context(), opts.storePath, repo)
	if err != nil {
		return err
	}

	verifiers, err := loadNotationVerifiersFromConfigFile(opts.configFilePath)
	if err != nil {
		return fmt.Errorf("Error loading verifiers: %v", err)
	}

	executor, err := ratify.NewExecutor(store, verifiers, nil)
	if err != nil {
		return fmt.Errorf("Error creating executor: %v", err)
	}

	validateOpts := ratify.ValidateArtifactOptions{
		Subject: opts.subject,
	}
	result, err := executor.ValidateArtifact(cmd.Context(), validateOpts)
	if err != nil {
		return fmt.Errorf("Error validating artifact: %v", err)
	}

	printReport(result, repo)
	return nil
}

func createStore(ctx context.Context, storePath, repo string) (*ratify.StoreMux, error) {
	fileInfo, err := os.Stat(storePath)
	if os.IsNotExist(err) {
		return nil, errors.New("store path does not exist")
	}
	if err != nil {
		return nil, fmt.Errorf("Error checking store path: %v", err)
	}

	var store *ratify.OCIStore
	if fileInfo.IsDir() {
		fs := os.DirFS(storePath)
		store, err = ratify.NewOCIStoreFromFS(ctx, "store", fs)
		if err != nil {
			return nil, fmt.Errorf("Error creating store from FS: %v", err)
		}
	} else {
		store, err = ratify.NewOCIStoreFromTar(ctx, "store", storePath)
		if err != nil {
			return nil, fmt.Errorf("Error creating store from tar: %v", err)
		}
	}

	mux := ratify.NewStoreMux("multiplexer")
	if err := mux.Register(repo, store); err != nil {
		return nil, fmt.Errorf("Error registering store: %v", err)
	}
	return mux, nil
}

func loadNotationVerifiersFromConfigFile(configFilePath string) ([]ratify.Verifier, error) {
	body, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("Error reading config file: %w", err)
	}

	verifierOpts, err := verifier.ParseNotationVerifierOptions(body)
	if err != nil {
		return nil, fmt.Errorf("Error loading verifier options: %w", err)
	}

	verifier, err := notation.NewVerifier(verifierOpts)
	if err != nil {
		return nil, fmt.Errorf("Error creating verifier: %w", err)
	}

	return []ratify.Verifier{verifier}, nil
}

func printReport(result *ratify.ValidationResult, repo string) {
	convertedResult := convertValidationResult(result, repo)
	jsonData, err := json.MarshalIndent(convertedResult, "", "    ")
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}
	fmt.Println(string(jsonData))
}

func convertValidationResult(result *ratify.ValidationResult, repo string) *validationResult {
	return &validationResult{
		Succeeded:       result.Succeeded,
		ArtifactReports: convertValidationReports(result.ArtifactReports, repo),
	}
}

func convertValidationReports(reports []*ratify.ValidationReport, repo string) []validationReport {
	var convertedReports []validationReport
	for _, report := range reports {
		convertedReport := validationReport{
			Subject:         report.Subject,
			Artifact:        repo + "@" + report.Artifact.Digest.String(),
			ArtifactType:    report.Artifact.ArtifactType,
			Results:         convertVerificationResults(report.Results),
			ArtifactReports: convertValidationReports(report.ArtifactReports, repo),
		}
		convertedReports = append(convertedReports, convertedReport)
	}
	return convertedReports
}

func convertVerificationResults(results []*ratify.VerificationResult) []verificationResult {
	var convertedResults []verificationResult
	for _, result := range results {
		detail := result.Detail
		if result.Err != nil {
			detail = result.Err.Error()
		}
		convertedResult := verificationResult{
			Succeeded:    result.Err == nil,
			Description:  result.Description,
			VerifierName: result.Verifier.Name(),
			VerifierType: result.Verifier.Type(),
			Detail:       detail,
		}
		convertedResults = append(convertedResults, convertedResult)
	}
	return convertedResults
}
