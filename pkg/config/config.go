package config

import (
	"context"

	"github.com/aserto-dev/aserto-idp-plugin-azuread/pkg/azureclient"
	"github.com/aserto-dev/idp-plugin-sdk/plugin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// values set by linker using ldflag -X.
var (
	ver    string // nolint:gochecknoglobals // set by linker
	date   string // nolint:gochecknoglobals // set by linker
	commit string // nolint:gochecknoglobals // set by linker
)

func GetVersion() (string, string, string) {
	return ver, date, commit
}

type AzureADConfig struct {
	Tenant       string `description:"AzureAD tenant" kind:"attribute" mode:"normal" readonly:"false" name:"tenant"`
	ClientID     string `description:"AzureAD Client ID" kind:"attribute" mode:"normal" readonly:"false" name:"client-id"`
	ClientSecret string `description:"AzureAD Client Secret" kind:"attribute" mode:"normal" readonly:"false" name:"client-secret"`
	RefreshToken string `description:"AzureAD Refresh Token" kind:"attribute" mode:"normal" readonly:"false" name:"refresh-token"`
	UserPID      string `description:"AzureAD User PID of the user you want to read" kind:"attribute" mode:"normal" readonly:"false" name:"user-pid"`
	UserEmail    string `description:"AzureAD User email of the user you want to read" kind:"attribute" mode:"normal" readonly:"false" name:"user-email"`
}

func (c *AzureADConfig) Validate(operation plugin.OperationType) error {
	var client *azureclient.AzureADClient
	var err error

	if c.Tenant == "" {
		return status.Error(codes.InvalidArgument, "no tenant was provided")
	}

	if c.ClientID == "" {
		return status.Error(codes.InvalidArgument, "no client id was provided")
	}

	if c.ClientSecret == "" {
		return status.Error(codes.InvalidArgument, "no client secret was provided")
	}

	if c.UserPID != "" && c.UserEmail != "" {
		return status.Error(codes.InvalidArgument, "an user PID and an user email were provided; please specify only one")
	}

	if c.RefreshToken != "" {
		client, err = azureclient.NewAzureADClientWithRefreshToken(
			context.Background(),
			c.Tenant,
			c.ClientID,
			c.ClientSecret,
			c.RefreshToken)
	} else {
		client, err = azureclient.NewAzureADClient(
			context.Background(),
			c.Tenant,
			c.ClientID,
			c.ClientSecret)
	}
	if err != nil {
		return status.Errorf(codes.Internal, "failed to connect to AzureAD, %s", err.Error())
	}

	_, errReq := client.ListUsers()

	if errReq != nil {
		return status.Errorf(codes.Internal, "failed to retrieve users from AzureAD: %s", errReq.Error())
	}

	return nil
}

func (c *AzureADConfig) Description() string {
	return "AzureAD plugin"
}
