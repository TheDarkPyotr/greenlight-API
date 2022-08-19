package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type envelope map[string]interface{}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {

	jsonBody, err := json.MarshalIndent(data, "", "\t")
	if err != nil {

		return err
	}

	//for CLI visualization
	jsonBody = append(jsonBody, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(jsonBody)
	if err != nil {
		return err
	}

	return nil

}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {

	//limit max number of bytes per request
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(dst)
	if err != nil {
		//let's triage the error: only 5 type of error are possible
		var syntaxError *json.SyntaxError
		var unmarshallTypeError *json.UnmarshalTypeError
		var invalidUnmarshallError *json.InvalidUnmarshalError

		switch {
		//errors.As compare type of errors
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly formed JSON at character %d", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshallTypeError):
			if unmarshallTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshallTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type at character %d", unmarshallTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case errors.As(err, &invalidUnmarshallError):
			panic(err)

			//curl -d '{"title": "Moana", "rating":"PG"}' localhost:4000/v1/movies --> avoid non-existent field
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			field := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", field)

		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		default:
			return err

		}
	}

	//to avoid multiple JSON value in the same request
	//curl {title": "Moana"}{"title": "Top Gun"}' localhost:4000/v1/movies
	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil

}
func (app *application) readIDParam(r *http.Request) (int64, error) {

	parameters := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(parameters.ByName("id"), 10, 64)
	if err != nil || id < 1 {

		return 0, errors.New("Invalid ID Param")
	}

	return id, nil

}
