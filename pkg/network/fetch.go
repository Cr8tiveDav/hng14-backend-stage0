package network

import (
	"encoding/json"
	"fmt"
	"hng14-stage0-api-data-processing/pkg/utils"
	"net/http"
)

func FetchData(url, apiName string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return utils.External(fmt.Sprintf("%s returned an invalid response", apiName), err)
	}

	// Defer response body close
	defer resp.Body.Close()

	// Check if the API call was successful
	if resp.StatusCode != http.StatusOK {
		return utils.External(fmt.Sprintf("%s returned an invalid response", apiName), err)
	}

	// Decode response from the API into the target pointer struct
	err = json.NewDecoder(resp.Body).Decode(target)
	if err != nil {
		return utils.Internal("Failed to decode response", err)
	}

	return nil
}
