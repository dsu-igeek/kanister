# vSphere snapshot commands

This is an application intended to copy vSphere snapshot data to an object store or from an object store to a vSphere disk.

## Usage

### push
- **description:** Writes data from a vSphere snapshot to an object store.
- **inputs:**
  - (-i) snapshotID
  - (-R) [profile](https://docs.kanister.io/architecture.html#profiles) that describes a remote object store
  - (-C) vsphere credentials
  - (-p) an optional path within an object store

- **example usage:**
```bash
LD_LIBRARY_PATH=/opt/vddk/lib64 bin/amd64/vsnap_copy push -i ivd:asdfaf:adfaf -R '{"apiVersion":"cr.kanister.io/v1alpha1","credential":{"secret":{"apiVersion":"v1","group":"","kind":"Secret","name":"XXXX","namespace":"kasten-io","resource":""},"type":"secret"},"kind":"Profile","location":{"bucket":"XXXX","endpoint":"","prefix":"","region":"us-west-1","type":"s3Compliant"},"skipSSLVerify":false}' -C '{ "vchost":"host", "vcuser":"user", "vcpass":"pass", "s3urlbase": "something"}'
```

### pull
Future