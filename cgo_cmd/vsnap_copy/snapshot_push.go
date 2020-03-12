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
	"context"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/kanisterio/kanister/pkg/location"
	"github.com/kanisterio/kanister/pkg/param"
)

const (
	snapIDFlagName = "snap"
)

func newSnapshotPushCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "push",
		Short: "Push a source file or stdin stream to s3-compliant object storage",
		Args:  cobra.ExactArgs(0),
		// TODO: Example invocations
		RunE: func(c *cobra.Command, args []string) error {
			return runSnapshotPush(c, args)
		},
	}
	cmd.PersistentFlags().StringP(snapIDFlagName, "i", "", "Specify the snapshot ID (required)")
	_ = cmd.MarkPersistentFlagRequired(snapIDFlagName)
	return cmd
}

func runSnapshotPush(cmd *cobra.Command, args []string) error {
	snapshotID := cmd.Flag(snapIDFlagName).Value.String()
	profile, err := unmarshalProfileFlag(cmd)
	if err != nil {
		return err
	}
	path := pathFlag(cmd)
	config, err := GetVsphereCreds(cmd)
	if err != nil {
		return err
	}
	ctx := context.Background()
	return copySnapshotToObjectStore(ctx, config, profile, snapshotID, path)
}

func copySnapshotToObjectStore(ctx context.Context, config *VSphereCreds, profile *param.Profile, snapshot string, path string) error {
	reader, err := GetDataReaderFromSnapshot(ctx, config, snapshot)
	if err != nil {
		return errors.Wrap(err, "Failed to get reader from snapshot")
	}
	return location.Write(ctx, reader, *profile, path)
}
