package types

import (
	"fmt"
	"reflect"
	"strings"
)

func RequireFields(m interface{}, params ...string) error {
	var missing []string
	mval := reflect.ValueOf(m)
	for _, el := range params {
		val := reflect.Indirect(mval).FieldByName(el)
		valType := val.Type().Name()
		if valType == "string" && val.Interface().(string) == "" ||
			valType == "uint64" && val.Interface().(uint64) == 0 {
			missing = append(missing, el)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing parameters: %s", strings.Join(missing, ","))
	}
	return nil
}
