package testutils

import (
	"time"

	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func CreateTestAPIUser(id, displayName, email, picture string) *api.User {
	user := api.User{
		Id:          id,
		DisplayName: displayName,
		Email:       email,
		Picture:     picture,
		Identities:  make(map[string]*api.IdentitySource),
		Attributes: &api.AttrSet{
			Properties:  &structpb.Struct{Fields: make(map[string]*structpb.Value)},
			Roles:       []string{},
			Permissions: []string{},
		},
		Applications: make(map[string]*api.AttrSet),
		Metadata: &api.Metadata{
			CreatedAt: timestamppb.New(time.Now()),
			UpdatedAt: timestamppb.New(time.Now()),
		},
	}

	return &user
}

func CreateTestAzureADUser(id, displayName, email, picture, phoneNo, userName string) *models.User {

	user := models.NewUser()
	user.SetId(&id)
	user.SetDisplayName(&displayName)
	user.SetMail(&email)
	t := time.Now()
	user.SetCreatedDateTime(&t)

	return user
}
