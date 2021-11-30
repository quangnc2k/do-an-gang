package something

import (
	"log"
	"strings"
)

func ExtractFromJsonMap(m map[string]interface{}, key string) (v interface{}) {
	index := strings.IndexAny(key,`.`)
	if index >= 0 {
		pre := key[:index]
		key = key[index+1:]
		n, ok := m[pre].(map[string]interface{})
		if !ok {
			log.Fatalln("error traversing map object, check key: ", key)
		}
		v = ExtractFromJsonMap(n, key)
		if len(n) == 0 {
			delete(m, pre)
		}
	} else {
		v = m[key]
		delete(m, key)
	}

	return v
}
