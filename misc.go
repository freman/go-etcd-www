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
			<meta name="viewport" content="width=device-width, initial-scale=1">
			<title>etcd-www</title>
			<link rel="stylesheet" href="//yui.yahooapis.com/pure/0.5.0/pure-min.css">
			<link rel="stylesheet" href="//maxcdn.bootstrapcdn.com/font-awesome/4.3.0/css/font-awesome.min.css">
		</head>
		<body>
			<h2>etcd</h2>
			<p><b>{{ if hasParent . }}<a href="{{ .Node.Key | dir }}">..</a>{{ end }} {{ .Node.Key | basename }}</b></p>{{ if .Node.Dir }}{{ range .Node.Nodes }}
			<div class="pure-g">
				{{ if .Dir }}{{ template "folder" . }}{{ else }}{{ template "entry" . }}{{ end }}
			</div>{{ end }}{{ end }}
			<div>
				<p><b>© 2015 Fremnet. All rights reserved.</b></p>
			</div>
		</body>
	</html>
`

const folderTemplateHTML = `
				<div class="pure-u-1-24"><i class="fa fa-folder"></i></div>
				<div class="pure-u-23-24"><a href="{{ .Key }}">{{ .Key | basename }}</a></div>
`

const entryTemplateHTML = `
				<div class="pure-u-1-24"><i class="fa fa-file"></i></div>
				<div class="pure-u-1-2"><a href="{{ .Key }}">{{ .Key | basename }}</a></div>
				<div class="pure-u-1-24">=</div>
				<div class="pure-u-10-24">{{ .Value }}</div>
`
