package app

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	chirender "github.com/go-chi/render"
	"github.com/pkg/errors"
)

// CreatePasteRequest represent new paste obtaned from user
type CreatePasteRequest struct {
	Text      string        `json:"text"`
	Syntax    string        `json:"syntax"`
	UserToken string        `json:"token,omitempty"`
	TTL       time.Duration `json:"ttl,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type PreviewResponse struct {
	Body string `json:"body"`
}

type CreateResponse struct {
	Path string `json:"path"`
}

func parseCreateRequestJson(r *http.Request) (*CreatePasteRequest, error) {
	if chirender.GetRequestContentType(r) == chirender.ContentTypeJSON {
		parsedReq := &CreatePasteRequest{}
		err := json.NewDecoder(r.Body).Decode(parsedReq)
		return parsedReq, err
	}

	return nil, errors.Errorf("can't parse create request: unknown content type %q", r.Header.Get("Content-Type"))
}

func (app *App) handleApiPreview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	createReq, err := parseCreateRequestJson(r)
	if err != nil {
		app.apiResponseWithCode(w, http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}
	if err = app.validatePasteRequest(createReq); err != nil {
		app.apiResponseWithCode(w, http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}
	doc, err := app.converter.Convert(strings.NewReader(createReq.Text), createReq.Syntax)
	if err != nil {
		app.apiResponseWithCode(w, http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}
	app.apiResponse(w, PreviewResponse{doc.Body})
}

func (app *App) handleApiCreateDocument(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	createReq, err := parseCreateRequestJson(r)
	if err != nil {
		app.apiResponseWithCode(w, http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}
	if err = app.validatePasteRequest(createReq); err != nil {
		app.apiResponseWithCode(w, http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}
	docID, err := app.savePaste(createReq)
	if err != nil {
		app.apiResponseWithCode(w, http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}
	app.apiResponse(w, CreateResponse{"/p/" + docID})

}

func (app *App) handleApiNotFound(w http.ResponseWriter, r *http.Request) {
	app.apiResponseWithCode(w, http.StatusNotFound, ErrorResponse{"Page not found"})
}

func (app *App) apiResponse(w http.ResponseWriter, resp interface{}) {
	app.apiResponseWithCode(w, http.StatusOK, resp)
}

func (app *App) apiResponseWithCode(w http.ResponseWriter, code int, resp interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	err := enc.Encode(resp)
	if err != nil {
		log.Printf("[ERROR] can't write data: %s", err)
	}
}
