## jigsaw

Automatically generate a sequence diagram from JSON of Trace in Datadog.

⚠️ Only gRPC calls appear in the sequence diagram.

### Example

#### w/ response

```bash
$ jigsaw -config ./example/config.yaml ./example/trace.json
@startuml
"v1-service" -> "v2-service": Ping Request
"v2-service" -> "v3-service": Pong Request
"v2-service" <-- "v3-service": Pong Response
"v1-service" <-- "v2-service": Ping Response
@enduml
```

![output](https://user-images.githubusercontent.com/8219560/127803698-763b9343-5429-417a-89b9-492e88ed08ff.png)


#### w/o response

```bash
$ jigsaw -config ./example/config.yaml -no-response ./example/trace.json
@startuml
"v1-service" -> "v2-service": Ping Request
"v2-service" -> "v3-service": Pong Request
@enduml
```

![output](https://user-images.githubusercontent.com/8219560/127775036-b13113ff-496c-489c-8b1d-a6a756c62d97.png)

### Usage

You can get a trace as a JSON via `https://app.datadoghq.com/api/v1/trace/TRACE_ID`.

```bash
$ go get -u github.com/upamune/jigsaw
$ jigsaw trace.json -c config.yaml
```

```bash
$ cat config.yaml
type: plantuml
include_services:
  - foo-service
  - bar-service
exclude_grpc_services:
  - /foo.bar.v0.Service
grpc_serivce_alias:
  /foo.bar.v1.Service: v1-serivce
  /foo.bar.v2.Service: v2-serivce
```
