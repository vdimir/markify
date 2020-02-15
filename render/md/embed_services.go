package md

import "html/template"

// EmbedTweet extension that allows embed tweets via oembed API
var EmbedTweet = &oembedExtender{
	name:        "tweet",
	urlTemplate: "https://api.twitter.com/1/statuses/oembed.json?id=%s&omit_script=true",
}

// EmbedInstagram extension that allows embed intagram posts via oembed API
var EmbedInstagram = &oembedExtender{
	name:        "instagram",
	urlTemplate: "https://api.instagram.com/oembed/?url=https://www.instagram.com/p/%s/&amp;maxwidth=420&amp;omitscript=true",
}

const gistTemplateStr = `<script type="application/javascript" ` +
	`src="https://gist.github.com/{{ index . 0 }}/{{ index . 1 }}.js` +
	`{{if len . | eq 3 }}?file={{ index . 2 }}{{end}}">` +
	`</script>`

// EmbedGist extension that allows embed gists
var EmbedGist = &tplEmbedExtender{
	name: "gist",
	tpl:  template.Must(template.New("gist").Parse(gistTemplateStr)),
}
