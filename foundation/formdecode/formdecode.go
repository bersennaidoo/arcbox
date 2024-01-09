package formdecode

import "github.com/gorilla/schema"

func New() *schema.Decoder {

	decoder := schema.NewDecoder()

	return decoder
}
