package srv

import (
	"io"
	"testing"

	"github.com/aserto-dev/aserto-idp-plugin-azuread/pkg/config"
	azureADTestUtils "github.com/aserto-dev/aserto-idp-plugin-azuread/pkg/testutils"
	"github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/aserto-dev/go-utils/testutil"
	"github.com/aserto-dev/idp-plugin-sdk/plugin"
	"github.com/stretchr/testify/require"
)

func CreateConfig() config.AzureADConfig {
	return config.AzureADConfig{
		Tenant:       testutil.VaultValue("azuread-idp-test-account.tenant"),
		ClientID:     testutil.VaultValue("azuread-idp-test-account.client-id"),
		ClientSecret: testutil.VaultValue("azuread-idp-test-account.client-secret"),
	}
}

func TestOpen(t *testing.T) {
	assert := require.New(t)

	cfg := CreateConfig()
	err := cfg.Validate(plugin.OperationTypeRead)
	assert.Nil(err)

	azureADPlugin := NewAzureADPlugin()
	err = azureADPlugin.Open(&cfg, plugin.OperationTypeRead)
	assert.Nil(err)

	stats, err := azureADPlugin.Close()
	assert.Nil(err)
	assert.Nil(stats)
}

func TestWrite(t *testing.T) {
	t.Skip()
	assert := require.New(t)

	apiUser := azureADTestUtils.CreateTestAPIUser("2ff319e101e1", "Test User", "user@test.com", "https://github.com/aserto-demo/contoso-ad-sample/raw/main/UserImages/Euan%20Garden.jpg")
	cfg := CreateConfig()
	err := cfg.Validate(plugin.OperationTypeWrite)
	assert.Nil(err)

	azureADPlugin := NewAzureADPlugin()
	err = azureADPlugin.Open(&cfg, plugin.OperationTypeWrite)
	assert.Nil(err)

	err = azureADPlugin.Write(apiUser)
	assert.Nil(err)

	stats, err := azureADPlugin.Close()
	assert.Nil(err)
	assert.NotNil(stats)
	assert.Equal(int32(1), stats.Received)
	assert.Equal(int32(1), stats.Created)
	assert.Equal(int32(0), stats.Errors)
}

func TestReadInvalidUserID(t *testing.T) {
	t.Skip()
	assert := require.New(t)

	cfg := CreateConfig()
	cfg.UserPID = "somerandomID"
	err := cfg.Validate(plugin.OperationTypeRead)
	assert.Nil(err)

	azureADPlugin := NewAzureADPlugin()
	err = azureADPlugin.Open(&cfg, plugin.OperationTypeRead)
	assert.Nil(err)

	users, err := azureADPlugin.Read()
	assert.NotNil(err)
	assert.Equal(0, len(users))

	_, err = azureADPlugin.Read()
	assert.NotNil(err)
	assert.Equal(io.EOF, err)

	stats, err := azureADPlugin.Close()
	assert.Nil(err)
	assert.Nil(stats)
}

func TestReadUserByID(t *testing.T) {
	t.Skip()
	assert := require.New(t)

	cfg := CreateConfig()
	cfg.UserPID = "2ff319e101e1"
	err := cfg.Validate(plugin.OperationTypeRead)
	assert.Nil(err)

	azureADPlugin := NewAzureADPlugin()
	err = azureADPlugin.Open(&cfg, plugin.OperationTypeRead)
	assert.Nil(err)

	users, err := azureADPlugin.Read()
	assert.Nil(err)
	assert.Equal(1, len(users))
	assert.Equal(users[0].GetId(), "2ff319e101e1")
	assert.Equal(users[0].GetDisplayName(), "Test User")

	_, err = azureADPlugin.Read()
	assert.NotNil(err)
	assert.Equal(io.EOF, err)

	stats, err := azureADPlugin.Close()
	assert.Nil(err)
	assert.Nil(stats)
}

func TestReadInvalidUserEmail(t *testing.T) {
	assert := require.New(t)

	cfg := CreateConfig()
	cfg.UserEmail = "invalidID"
	err := cfg.Validate(plugin.OperationTypeRead)
	assert.Nil(err)

	azureADPlugin := NewAzureADPlugin()
	err = azureADPlugin.Open(&cfg, plugin.OperationTypeRead)
	assert.Nil(err)

	users, err := azureADPlugin.Read()
	assert.NotNil(err)
	assert.Equal(0, len(users))

	_, err = azureADPlugin.Read()
	assert.NotNil(err)
	assert.Equal(io.EOF, err)

	stats, err := azureADPlugin.Close()
	assert.Nil(err)
	assert.Nil(stats)
}

func TestReadUserByEmail(t *testing.T) {
	assert := require.New(t)

	cfg := CreateConfig()
	cfg.UserEmail = "omri@aserto.com"
	err := cfg.Validate(plugin.OperationTypeRead)
	assert.Nil(err)

	azureADPlugin := NewAzureADPlugin()
	err = azureADPlugin.Open(&cfg, plugin.OperationTypeRead)
	assert.Nil(err)

	users, err := azureADPlugin.Read()
	assert.Nil(err)
	assert.Equal(1, len(users))
	assert.Equal(users[0].GetEmail(), "omri@aserto.com")
	assert.Equal(users[0].GetDisplayName(), "Omri Gazitt")

	_, err = azureADPlugin.Read()
	assert.NotNil(err)
	assert.Equal(io.EOF, err)

	stats, err := azureADPlugin.Close()
	assert.Nil(err)
	assert.Nil(stats)
}

func TestRead(t *testing.T) {
	assert := require.New(t)

	cfg := CreateConfig()
	err := cfg.Validate(plugin.OperationTypeRead)
	assert.Nil(err)

	azureADPlugin := NewAzureADPlugin()
	err = azureADPlugin.Open(&cfg, plugin.OperationTypeRead)
	assert.Nil(err)

	users, err := azureADPlugin.Read()
	assert.Nil(err)
	assert.Equal(12, len(users))

	_, err = azureADPlugin.Read()
	assert.NotNil(err)
	assert.Equal(io.EOF, err)

	stats, err := azureADPlugin.Close()
	assert.Nil(err)
	assert.Nil(stats)
}

func TestDelete(t *testing.T) {
	t.Skip()
	assert := require.New(t)

	cfg := CreateConfig()
	err := cfg.Validate(plugin.OperationTypeRead)
	assert.Nil(err)

	azureADPlugin := NewAzureADPlugin()
	err = azureADPlugin.Open(&cfg, plugin.OperationTypeRead)
	assert.Nil(err)

	users, err := azureADPlugin.Read()
	assert.Nil(err)
	assert.Less(1, len(users))

	var testUser *api.User

	for _, user := range users {
		if user.DisplayName == "Test User" {
			testUser = user
		}
	}
	assert.NotNil(testUser)

	stats, err := azureADPlugin.Close()
	assert.Nil(err)
	assert.Nil(stats)

	err = azureADPlugin.Open(&cfg, plugin.OperationTypeDelete)
	assert.Nil(err)

	err = azureADPlugin.Delete(testUser.Id)
	assert.Nil(err)

	stats, err = azureADPlugin.Close()
	assert.Nil(err)
	assert.Nil(stats)
}
