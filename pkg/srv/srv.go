package srv

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/aserto-dev/aserto-idp-plugin-azuread/pkg/azureclient"
	"github.com/aserto-dev/aserto-idp-plugin-azuread/pkg/config"
	"github.com/aserto-dev/aserto-idp-plugin-azuread/pkg/transform"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/aserto-dev/idp-plugin-sdk/plugin"
	"gopkg.in/auth0.v5/management"
)

type AzureADPlugin struct {
	Config       *config.AzureADConfig
	azureClient  *azureclient.AzureADClient
	page         int
	finishedRead bool
	totalSize    int64
	jobs         []management.Job
	users        []map[string]interface{}
	connectionID string
	wg           sync.WaitGroup
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

	azureClient, err := azureclient.NewAzureADClient(
		context.Background(),
		azureadConfig.Tenant,
		azureadConfig.ClientID,
		azureadConfig.ClientSecret)
	if err != nil {
		return err
	}

	a.azureClient = azureClient

	return nil
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
		fmt.Printf("User: %s\n", *user.GetDisplayName())
		fmt.Printf("  ID: %s\n", *user.GetId())

		noEmail := "NO EMAIL"
		email := user.GetMail()
		if email == nil {
			email = &noEmail
		}
		fmt.Printf("  Email: %s\n", *email)
		u := transform.Transform(user)
		users = append(users, u)
	}

	a.finishedRead = true

	return users, errs
}

func (a *AzureADPlugin) readByPID(id string) (*api.User, error) {

	aadUsers, err := a.azureClient.GetUser(id)
	if err != nil {
		return nil, err
	}

	for _, user := range aadUsers.GetValue() {
		if user == nil {
			return nil, fmt.Errorf("failed to get user by pid %s", id)
		}
		return transform.Transform(user), nil
	}

	return nil, fmt.Errorf("failed to get user by pid %s", id)
}

func (a *AzureADPlugin) readByEmail(email string) ([]*api.User, error) {
	var users []*api.User

	aadUsers, err := a.azureClient.GetUser(email)
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
