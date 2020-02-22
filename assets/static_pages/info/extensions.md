## Markdown extensions

Additionally to [pure markdown](/info/markdown) we support some pre-defined *shortcodes*
used to embed some content on page.

Raw code on this page can be found [here](extensions/raw).

{{ toc 3 6 }}

### Table of Contents

To insert Table of Contents insert `{{ toc [<highest_level> <lowest_level>] }}`

Parameter `highest_level` used to filter headers of highest levels and `lowest_level` for lower levels ones.

Insert table of contests:
```
{{ toc 3 6 }}
```

You can see result at the [beginning](#markdown-extenstions) of this page.

### Tweet

To embed tweet insert `{{ tweet <id> }}`

Example:

```
{{ tweet 1224047348109795330 }}
```
Will be displayed as:

{{ tweet 1224047348109795330 }}


### Instagram

To embed Instagram post insert `{{ instagram <id> }}`

Example:

```
{{ instagram B7PFiANKX5r }}
```
Will be displayed as:

{{ instagram B7PFiANKX5r }}

### GitHub Gist
To embed GitHub gist insert `{{ gist <userid> <gistid> [file] }}`

Parameter `file` is optional and used when gist contain more than one file.

Example:

```
{{ gist gvanrossum 18bdf248a679155f1381 echo_server_tulip.py }}
```

Will be displayed as:

{{ gist gvanrossum 18bdf248a679155f1381 echo_server_tulip.py }}
