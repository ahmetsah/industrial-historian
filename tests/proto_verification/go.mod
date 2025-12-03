module verify

go 1.23

toolchain go1.24.11

require (
	github.com/ahmet/historian/go-services/pkg/proto v0.0.0
	google.golang.org/protobuf v1.36.10
)

replace github.com/ahmet/historian/go-services/pkg/proto => ../../go-services/pkg/proto
