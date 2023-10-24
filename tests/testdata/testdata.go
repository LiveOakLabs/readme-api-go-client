package testdata

import "encoding/json"

func ToJSON(data interface{}) string {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	return string(dataJSON)
}
