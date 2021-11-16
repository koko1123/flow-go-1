package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/onflow/flow-go/access"
	"github.com/onflow/flow-go/engine/access/rest/generated"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"strings"
)

type StatusError interface {
	error
	Status() int
}

type RestError struct {
	status      int
	userMessage string
	err         error
}

// NewRestError creates an error returned to user with provided status
// user displayed message and internal error
func NewRestError(status int, msg string, err error) *RestError {
	return &RestError{
		status:      status,
		userMessage: msg,
		err:         err,
	}
}

// NewBadRequest creates a new bad request rest error.
func NewBadRequest(msg string, err error) *RestError {
	return &RestError{
		status:      http.StatusBadRequest,
		userMessage: msg,
		err:         err,
	}
}

func (e *RestError) Status() int {
	return e.status
}

func (e *RestError) Error() string {
	return e.err.Error()
}

type Handler struct {
	logger      zerolog.Logger
	backend     access.API
	method      string
	pattern     string
	name        string
	handlerFunc func(
		w http.ResponseWriter,
		r *http.Request,
		vars map[string]string,
		backend access.API,
		logger zerolog.Logger,
	) (interface{}, StatusError)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// execute handler function and check for error
	response, err := h.handlerFunc(w, r, mux.Vars(r), h.backend, h.logger)
	if err != nil {
		switch e := err.(type) {
		case StatusError:
			errorResponse(w, e.Status(), e.Error(), h.logger)
		default:
			errorResponse(w, http.StatusInternalServerError, e.Error(), h.logger)
		}

		// stop going further
		return
	}

	encodedResponse, encErr := json.Marshal(response)
	if encErr != nil {
		h.logger.Error().Err(err).Msg("failed to encode response")
		errorResponse(w, http.StatusInternalServerError, "error generating response", h.logger)
		return
	}

	_, writeErr := w.Write(encodedResponse)
	if writeErr != nil {
		h.logger.Error().Err(err).Msg("failed to write response")
		errorResponse(w, http.StatusInternalServerError, "error generating response", h.logger)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// errorResponse sends an HTTP error response to the client with the given return code
// and a model error with the given response message in the response body
func errorResponse(
	w http.ResponseWriter,
	returnCode int,
	responseMessage string,
	logger zerolog.Logger,
) {
	modelError := generated.ModelError{
		Code:    int32(returnCode),
		Message: responseMessage,
	}
	encodedError, err := json.Marshal(modelError)
	if err != nil {
		logger.Error().Str("response_message", responseMessage).Msg("failed to json encode error message")
		return
	}

	w.WriteHeader(returnCode)
	_, err = w.Write(encodedError)
	if err != nil {
		logger.Error().Err(err).Msg("failed to send error response")
	}
}

func jsonDecode(body io.ReadCloser, dst interface{}) error {
	// validate size

	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &badRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			return &badRequest{status: http.StatusBadRequest, msg: msg}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &badRequest{status: http.StatusBadRequest, msg: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &badRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &badRequest{status: http.StatusBadRequest, msg: msg}

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return &badRequest{status: http.StatusRequestEntityTooLarge, msg: msg}

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain h single JSON object"
		return &badRequest{status: http.StatusBadRequest, msg: msg}
	}

	return nil
}
