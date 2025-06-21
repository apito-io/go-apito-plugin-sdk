module github.com/apito-io/go-apito-plugin-sdk

go 1.22.5

toolchain go1.23.3

require (
	github.com/hashicorp/go-plugin v1.6.0
	gitlab.com/apito.io/buffers v0.0.0
	google.golang.org/grpc v1.65.0
	google.golang.org/protobuf v1.34.1
)

require (
	github.com/fatih/color v1.13.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/hashicorp/go-hclog v1.5.0 // indirect
	github.com/hashicorp/yamux v0.1.1 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mitchellh/go-testing-interface v0.0.0-20171004221916-a61a99592b77 // indirect
	github.com/oklog/run v1.0.0 // indirect
	golang.org/x/net v0.27.0 // indirect
	golang.org/x/sys v0.23.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240604185151-ef581f913117 // indirect
)

replace gitlab.com/apito.io/buffers => ../buffers
