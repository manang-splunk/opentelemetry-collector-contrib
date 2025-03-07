// Copyright 2020, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package k8sobserver // import "github.com/open-telemetry/opentelemetry-collector-contrib/extension/observer/k8sobserver"

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/extension/extensionhelper"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/k8sconfig"
)

const (
	// The value of extension "type" in configuration.
	typeStr config.Type = "k8s_observer"
)

// NewFactory should be called to create a factory with default values.
func NewFactory() component.ExtensionFactory {
	return extensionhelper.NewFactory(
		typeStr,
		createDefaultConfig,
		createExtension)
}

// CreateDefaultConfig creates the default configuration for the extension.
func createDefaultConfig() config.Extension {
	return &Config{
		ExtensionSettings: config.NewExtensionSettings(config.NewComponentID(typeStr)),
		APIConfig:         k8sconfig.APIConfig{AuthType: k8sconfig.AuthTypeServiceAccount},
	}
}

// CreateExtension creates the extension based on this config.
func createExtension(
	ctx context.Context,
	params component.ExtensionCreateSettings,
	cfg config.Extension,
) (component.Extension, error) {
	oCfg := cfg.(*Config)

	clientset, err := k8sconfig.MakeClient(oCfg.APIConfig)
	if err != nil {
		return nil, err
	}

	listWatch := cache.NewListWatchFromClient(
		clientset.CoreV1().RESTClient(), "pods", v1.NamespaceAll,
		fields.OneTermEqualSelector("spec.nodeName", oCfg.Node))

	return newObserver(params.Logger, oCfg, listWatch)
}
