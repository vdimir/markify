## Markdown

**Markify** uses [goldmark](https://github.com/yuin/goldmark) markdown implementation
and [![MarkdownMark](https://raw.githubusercontent.com/dcurtis/markdown-mark/master/png/32x20.png) CommonMark](https://commonmark.org) complaint.
Quick markdown tutorial can be found at [commonmark help](https://commonmark.org/help).
Also [GFM](https://github.github.com/gfm/) extension are enabled
and *tables*, *strikethrough*, *autolinks* and *task lists* can be used too.

Raw code on this page can be found [here](markdown/raw).

{{ toc 3 4 }}

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

Some punctuation can be replaced with typographic entities:

```
Ellipsis...

<<Text in quotes>>

Double dash -- to insert en dash

Tripple dash --- to insert em dash
```

Ellipsis...

<<Text in quotes>>

Double dash -- to insert en dash

Triple dash --- to insert em dash

### Extensions

Additionally to basic markdown syntax markify supports some extensions like Table of Contents, social media embedding. [Read more](/info/extensions)
