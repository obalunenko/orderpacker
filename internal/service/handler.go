package service

import (
	"encoding/json"
	"io"
	log "log/slog"
	"net/http"

	"github.com/obalunenko/orderpacker/internal/packer"
)

func NewRouter(p *packer.Packer) *http.ServeMux {
	mux := http.NewServeMux()

	handler := PackHandler(p)
	handler = logRequestMiddleware(handler)
	handler = logResponseMiddleware(handler)
	handler = requestIDMiddleware(handler)
	handler = recoverMiddleware(handler)

	mux.Handle("/pack", handler)

	return mux
}

func PackHandler(p *packer.Packer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
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
				log.Error("Error closing request body", "error", err)
			}
		}()

		var req PackRequest

		if err = json.Unmarshal(b, &req); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

			return
		}

		order := p.PackOrder(req.Items)

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
	})
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

	return resp
}
