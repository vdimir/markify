package render

import (
	"html"
	"io"
	"io/ioutil"
)

type plainText struct {
}

func (r *plainText) Convert(reader io.Reader) (*Document, error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return &Document{
		Body: "<pre><code>" + html.EscapeString(string(data)) + "</code></pre>",
	}, nil
}
