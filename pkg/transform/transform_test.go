package transform_test

import (
	"reflect"
	"testing"

	"github.com/aserto-dev/aserto-idp-plugin-azuread/pkg/msgraph/models"
	azureADTestUtils "github.com/aserto-dev/aserto-idp-plugin-azuread/pkg/testutils"
	"github.com/aserto-dev/aserto-idp-plugin-azuread/pkg/transform"
	"github.com/stretchr/testify/require"
)

func TestTransformToAzureAD(t *testing.T) {
	assert := require.New(t)
	apiUser := azureADTestUtils.CreateTestAPIUser("1", "Name", "email", "pic")

	azureadUser := transform.ToAzureAD(apiUser)

	assert.True(reflect.TypeOf(azureadUser) == reflect.TypeOf(models.NewUser()), "the returned object should be *models.User")
	assert.Equal("Name", *(*azureadUser).GetDisplayName(), "should correctly detect the display name")
	assert.Equal("email", *(*azureadUser).GetMail(), "should correctly populate the email")
}

func TestTransform(t *testing.T) {
	assert := require.New(t)
	azureadUser := azureADTestUtils.CreateTestAzureADUser("1", "Name", "email", "pic", "+40722332233", "userName")

	apiUser := *transform.Transform(azureadUser)

	assert.Equal("1", apiUser.Id, "should correctly populate the id")
	assert.Equal("Name", apiUser.DisplayName, "should correctly detect the displayname")
	assert.Equal("email", apiUser.Email, "should correctly populate the email")
}
