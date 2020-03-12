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
	"encoding/json"
	"fmt"
	"io"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vmware-tanzu/astrolabe/pkg/astrolabe"
	"github.com/vmware-tanzu/astrolabe/pkg/ivd"
)

type SnapshotManager struct {
	config  *VSphereCreds
	ivdPETM *ivd.IVDProtectedEntityTypeManager
}

func newSnapshotManager(config *VSphereCreds) (*SnapshotManager, error) {
	params := map[string]interface{}{
		"vcHost":     config.Host,
		"vcUser":     config.User,
		"vcPassword": config.Pass,
	}
	ivdPETM, err := ivd.NewIVDProtectedEntityTypeManagerFromConfig(params, config.S3UrlBase, logrus.New())
	if err != nil {
		return nil, fmt.Errorf("Unable to create ivd Protected Entity Manager from config %s", err.Error())
	}
	return &SnapshotManager{
		config:  config,
		ivdPETM: ivdPETM,
	}, nil
}

type VSphereCreds struct {
	Host      string `json:"vchost"`
	User      string `json:"vcuser"`
	Pass      string `json:"vcpass"`
	S3UrlBase string `json:"s3urlbase"`
}

func (v *VSphereCreds) Unmarshal(creds []byte) error {
	return json.Unmarshal(creds, v)
}

func (v *VSphereCreds) Validate() error {
	if v.Host == "" {
		return fmt.Errorf("missing endpoint value")
	}
	if v.User == "" {
		return fmt.Errorf("missing username value")
	}
	if v.Pass == "" {
		return fmt.Errorf("missing password value")
	}
	if v.S3UrlBase == "" {
		return fmt.Errorf("missing s3URLBase value")
	}
	return nil
}

func GetVsphereCreds(cmd *cobra.Command) (*VSphereCreds, error) {
	creds := &VSphereCreds{}
	if err := creds.Unmarshal([]byte(cmd.Flag(vSphereCreds).Value.String())); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal vsphere credentials")
	}
	err := creds.Validate()
	return creds, errors.Wrap(err, "failed to validate vSphere credentials")
}

func GetDataReaderFromSnapshot(ctx context.Context, config *VSphereCreds, snapshotID string) (io.ReadCloser, error) {
	snapManager, err := newSnapshotManager(config)
	if err != nil {
		return nil, err
	}
	// expecting a snapshot id of the form type:volumeID:snapshotID
	peID, err := astrolabe.NewProtectedEntityIDFromString(snapshotID)
	if err != nil {
		return nil, err
	}
	pe, err := snapManager.ivdPETM.GetProtectedEntity(ctx, peID)
	if err != nil {
		return nil, err
	}

	return pe.GetDataReader(ctx)
}
