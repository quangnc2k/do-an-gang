package something

import "encoding/json"

func CombineAsMetadata(args ...interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	for _, arg := range args {
		var tmp map[string]interface{}

		data, err := json.Marshal(arg)
		if err != nil {
			continue
		}
		err = json.Unmarshal(data, &tmp)
		if err != nil {
			continue
		}

		for k, v := range tmp {
			m[k] = v
		}
	}

	return m
}
