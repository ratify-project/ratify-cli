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
	"io"
	"os"
	"strings"
	"testing"

	"github.com/ratify-project/ratify-cli/v2/internal/version"
)

func TestVersionCommand(t *testing.T) {
	// set git commit
	version.GitCommit = "testCommit"
	cmd := versionCommand()

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Unexpected error executing version command: %v", err)
	}

	w.Close()
	os.Stdout = oldStdout
	out, _ := io.ReadAll(r)
	output := string(out)
	if !strings.Contains(output, "Version:") {
		t.Errorf("Expected output to contain 'Version:', got: %s", output)
	}
	if !strings.Contains(output, "Go version:") {
		t.Errorf("Expected output to contain 'Go version:', got: %s", output)
	}
	if !strings.Contains(output, "Git commit:") {
		t.Errorf("Expected output to contain 'Git commit:', got: %s", output)
	}
}
