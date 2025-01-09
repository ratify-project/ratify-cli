// Copyright The Ratify Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/ratify-project/ratify-cli/v2/cmd/ratify/root"
	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "ratify",
	Short: "Ratify is a reference artifact tool for managing and verifying reference artifacts",
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	if err := root.New().ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
