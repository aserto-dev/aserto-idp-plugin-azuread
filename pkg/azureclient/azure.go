package azureclient

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	auth "github.com/microsoft/kiota-authentication-azure-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	adusers "github.com/microsoftgraph/msgraph-sdk-go/users"
)

type AzureADClient struct {
	clientSecretCredential *azidentity.ClientSecretCredential
	appClient              *msgraphsdk.GraphServiceClient
}

func NewAzureADClient(ctx context.Context, tenant string, clientID string, clientSecret string) (*AzureADClient, error) {
	c := &AzureADClient{}

	credential, err := azidentity.NewClientSecretCredential(tenant, clientID, clientSecret, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create an Azure secret credential: %s", err.Error())
	}

	c.clientSecretCredential = credential

	// Create an auth provider using the credential
	authProvider, err := auth.NewAzureIdentityAuthenticationProviderWithScopes(c.clientSecretCredential, []string{
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
	c.appClient = client

	return c, nil
}

func (c *AzureADClient) ListUsers() (models.UserCollectionResponseable, error) {
	var topValue int32 = 25
	query := adusers.UsersRequestBuilderGetQueryParameters{
		// Only request specific properties
		Select: []string{"displayName", "id", "mail"},
		// Get at most 25 results
		Top: &topValue,
		// Sort by display name
		Orderby: []string{"displayName"},
	}

	return c.appClient.Users().
		Get(context.Background(),
			&adusers.UsersRequestBuilderGetRequestConfiguration{
				QueryParameters: &query,
			})
}

func (c *AzureADClient) GetUser(name string) (models.UserCollectionResponseable, error) {
	var topValue int32 = 1
	query := adusers.UsersRequestBuilderGetQueryParameters{
		// Only request specific properties
		Select: []string{"displayName", "id", "mail"},
		// Get at most 1 results
		Top: &topValue,
		// Sort by display name
		Orderby: []string{"displayName"},
		Filter:  &name,
	}

	return c.appClient.Users().
		Get(context.Background(),
			&adusers.UsersRequestBuilderGetRequestConfiguration{
				QueryParameters: &query,
			})
}
