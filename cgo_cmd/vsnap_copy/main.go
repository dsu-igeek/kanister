// Copyright 2020 The Kanister Authors.
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

package main

import (
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/kanisterio/kanister/pkg/log"
	"github.com/kanisterio/kanister/pkg/param"
)

func main() {
	Execute()
}

const (
	pathFlagName    = "path"
	profileFlagName = "profile"
	vSphereCreds    = "vcreds"
)

func Execute() {
	root := newRootCommand()
	if err := root.Execute(); err != nil {
		log.WithError(err).Print("vsnapcopy failed to execute")
		os.Exit(1)
	}
}

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vsnapcopy",
		Short: "push, pull from object storage",
	}
	cmd.AddCommand(newSnapshotPushCommand())
	//cmd.AddCommand(newSnapshotPullCommand())
	cmd.PersistentFlags().StringP(profileFlagName, "R", "", "Profile describing a remote store as a JSON string (required)")
	cmd.PersistentFlags().StringP(vSphereCreds, "C", "", "VSphereCredentials as a JSON string (required)")
	cmd.PersistentFlags().StringP(pathFlagName, "p", "", "Specify a path within the object store (optional)")
	_ = cmd.MarkPersistentFlagRequired(profileFlagName)
	_ = cmd.MarkPersistentFlagRequired(vSphereCreds)
	return cmd
}

func ParsePathFlag(cmd *cobra.Command) string {
	return cmd.Flag(pathFlagName).Value.String()
}

func ParseRemoteStoreFlag(cmd *cobra.Command) (*param.Profile, error) {
	p := &param.Profile{}
	err := p.Unmarshal([]byte(cmd.Flag(profileFlagName).Value.String()))
	return p, errors.Wrap(err, "Failed to unmarshal profile")
}

func ParseVsphereCredFlag(cmd *cobra.Command) (*VSphereCreds, error) {
	creds := &VSphereCreds{}
	if err := creds.Unmarshal([]byte(cmd.Flag(vSphereCreds).Value.String())); err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal vSphere credentials")
	}
	err := creds.Validate()
	return creds, errors.Wrap(err, "Failed to validate vSphere credentials")
}
