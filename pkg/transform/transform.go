package transform

import (
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

const (
	Provider = "azuread"
)

func ToAzureAD(in *api.User) *models.User {

	user := models.NewUser()
	user.SetId(&in.Id)
	user.SetDisplayName(&in.DisplayName)
	user.SetMail(&in.Email)

	return user
}

// Transform AzureAD user definition into Aserto Edge User object definition.
func Transform(in models.Userable) *api.User {

	user := api.User{
		Id:          *in.GetId(),
		DisplayName: *in.GetDisplayName(),
		//Picture:     *in.GetPhoto().GetId(),
		Identities: make(map[string]*api.IdentitySource),
		Attributes: &api.AttrSet{
			Properties:  &structpb.Struct{Fields: make(map[string]*structpb.Value)},
			Roles:       []string{},
			Permissions: []string{},
		},
		Applications: make(map[string]*api.AttrSet),
		Metadata: &api.Metadata{
			CreatedAt: timestamppb.New(*in.GetCreatedDateTime()),
			UpdatedAt: timestamppb.New(*in.GetCreatedDateTime()),
		},
	}

	email := in.GetMail()
	if email != nil {
		user.Email = *email
	} else {
		user.Email = ""
	}

	user.Identities[*in.GetId()] = &api.IdentitySource{
		Kind:     api.IdentityKind_IDENTITY_KIND_PID,
		Provider: Provider,
		Verified: true,
	}

	if email != nil && *email != "" {
		user.Identities[*email] = &api.IdentitySource{
			Kind:     api.IdentityKind_IDENTITY_KIND_EMAIL,
			Provider: Provider,
			Verified: true,
		}
	}

	phone := in.GetMobilePhone()
	if phone != nil && *phone != "" {
		user.Identities[*phone] = &api.IdentitySource{
			Kind:     api.IdentityKind_IDENTITY_KIND_PHONE,
			Provider: Provider,
			Verified: false,
		}
	}

	return &user
}
