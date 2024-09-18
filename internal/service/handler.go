package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"

	log "github.com/obalunenko/logger"

	"github.com/obalunenko/orderpacker/internal/packer"
	"github.com/obalunenko/orderpacker/internal/service/assets"
)

// ErrEmptyItems is returned when items is zero or empty.
var ErrEmptyItems = errors.New("empty items")

func NewRouter(p *packer.Packer) *http.ServeMux {
	mux := http.NewServeMux()

	mw := []func(http.Handler) http.Handler{
		logRequestMiddleware,
		logResponseMiddleware,
		requestIDMiddleware,
		recoverMiddleware,
		loggerMiddleware,
		corsMiddleware,
	}

	mwApply := func(h http.Handler) http.Handler {
		for i := range mw {
			h = mw[i](h)
		}

		return h
	}

	mux.Handle("/", mwApply(indexHandler()))
	mux.Handle("/favicon.ico", mwApply(faviconHandler()))

	// Group api/v1 routes.
	mux.Handle("/api/v1/pack", mwApply(packHandler(p)))

	return mux
}

func indexHandler() http.HandlerFunc {
	homePageHTML := string(assets.MustLoad("index.gohtml"))
	homePageTmpl := template.Must(template.New("index").Parse(homePageHTML))

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		if err := homePageTmpl.Execute(w, nil); err != nil {
			http.Error(w, "failed to execute template", http.StatusInternalServerError)

			return
		}
	}
}

func faviconHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}
}

// packHandler - handler for /pack endpoint.
//
//	@Summary		Get the number of packs needed to ship to a customer
//	@Tags			pack
//	@Description	Calculates the number of packs needed to ship to a customer
//	@ID				orderpacker-pack	post
//	@Accept			json
//	@Produce		json
//	@Param			data	body		PackRequest				true	"Request data"
//	@Success		200		{object}	PackResponse			"Successful response with packs data"
//	@Failure		400		{object}	badRequestError			"Invalid request data
//	@Failure		405		{object}	methodNotAllowedError	"Method not allowed"
//	@Failure		500		{object}	internalServerError		"Internal server error"
//	@Router			/api/v1/pack [post]
func packHandler(p *packer.Packer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			makeResponse(
				r.Context(),
				w,
				http.StatusMethodNotAllowed,
				PackResponse{},
				errors.New(http.StatusText(http.StatusMethodNotAllowed)),
			)

			return
		}

		b, err := io.ReadAll(r.Body)
		if err != nil {
			makeResponse(
				r.Context(),
				w,
				http.StatusBadRequest,
				PackResponse{},
				fmt.Errorf("failed to read request body: %w", err),
			)

			return
		}

		defer func() {
			if err = r.Body.Close(); err != nil {
				log.WithError(r.Context(), err).Error("Error closing request body")
			}
		}()

		var req PackRequest

		if err = json.Unmarshal(b, &req); err != nil {
			makeResponse(r.Context(), w, http.StatusBadRequest, PackResponse{}, fmt.Errorf("failed to unmarshal request: %w", err))

			return
		}

		items, err := fromAPIRequest(req)
		if err != nil {
			makeResponse(
				r.Context(),
				w,
				http.StatusBadRequest,
				PackResponse{},
				fmt.Errorf("invalid request: %w", err),
			)

			return
		}

		order := p.PackOrder(r.Context(), items)

		resp := toAPIResponse(order)

		b, err = json.Marshal(resp)
		if err != nil {
			makeResponse(
				r.Context(),
				w,
				http.StatusInternalServerError,
				PackResponse{},
				fmt.Errorf("failed to marshal response: %w", err),
			)

			return
		}

		w.Header().Set("Content-Type", "application/json")

		makeResponse(r.Context(), w, http.StatusOK, resp, nil)
	}
}

func makeResponse(ctx context.Context, w http.ResponseWriter, code int, resp PackResponse, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	var response any

	response = resp

	if err != nil {
		log.WithError(ctx, err).Error("Error processing request")

		response = newHTTPError(ctx, code, err.Error())
	}

	if err = json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
