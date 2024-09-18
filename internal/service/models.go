package service

import (
	"context"
	"net/http"

	log "github.com/obalunenko/logger"
)

// PackRequest represents a request to pack items.
type PackRequest struct {
	Items uint `json:"items" format:"uint" example:"543"`
}

// Pack represents a pack of items.
type Pack struct {
	Box      uint `json:"box" format:"uint" example:"50"`
	Quantity uint `json:"quantity" format:"uint" example:"3"`
}

// PackResponse represents a response to a pack request.
type PackResponse struct {
	Packs []Pack `json:"packs,omitempty"`
}

// HTTPError represents an HTTP error.
type HTTPError interface {
	// StatusCode returns the status code of the error.
	StatusCode() int
	// Message returns the message of the error.
	Message() string
}

func newHTTPError(ctx context.Context, code int, msg string) HTTPError {
	switch code {
	case http.StatusBadRequest:
		return newBadRequestError(msg)
	case http.StatusMethodNotAllowed:
		return newMethodNotAllowedError(msg)
	case http.StatusInternalServerError:
		return newInternalServerError(msg)
	default:
		log.WithField(ctx, "code", code).Warn("Unknown error code")

		return newInternalServerError(msg)
	}
}
func newBadRequestError(msg string) HTTPError {
	return badRequestError{
		Code: http.StatusBadRequest,
		Msg:  msg,
	}
}

type badRequestError struct {
	Code int    `json:"code" example:"400"`
	Msg  string `json:"message" example:"Bad request"`
}

func (e badRequestError) StatusCode() int {
	return e.Code
}

func (e badRequestError) Message() string {
	return e.Msg
}

type internalServerError struct {
	Code int    `json:"code" example:"500"`
	Msg  string `json:"message" example:"Internal server error"`
}

func newInternalServerError(msg string) HTTPError {
	return internalServerError{
		Code: http.StatusInternalServerError,
		Msg:  msg,
	}
}

func (e internalServerError) StatusCode() int {
	return e.Code
}

func (e internalServerError) Message() string {
	return e.Msg
}

type methodNotAllowedError struct {
	Code int    `json:"code" example:"405"`
	Msg  string `json:"message" example:"Method not allowed"`
}

func newMethodNotAllowedError(msg string) HTTPError {
	return methodNotAllowedError{
		Code: http.StatusMethodNotAllowed,
		Msg:  msg,
	}
}

func (e methodNotAllowedError) StatusCode() int {
	return e.Code
}

func (e methodNotAllowedError) Message() string {
	return e.Msg
}
