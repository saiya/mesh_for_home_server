package interfaces

import (
	"encoding/json"
	"fmt"
)

func toJSON(any interface{}) string {
	bytes, err := json.Marshal(any)
	if err != nil {
		return fmt.Sprintf("%v", err)
	}
	return string(bytes)
}
