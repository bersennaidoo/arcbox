package formdecode

import "github.com/go-playground/form/v4"

func New() *form.Decoder {

	decoder := form.NewDecoder()

	return decoder
}
