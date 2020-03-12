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
	"github.com/vmware-tanzu/astrolabe/pkg/astrolabe"
	"github.com/vmware-tanzu/astrolabe/pkg/ivd"
)

type SnapshotManager struct {
	config  *VSphereCreds
	ivdPETM *ivd.IVDProtectedEntityTypeManager
}

func newSnapshotManager(creds *VSphereCreds) (*SnapshotManager, error) {
	params := map[string]interface{}{
		"vcHost":     creds.Host,
		"vcUser":     creds.User,
		"vcPassword": creds.Pass,
	}
	ivdPETM, err := ivd.NewIVDProtectedEntityTypeManagerFromConfig(params, creds.S3UrlBase, logrus.New())
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create ivd Protected Entity Manager from config")
	}
	return &SnapshotManager{
		config:  creds,
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
	if v.Host == "" || v.User == "" || v.Pass == "" || v.S3UrlBase == "" {
		return fmt.Errorf("invalid vSphere credentials")
	}
	return nil
}

func GetDataReaderFromSnapshot(ctx context.Context, creds *VSphereCreds, snapshotID string) (io.ReadCloser, error) {
	snapManager, err := newSnapshotManager(creds)
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
