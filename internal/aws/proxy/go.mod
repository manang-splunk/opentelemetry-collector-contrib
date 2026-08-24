module github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/proxy

go 1.17

require (
	github.com/aws/aws-sdk-go v1.42.14
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal v0.40.0
	github.com/stretchr/testify v1.12.1
	go.opentelemetry.io/collector v0.40.0
	go.uber.org/zap v1.19.1
)

require (
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
)

replace github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal => ../../../internal/coreinternal
