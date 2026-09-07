module github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/proxy

go 1.26.0

require (
	github.com/aws/aws-sdk-go v1.42.14
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal v0.40.0
	github.com/stretchr/testify v1.12.1
	go.opentelemetry.io/collector/config/confignet v1.66.0
	go.opentelemetry.io/collector/config/configtls v1.66.0
	go.uber.org/zap v1.28.0
)

require (
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/foxboron/go-tpm-keyfiles v0.0.0-20250903184740-5d135037bd4d // indirect
	github.com/go-viper/mapstructure/v2 v2.5.0 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/google/go-tpm v0.9.8 // indirect
	github.com/hashicorp/go-version v1.9.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/knadh/koanf v1.5.0 // indirect
	github.com/knadh/koanf/v2 v2.3.6 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	go.opentelemetry.io/collector/config/configopaque v1.66.0 // indirect
	go.opentelemetry.io/collector/confmap v1.66.0 // indirect
	go.opentelemetry.io/collector/featuregate v1.66.0 // indirect
	go.opentelemetry.io/collector/internal/testutil v0.160.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/crypto v0.52.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
)

replace github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal => ../../../internal/coreinternal
