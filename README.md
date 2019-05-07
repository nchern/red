[![Go Report Card](https://goreportcard.com/badge/github.com/nchern/red)](https://goreportcard.com/report/github.com/nchern/red)

RED (Request EDitor)
==

RED is a Console User Interface(CUI) that helps to make HTTP requests to different APIs from your favourite console text editor. You can think of it like a notebook for API requests.
Alpha version, tested with Vim only ;)
In the current version empty or JSON request body is supported. 

Install
===

```go
go get github.com/nchern/red/...
```

Usage
===

```bash
red # Opens editor for editing requests in a query file

red run # runs the query from a query file

red example # Outputs self-explainable query file body. You can get an idea about query file syntax from it 
```

Screenshots
====
![Overview](https://github.com/nchern/red/blob/master/screenshots/general.png)

![Query selection example](https://github.com/nchern/red/blob/master/screenshots/selection.png)

Vim hints
===
A couple of convenient shortcuts could be added to `~/.vimrc`
```
" runs a query
command! -range=% DoQuery :<line1>,<line2>!red run
```

Other editors support
===
Exists in theory, not tested. However, there are a couple of env variables to control an editor invocation: 
 * `EDITOR` - an editor command name, defaults to `vim`
 * `EDITOR_FLAGS` - flags to pass to the editor

The utility calls then an editor with the following command line:
```bash
<EDITOR> <EDITOR_FLAGS> <query-filename> 
```


Inspired by
===
 * [Sense](https://chrome.google.com/webstore/detail/sense-beta/lhjgkmllcaadmopgmanpapmpjgmfcfig?hl=en)
 * [Postman](https://chrome.google.com/webstore/detail/postman/fhbjgbiflinjbdggehcddcbncdddomop?hl=en)
