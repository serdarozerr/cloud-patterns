package pkg

import (
	"encoding/json"
	"io"
)

type Validator interface {
	Validate() map[string]string
}

func ValidateEmptyField(field interface{}) bool {

	switch v := field.(type) {
	case string:
		return v == ""
	case int:
		return v == 0
	default:
		return false
	}

}

func ValidateStruct[T Validator](rBody io.ReadCloser) (data T, errors map[string]string, err error) {

	err = json.NewDecoder(rBody).Decode(&data)
	if err != nil {
		return data, nil, err
	}

	errors = data.Validate()
	if len(errors) > 0 {
		return data, errors, nil
	}

	return data, nil, nil

}
