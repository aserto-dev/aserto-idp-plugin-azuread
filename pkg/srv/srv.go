package srv

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/aserto-dev/aserto-idp-plugin-azuread/pkg/azureclient"
	"github.com/aserto-dev/aserto-idp-plugin-azuread/pkg/config"
	"github.com/aserto-dev/aserto-idp-plugin-azuread/pkg/transform"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/aserto-dev/idp-plugin-sdk/plugin"
)

type AzureADPlugin struct {
	Config       *config.AzureADConfig
	azureClient  *azureclient.AzureADClient
	page         int
	finishedRead bool
	op           plugin.OperationType
}

func NewAzureADPlugin() *AzureADPlugin {
	return &AzureADPlugin{
		Config: &config.AzureADConfig{},
	}
}

func (a *AzureADPlugin) GetConfig() plugin.Config {
	return &config.AzureADConfig{}
}

func (a *AzureADPlugin) GetVersion() (string, string, string) {
	return config.GetVersion()
}

func (a *AzureADPlugin) Open(cfg plugin.Config, operation plugin.OperationType) error {
	azureadConfig, ok := cfg.(*config.AzureADConfig)
	if !ok {
		return errors.New("invalid config")
	}

	a.Config = azureadConfig
	a.page = 0
	a.finishedRead = false
	a.op = operation

	var err error
	if azureadConfig.RefreshToken != "" {
		a.azureClient, err = azureclient.NewAzureADClientWithRefreshToken(
			context.Background(),
			azureadConfig.Tenant,
			azureadConfig.ClientID,
			azureadConfig.ClientSecret,
			azureadConfig.RefreshToken)
		return err
	}

	a.azureClient, err = azureclient.NewAzureADClient(
		context.Background(),
		azureadConfig.Tenant,
		azureadConfig.ClientID,
		azureadConfig.ClientSecret)
	return err
}

func (a *AzureADPlugin) Read() ([]*api.User, error) {
	if a.finishedRead {
		return nil, io.EOF
	}

	var errs error
	var users []*api.User

	if a.Config.UserPID != "" {
		user, err := a.readByPID(a.Config.UserPID)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
		return users, nil
	}

	if a.Config.UserEmail != "" {
		return a.readByEmail(a.Config.UserEmail)
	}

	aadUsers, err := a.azureClient.ListUsers()
	if err != nil {
		return nil, err
	}

	for _, user := range aadUsers.GetValue() {
		u := transform.Transform(user)
		users = append(users, u)
	}

	a.finishedRead = true

	return users, errs
}

func (a *AzureADPlugin) readByPID(id string) (*api.User, error) {

	aadUsers, err := a.azureClient.GetUserByID(id)
	a.finishedRead = true
	if err != nil {
		return nil, err
	}

	users := aadUsers.GetValue()
	if len(users) == 0 {
		return nil, fmt.Errorf("failed to get user by pid %s", id)
	}
	return transform.Transform(users[0]), nil
}

func (a *AzureADPlugin) readByEmail(email string) ([]*api.User, error) {
	var users []*api.User

	aadUsers, err := a.azureClient.GetUserByEmail(email)
	a.finishedRead = true
	if err != nil {
		return nil, err
	}

	azureadUsers := aadUsers.GetValue()
	if len(azureadUsers) < 1 {
		return nil, fmt.Errorf("failed to get user by email %s", email)
	}

	for _, user := range azureadUsers {
		apiUser := transform.Transform(user)
		users = append(users, apiUser)
	}

	return users, nil
}

func (a *AzureADPlugin) Write(user *api.User) error {
	return nil
}

func (a *AzureADPlugin) Delete(userID string) error {
	return nil
}

func (a *AzureADPlugin) Close() (*plugin.Stats, error) {
	return nil, nil
}
