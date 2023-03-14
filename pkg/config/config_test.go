package config_test

import (
	"testing"

	"github.com/aserto-dev/aserto-idp-plugin-azuread/pkg/config"
	"github.com/aserto-dev/idp-plugin-sdk/plugin"
	"github.com/stretchr/testify/require"
)

func TestValidateWithEmptyTenant(t *testing.T) {
	assert := require.New(t)
	cfg := config.AzureADConfig{
		Tenant:       "",
		ClientID:     "id",
		ClientSecret: "secret",
	}

	err := cfg.Validate(plugin.OperationTypeRead)

	assert.NotNil(err)
	assert.Equal("rpc error: code = InvalidArgument desc = no tenant was provided", err.Error())
}

func TestValidateWithEmptyClientID(t *testing.T) {
	assert := require.New(t)
	cfg := config.AzureADConfig{
		Tenant:       "tenant",
		ClientID:     "",
		ClientSecret: "secret",
	}

	err := cfg.Validate(plugin.OperationTypeRead)

	assert.NotNil(err)
	assert.Equal("rpc error: code = InvalidArgument desc = no client id was provided", err.Error())
}

func TestValidateWithEmptyClientSecret(t *testing.T) {
	assert := require.New(t)
	cfg := config.AzureADConfig{
		Tenant:       "tenant",
		ClientID:     "id",
		ClientSecret: "",
	}

	err := cfg.Validate(plugin.OperationTypeRead)

	assert.NotNil(err)
	assert.Equal("rpc error: code = InvalidArgument desc = no client secret was provided", err.Error())
}

func TestValidateWithInvalidCredentials(t *testing.T) {
	assert := require.New(t)
	cfg := config.AzureADConfig{
		Tenant:       "tenant",
		ClientID:     "id",
		ClientSecret: "secret",
	}

	err := cfg.Validate(plugin.OperationTypeWrite)

	assert.NotNil(err)
	assert.Contains(err.Error(), "Internal desc = failed to retrieve users from AzureAD")
}

func TestValidateWithUserIDAndEmail(t *testing.T) {
	assert := require.New(t)
	cfg := config.AzureADConfig{
		Tenant:       "tenant",
		ClientID:     "id",
		ClientSecret: "secret",
		UserPID:      "someID",
		UserEmail:    "test@email.com",
	}

	err := cfg.Validate(plugin.OperationTypeWrite)

	assert.NotNil(err)
	assert.Contains(err.Error(), "rpc error: code = InvalidArgument desc = an user PID and an user email were provided; please specify only one")
}

func TestDescription(t *testing.T) {
	assert := require.New(t)
	cfg := config.AzureADConfig{
		Tenant:       "tenant",
		ClientID:     "id",
		ClientSecret: "secret",
	}

	description := cfg.Description()

	assert.Equal("AzureAD plugin", description)
}
