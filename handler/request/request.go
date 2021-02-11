package request

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func BindJSON(req *http.Request, i interface{}) error {
	if err := json.NewDecoder(req.Body).Decode(i); err != nil {
		if ute, ok := err.(*json.UnmarshalTypeError); ok {
			return fmt.Errorf("%w: %v", err, fmt.Sprintf("Unmarshal type error: expected=%v, got=%v, field=%v, offset=%v", ute.Type, ute.Value, ute.Field, ute.Offset))
		} else if se, ok := err.(*json.SyntaxError); ok {
			return fmt.Errorf("%w: %v", err, fmt.Sprintf("Syntax error: offset=%v, error=%v", se.Offset, se.Error()))
		}
		return err
	}
	return nil
}
