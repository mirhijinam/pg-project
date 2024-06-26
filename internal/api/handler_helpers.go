package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type envelope map[string]interface{}

func isAdmin(t string) bool {
	adminToken := os.Getenv("ADMIN_TOKEN")

	switch t {
	case adminToken:
		return true
	default:
		return false
	}

}

func getSudoNameCommand(cmdRaw string) (name string, isSudo bool) {
	trimmed := strings.TrimSpace(cmdRaw)
	space := regexp.MustCompile(`\s+`)
	cleaned := space.ReplaceAllString(trimmed, " ")

	cmdRawSlice := strings.Split(cleaned, " ")

	name = strings.Join(cmdRawSlice[:1], " ")
	isSudo = false

	if cmdRawSlice[0] == "sudo" {
		name = strings.Join(cmdRawSlice[:2], " ")
		isSudo = true
	}

	return
}

func readJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	maxBytesBody := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytesBody))
	dec := json.NewDecoder(r.Body)

	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesBody)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

func writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {

	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')
	for key, val := range headers {
		w.Header()[key] = val
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}
