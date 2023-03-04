package azureclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

type RefreshTokenCredential struct {
	clientID     string
	refreshToken string
	tenantID     string
}

func NewRefreshTokenCredential(ctx context.Context, tenantID, clientID, refreshToken string) (*RefreshTokenCredential, error) {
	c := &RefreshTokenCredential{
		clientID:     clientID,
		tenantID:     tenantID,
		refreshToken: refreshToken,
	}
	return c, nil
}

func (c *RefreshTokenCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return getAccessToken(c.tenantID, c.clientID, c.refreshToken)
}

func getAccessToken(tenantID string, clientID string, refreshToken string) (azcore.AccessToken, error) {
	accessToken := azcore.AccessToken{}

	url := "https://login.microsoftonline.com/" + tenantID + "/oauth2/v2.0/token"
	data := fmt.Sprintf("grant_type=refresh_token&client_id=%s&refresh_token=%s",
		clientID, refreshToken)
	payload := strings.NewReader(data)

	// create the request and execute it
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return accessToken, err
	}

	// process the response
	defer res.Body.Close()
	var responseData map[string]interface{}
	body, _ := ioutil.ReadAll(res.Body)

	// unmarshal the json into a string map
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return accessToken, err
	}

	// retrieve the access token and expiration
	accessToken.Token = responseData["access_token"].(string)
	expiresIn := responseData["expires_in"].(int)
	accessToken.ExpiresOn = time.Now().Add(time.Second * time.Duration(expiresIn))
	return accessToken, nil
}
