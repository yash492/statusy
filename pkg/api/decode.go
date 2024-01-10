package api

import (
	"encoding/json"
	"net/http"
)

type Validater interface {
	Validate() error
}

func Decode(r *http.Request, v any) error {

	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}

	if validater, ok := v.(Validater); ok {
		if err := validater.Validate(); err != nil {
			return err
		}
	}

	return nil

}
