// Copyright 2020 Coinbase, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/coinbase/rosetta-sdk-go/fetcher"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/spf13/cobra"
)

const (
	// fileMode 0600 indicates that the user/owner can read and write
	// but can't execute.
	fileMode = 0600
)

var (
	createConfigurationCmd = &cobra.Command{
		Use:   "create:configuration",
		Short: "Generate a static configuration file for the Asserter",
		Long: `In production deployments, it is useful to initialize the response
Asserter (https://github.com/coinbase/rosetta-sdk-go/tree/master/asserter) using
a static configuration instead of intializing a configuration dynamically
from the node. This allows a client to error on new types/statuses that may
have been added in an update instead of silently erroring.

To use this command, simply provide an absolute path as the argument for where
the configuration file should be saved (in JSON). Populate the optional
--server-url flag with the url of the server to generate the configuration
from.`,
		Run:  runCreateConfigurationCmd,
		Args: cobra.ExactArgs(1),
	}
)

func runCreateConfigurationCmd(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	// Create a new fetcher
	newFetcher := fetcher.New(
		ServerURL,
	)

	// Initialize the fetcher's asserter
	_, _, err := newFetcher.InitializeAsserter(ctx)
	if err != nil {
		log.Fatal(err)
	}

	configuration, err := newFetcher.Asserter.ClientConfiguration()
	if err != nil {
		log.Fatal(fmt.Errorf("%w: unable to generate spec", err))
	}

	specString := types.PrettyPrintStruct(configuration)
	log.Printf("Spec File: %s\n", specString)

	if err := ioutil.WriteFile(path.Clean(args[0]), []byte(specString), os.FileMode(fileMode)); err != nil {
		log.Fatal(err)
	}
}
