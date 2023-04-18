# aserto-idp-plugin-azuread
IDP Plugin for Aserto

## msgraph sdk generation

```
dotnet tool install --global Microsoft.OpenApi.Hidi
dotnet tool install --global Microsoft.OpenApi.Kiota
wget https://raw.githubusercontent.com/microsoftgraph/msgraph-metadata/master/openapi/v1.0/openapi.yaml
hidi transform -d openapi.yaml -c postman.json -o ./msgraph.yaml
kiota generate -l go -d msgraph.yaml -n github.com/aserto-dev/aserto-idp-plugin-azuread/pkg/msgraph -o ./pkg/msgraph -c msgraph
```