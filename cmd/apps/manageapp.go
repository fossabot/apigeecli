// Copyright 2020 Google LLC
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

package apps

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/client/apps"
)

//ManageCmd to create developer keys
var ManageCmd = &cobra.Command{
	Use:   "manage",
	Short: "Approve or revoke a developer app",
	Long:  "Approve or revoke a developer app",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = apps.Manage(name, action)
		return
	},
}

func init() {

	ManageCmd.Flags().StringVarP(&name, "name", "n",
		"", "Developer app name")
	ManageCmd.Flags().StringVarP(&action, "action", "x",
		"revoke", "Action to perform - revoke or approve")

	_ = ManageCmd.MarkFlagRequired("name")
	_ = ManageCmd.MarkFlagRequired("action")
}