#!/bin/sh -x

# Copyright 2020 The Kanister Authors.
# 
# Copyright 2016 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

if [ -z "${PKG}" ]; then
    echo "PKG must be set"
    exit 1
fi
if [ -z "${VERSION}" ]; then
    echo "VERSION must be set"
    exit 1
fi

# gcc may not be installed
which gcc >/dev/null 2>&1
if [ $? -ne 0 ]; then
    echo "gcc not found"
    exit 1
fi

export GOPATH=/go
export GO111MODULE=on

# clone the astrolabe package if it has not been copied in yet
# this will persist out of the build container
VENDOR=${PWD}/vendor
VENDOR_ASTROLABE="${VENDOR}/github.com/vmware-tanzu/astrolabe"
VENDOR_GVDDK="${VENDOR_ASTROLABE}/vendor/github.com/vmware/gvddk"
mkdir -p $(dirname ${VENDOR_ASTROLABE})
if [ ! -d ${VENDOR_ASTROLABE} ]; then
    (cd $(dirname ${VENDOR_ASTROLABE}); git clone http://github.com/vmware-tanzu/astrolabe)
    # set up the local use of astrolabe (referenced from ALT_MOD below)
    (cd ${VENDOR_ASTROLABE}; if [ ! -f go.mod ] ; then go mod init ; fi)
    # set up the local use of gvddk (referenced from ALT_MOD belpw)
    (cd ${VENDOR_GVDDK}; if [ ! -f go.mod ] ; then go mod init ; fi)
fi

# set up an alternate go.mod file (needs go v1.14+)
# the file will get updated by go
ALT_MOD=./go_vsnap.mod
ALT_SUM=./go_vsnap.sum
cp go.mod ${ALT_MOD}
cat <<EOF >>${ALT_MOD}
require github.com/vmware-tanzu/astrolabe v0.0.0-00010101000000-000000000000
require github.com/vmware/gvddk v0.0.0-00010101000000-000000000000

replace github.com/vmware-tanzu/astrolabe => ${VENDOR_ASTROLABE}
replace github.com/vmware/gvddk => ${VENDOR_GVDDK}
EOF
cp go.sum ${ALT_SUM}

export CGO_ENABLED=1
export GO_EXTLINK_ENABLED=1
export CGO_LDFLAGS="-L/opt/vddk/lib64 -lvixDiskLib"
go install -v -modfile ${ALT_MOD} \
    -installsuffix "static" \
    -ldflags "-X ${PKG}/pkg/version.VERSION=${VERSION}" \
    ./cgo_cmd/vsnap_copy

# To run with cgo_shell:
#  LD_LIBRARY_PATH=/opt/vddk/lib64 bin/amd64/vsnap_copy
# To view the DLLs linked in /opt/vddk:
#  LD_LIBRARY_PATH=/opt/vddk/lib64/ ldd bin/amd64/vsnap_copy