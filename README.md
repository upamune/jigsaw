## jigsaw

### Usage

```bash
$ go run . trace.json -c config.yaml
```

```bash
$ cat config.yaml
include_services:
  - foo-service
  - bar-service
exclude_grpc_services:
  - /foo.bar.v0.Service
grpc_serivce_alias:
  /foo.bar.v1.Service: v1-serivce
  /foo.bar.v2.Service: v2-serivce
```

