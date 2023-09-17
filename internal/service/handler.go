package service

import (
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"net/http"
	"sort"

	log "github.com/obalunenko/logger"

	"github.com/obalunenko/orderpacker/internal/packer"
	"github.com/obalunenko/orderpacker/internal/service/assets"
)

func NewRouter(p *packer.Packer) *http.ServeMux {
	mux := http.NewServeMux()

	mw := []func(http.Handler) http.Handler{
		logRequestMiddleware,
		logResponseMiddleware,
		requestIDMiddleware,
		recoverMiddleware,
		loggerMiddleware,
	}

	mwApply := func(h http.Handler) http.Handler {
		for i := range mw {
			h = mw[i](h)
		}

		return h
	}

	mux.Handle("/", mwApply(indexHandler()))
	mux.Handle("/pack", mwApply(packHandler(p)))
	mux.Handle("/favicon.ico", mwApply(faviconHandler()))

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

func packHandler(p *packer.Packer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

			return
		}

		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

			return
		}

		defer func() {
			if err = r.Body.Close(); err != nil {
				log.WithError(r.Context(), err).Error("Error closing request body")
			}
		}()

		var req PackRequest

		if err = json.Unmarshal(b, &req); err != nil {
			log.WithError(r.Context(), err).Error("Error unmarshalling request body")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

			return
		}

		items, err := fromAPIRequest(req)
		if err != nil {
			log.WithError(r.Context(), err).Error("Invalid request body")
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		order := p.PackOrder(r.Context(), items)

		resp := toAPIResponse(order)

		b, err = json.Marshal(resp)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		if _, err = w.Write(b); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

			return
		}
	}
}

// ErrEmptyItems is returned when items is zero or empty.
var ErrEmptyItems = errors.New("empty items")

func fromAPIRequest(req PackRequest) (uint, error) {
	if req.Items == 0 {
		return 0, ErrEmptyItems
	}

	return req.Items, nil
}

func toAPIResponse(boxes []uint) PackResponse {
	var resp PackResponse

	orderMap := make(map[uint]uint)
	for i := range boxes {
		orderMap[boxes[i]]++
	}

	for k, v := range orderMap {
		resp.Packs = append(resp.Packs, Pack{
			Box:      k,
			Quantity: v,
		})
	}

	sort.Slice(resp.Packs, func(i, j int) bool {
		return resp.Packs[i].Box > resp.Packs[j].Box
	})

	return resp
}
