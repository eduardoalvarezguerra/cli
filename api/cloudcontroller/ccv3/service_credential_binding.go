package ccv3

import (
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3/internal"
	"code.cloudfoundry.org/cli/resources"
)

// GetServiceCredentialBindings queries the CC API with the specified query
// and returns a slice of ServiceCredentialBindings. Additionally if Apps are
// included in the API response (by having `include=app` in the query) then the
// App names will be added into each ServiceCredentialBinding for app bindings
func (client Client) GetServiceCredentialBindings(query ...Query) ([]resources.ServiceCredentialBinding, Warnings, error) {
	var result []resources.ServiceCredentialBinding

	included, warnings, err := client.MakeListRequest(RequestParams{
		RequestName:  internal.GetServiceCredentialBindingsRequest,
		Query:        query,
		ResponseBody: resources.ServiceCredentialBinding{},
		AppendToList: func(item interface{}) error {
			result = append(result, item.(resources.ServiceCredentialBinding))
			return nil
		},
	})

	if len(included.Apps) > 0 {
		appLookup := make(map[string]resources.Application)
		for _, app := range included.Apps {
			appLookup[app.GUID] = app
		}

		for i := range result {
			result[i].AppName = appLookup[result[i].AppGUID].Name
			result[i].AppSpaceGUID = appLookup[result[i].AppGUID].SpaceGUID
		}
	}

	return result, warnings, err
}
