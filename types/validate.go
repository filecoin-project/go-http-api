package types

import (
	"fmt"
	"reflect"
	"strings"
)

func RequireFields(m interface{}, params ...string) error {
	var missing []string
	mval := reflect.ValueOf(m)
	for _, param := range params {
		val := reflect.Indirect(mval).FieldByName(param)
		valKind := val.Kind().String()
		isMissing := false
		switch valKind {
		case "string":
			if val.Interface().(string) == "" {
				isMissing = true
			}
		case "uint64":
			if val.Interface().(uint64) == 0 {
				isMissing = true
			}
		case "ptr":
			if val.IsNil() {
				isMissing = true
			}
		case "slice":
			if val.IsNil() {
				isMissing = true
			}
		case "invalid":
			return fmt.Errorf("%s is not part of struct %s", param, reflect.TypeOf(m).Name())
		}
		if isMissing {
			missing = append(missing, param)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing parameters: %s", strings.Join(missing, ","))
	}
	return nil
}
