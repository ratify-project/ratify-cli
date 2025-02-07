/*
Copyright The Ratify Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package verifier

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/notaryproject/notation-go/dir"
	"github.com/notaryproject/notation-go/verifier/trustpolicy"
	"github.com/notaryproject/notation-go/verifier/truststore"
	"github.com/ratify-project/ratify-verifier-go/notation"
)

type notationVerifierOptions struct {
	Name           string                `json:"name"`
	TrustPolicyDoc *trustpolicy.Document `json:"trustPolicyDoc"`
	TrustStorePath string                `json:"trustStorePath"`
}

func ParseNotationVerifierOptions(data []byte) (*notation.VerifierOptions, error) {
	opts := &notationVerifierOptions{}
	if err := json.Unmarshal(data, opts); err != nil {
		return nil, err
	}
	fileInfo, err := os.Stat(opts.TrustStorePath)
	if err != nil {
		return nil, fmt.Errorf("Error reading trust store: %s, err: %w", opts.TrustStorePath, err)
	}
	if !fileInfo.IsDir() {
		return nil, fmt.Errorf("Trust store is not a directory: %s", opts.TrustStorePath)
	}

	return &notation.VerifierOptions{
		Name:           opts.Name,
		TrustPolicyDoc: opts.TrustPolicyDoc,
		TrustStore:     truststore.NewX509TrustStore(dir.NewSysFS(opts.TrustStorePath)),
	}, nil
}
