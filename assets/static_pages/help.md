# Markify

{{ toc 2 4 }}

## Markdown

**Markify** uses [goldmark](https://github.com/yuin/goldmark) markdown implementation
and [![MarkdownMark](https://raw.githubusercontent.com/dcurtis/markdown-mark/master/png/32x20.png) CommonMark](https://commonmark.org) complaint.
Quick markdown tutorial can be found at [commonmark help](https://commonmark.org/help).
Also [GFM](https://github.github.com/gfm/) extension are enabled
and *tables*, *strikethrough*, *autolinks* and *task lists* can be used too.

This page can be used as tutorial, raw code can be found [here](help/raw).

### Emphasis

```
Following text is _italic_
This is *italic too*
It is __bold__
also **bold**
```

Following text is _italic_

This is *italic too*

It is __bold__

also **bold**

### *Lists*

Lists can be ordered

1. one
1. two

or unordered

- one
- two

### *Tasklist*

- [x] one checked
- [ ] two unchecked

### *Blockquotes*

> some text
>
> second line
> continues there

### *Tables*

First Header | Second Header
------------ | -------------
Content from cell 1 | Content from cell 2
Content in the first column | Content in the second column
Content in the third line | Content in the third line too

### *Horizontal Rule*

---

### *Footnote*

That's some text with a footnote.[^1]

Some other text...

[^1]: And that's the footnote.

### *Typography*

Some punctuations can be replaced with typographic entities:

```
Ellipsis...

<<Text in quotes>>

Double dash -- to insert en dash

Tripple dash --- to insert em dash
```

Ellipsis...

<<Text in quotes>>

Double dash -- to insert en dash

Tripple dash --- to insert em dash


## Shortcodes

Additionally to pure markdown we support some pre-defined *shortcodes*
that used to embed some content on page.

### tweet

To embed tweet insert `{{ tweet <id> }}`

Example:

```
{{ tweet 1224047348109795330 }}
```
Will be displayed as:

{{ tweet 1224047348109795330 }}


### instagram

To embed Instagram post insert `{{ instagram <id> }}`

Example:

```
{{ instagram B7PFiANKX5r }}
```
Will be displayed as:

{{ instagram B7PFiANKX5r }}

### gist
To embed GitHub gist insert `{{ gist <userid> <gistid> [file] }}`

Parameter `file` is optional and used when gist contain more than one file.

Example:

```
{{ gist gvanrossum 18bdf248a679155f1381 echo_server_tulip.py }}
```
Will be displayed as:

{{ gist gvanrossum 18bdf248a679155f1381 echo_server_tulip.py }}


### table of contents

To insert Table of Contents insert `{{ toc [<highest_level> <lowest_level>] }}`

Parameter `highest_level` used to filter headers of highest levels and `lowest_level` for lower levels ones.

Example:

Insert full table of contests:
```
{{ toc }}
```
{{ toc }}

Table of contests with headers with levels from 2 to 3:
```
{{ toc 2 3 }}
```
{{ toc 2 3 }}

Also you could see table of contents display at the beginning of this document.
