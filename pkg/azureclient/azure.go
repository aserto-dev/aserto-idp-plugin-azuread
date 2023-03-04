package azureclient

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	auth "github.com/microsoft/kiota-authentication-azure-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	adusers "github.com/microsoftgraph/msgraph-sdk-go/users"
)

type AzureADClient struct {
	clientSecretCredential *azidentity.ClientSecretCredential
	refreshTokenCredential *RefreshTokenCredential
	appClient              *msgraphsdk.GraphServiceClient
}

func NewAzureADClientWithSecret(ctx context.Context, tenant, clientID, clientSecret string) (*AzureADClient, error) {
	c := &AzureADClient{}

	credential, err := azidentity.NewClientSecretCredential(tenant, clientID, clientSecret, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create an Azure secret credential: %s", err.Error())
	}

	c.appClient, err = getAppClient(credential)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func NewAzureADClientWithRefreshToken(ctx context.Context, tenant, clientID, refreshToken string) (*AzureADClient, error) {
	c := &AzureADClient{}

	credential, err := NewRefreshTokenCredential(ctx, tenant, clientID, refreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create Refresh Token credential: %s", err.Error())
	}

	c.appClient, err = getAppClient(credential)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *AzureADClient) ListUsers() (models.UserCollectionResponseable, error) {
	return c.listUsers("")
}

func (c *AzureADClient) GetUser(name string) (models.UserCollectionResponseable, error) {
	return c.listUsers(name)
}

func (c *AzureADClient) listUsers(filter string) (models.UserCollectionResponseable, error) {
	query := adusers.UsersRequestBuilderGetQueryParameters{
		Select:  []string{"displayName", "id", "mail", "createdDateTime", "mobilePhone"},
		Orderby: []string{"displayName"},
		Filter:  &filter,
	}
	return c.appClient.Users().
		Get(context.Background(),
			&adusers.UsersRequestBuilderGetRequestConfiguration{
				QueryParameters: &query,
			})
}

func getAppClient(credential azcore.TokenCredential) (*msgraphsdk.GraphServiceClient, error) {
	authProvider, err := auth.NewAzureIdentityAuthenticationProviderWithScopes(credential, []string{
		"https://graph.microsoft.com/.default",
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create Azure identity provider: %s", err.Error())
	}

	// Create a request adapter using the auth provider
	adapter, err := msgraphsdk.NewGraphRequestAdapter(authProvider)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create Azure AD Graph request adapter: %s", err.Error())
	}

	// Create a Graph client using request adapter
	client := msgraphsdk.NewGraphServiceClient(adapter)
	return client, nil
}
