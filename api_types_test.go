package linters

import (
	"testing"

	"github.com/golangci/plugin-module-register/register"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestPluginSDKv2(t *testing.T) {
	newPlugin, err := register.GetPlugin("linter-sdkv2")
	require.NoError(t, err)

	plugin, err := newPlugin(nil)
	require.NoError(t, err)

	analyzers, err := plugin.BuildAnalyzers()
	require.NoError(t, err)

	analysistest.Run(t, testdataDir(t), analyzers[0], "apiTypes/api/example")
	analysistest.Run(t, testdataDir(t), analyzers[1], "endpointStructFields")
}
