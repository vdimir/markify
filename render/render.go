package render

import (
	"github.com/pkg/errors"
	"github.com/vdimir/markify/render/markdown"
	"io"
	"io/ioutil"
)

type Document struct {
	Title   string
	Preview string
	Body    string
}

type DocConverter struct {
	md   *markdown.Converter
	code *plainText
}

func NewConverter() *DocConverter {
	return &DocConverter{
		md:   markdown.NewConverter(),
		code: &plainText{},
	}
}

func (r *DocConverter) SupportSyntax(syntax string) error {
	if syntax == "markdown" || syntax == "" {
		return nil
	}
	return errors.Errorf("syntax %q is not supported", syntax)
}

func (r *DocConverter) Convert(reader io.Reader, syntax string) (*Document, error) {
	if syntax == "markdown" {
		text, err := ioutil.ReadAll(reader)
		if err != nil {
			return nil, err
		}
		mdDoc, err := r.md.Convert(text)
		if err != nil {
			return nil, err
		}
		return &Document{
			Preview: mdDoc.Preview,
			Title:   mdDoc.Title,
			Body:    mdDoc.Body,
		}, nil
	}
	return r.code.Convert(reader)
}
