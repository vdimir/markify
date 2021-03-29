package app

import (
	"encoding/json"
	chirender "github.com/go-chi/render"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

type CreatePasteRequest struct {
	Text      string `json:"text"`
	Syntax    string `json:"syntax"`
	UserToken string `json:"token,omitempty"`
	Ttl       time.Duration
}

func ParseCreatePasteRequest(r *http.Request) (*CreatePasteRequest, error) {
	if chirender.GetRequestContentType(r) == chirender.ContentTypeForm {
		token := ""
		uidCookie, err := r.Cookie("user_id")
		if err == nil {
			token = uidCookie.Value
		}
		return &CreatePasteRequest{
			Text:      r.FormValue("data"),
			Syntax:    r.FormValue("syntax"),
			UserToken: token,
		}, nil
	}

	if chirender.GetRequestContentType(r) == chirender.ContentTypeJSON {
		parsedReq := &CreatePasteRequest{}
		err := json.NewDecoder(r.Body).Decode(parsedReq)
		return parsedReq, err
	}

	return nil, errors.Errorf("can't parse create request: unknown content type %q", r.Header.Get("Content-Type"))
}

