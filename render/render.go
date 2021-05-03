package render

import (
	"io"
	"io/ioutil"

	"github.com/vdimir/markify/render/markdown"
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
