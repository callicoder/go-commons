package errors

import (
	"encoding/json"
	"testing"

	"github.com/callicoder/go-commons/errors/codes"
	"github.com/stretchr/testify/assert"
)

func TestNewf(t *testing.T) {
	err := WithCode(codes.NotFound).Newf("Resource with id %d not found", 1234)
	assert.EqualError(t, err, "Resource with id 1234 not found")
	assert.EqualError(t, Cause(err), err.Error())
}

func TestWrap(t *testing.T) {
	err1 := New("Error1")
	err2 := Wrap(err1, "Error2")
	err3 := Wrap(err2, "Error3")

	assert.EqualError(t, Cause(err3), "Error1")
}

func TestWithDetailWithCode(t *testing.T) {
	err := WithDetails(Detail{
		Resource: "user",
		Field:    "email",
		Value:    "abc@@@g.com",
	}).WithDetails(Detail{
		Resource: "user",
		Field:    "name",
		Value:    "%#$1@",
	}).WithCode(codes.BadRequest).New("Invalid or missing parameters")

	assert.EqualError(t, err, "Invalid or missing parameters")

	expectedRes := `{"code":"bad_request","message":"Invalid or missing parameters","details":[{"resource":"user","field":"email","value":"abc@@@g.com"},{"resource":"user","field":"name","value":"%#$1@"}]}`
	res, err := json.Marshal(err)
	assert.Equal(t, string(res), expectedRes)
}
