package main

import (
	"path"
	"html/template"
	"github.com/coreos/go-etcd/etcd"
)

var templateFunctions = template.FuncMap{
	"basename" : path.Base,
	"dir" : path.Dir,
	"hasParent" : func (r *etcd.Response) bool {
		parent := path.Dir(r.Node.Key);
		return parent != "."
	},
}

const baseTemplateHTML = `
<!DOCTYPE html>
	<html>
		<head>
			<meta charset="utf-8">
			<title>Results</title>
			<link rel="stylesheet" href="http://yui.yahooapis.com/pure/0.5.0/pure-min.css">
		</head>
		<body>
			<h2>etcd</h2>
			<p><b>{{ if hasParent . }}<a href="{{ .Node.Key | dir }}">..</a>{{ end }} {{ .Node.Key | basename }}</b></p>
{{ if .Node.Dir }}{{ range .Node.Nodes }}
			<p><a href="{{ .Key }}">{{ .Key | basename }}</a> </p>{{ end }}{{ end }}
			<div>
				<p><b>© 2015 Fremnet. All rights reserved.</b></p>
			</div>
		</body>
	</html>
`

