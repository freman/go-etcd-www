package main

import (
	"path"
	"html/template"
	"github.com/coreos/go-etcd/etcd"
)

const version = "0.0.1";

var templateFunctions = template.FuncMap{
	"basename"  : path.Base,
	"dir"       : path.Dir,
	"paths"     : func (full string) []string {
		components := make([]string, 0)
		for len(full) > 1 {
			components = append([]string{full}, components...)
			full = path.Dir(full)
		}
		return components
	},
	"hasParent" : func (r *etcd.Response) bool {
		parent := path.Dir(r.Node.Key);
		return parent != "."
	},
	"version"   : func () string {
		return version
	},
}

const baseTemplateHTML = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<meta name="description" content="Browsable etcd.">
		<title>etcd-www</title>
		<link rel="stylesheet" href="//maxcdn.bootstrapcdn.com/bootstrap/3.3.2/css/bootstrap.min.css">
		<link rel="stylesheet" href="//maxcdn.bootstrapcdn.com/bootstrap/3.3.2/css/bootstrap-theme.min.css">
		<style type="text/css">
			html {position: relative; min-height: 100%;}
			body {margin-bottom: 15px;}
			.footer {position: absolute; bottom: 0; width: 100%; height: 15px; background-color: #f5f5f5;}
			.footer .text-muted {font-size: x-small; margin: 0;}
			#properties .row:nth-child(odd) {background-color: #fafafa;}
		</style>
	</head>
	<body>
		<div class="container">
			<ol class="breadcrumb">
				<li><a href="/">root</a></li>
{{ if .Node.Key }}{{ range paths .Node.Key }}{{ if . }}
				<li><a href="{{ . }}">{{ . | basename }}</a></li>
{{ end }}{{ end }}{{ end }}
			</ol>
{{ if .Node.Key }}

			<div id="properties" style="margin-bottom: 20px;">
				<div class="row">
					<div class="col-md-3">Key<span style="float: right">=</span></div>
					<div class="col-md-9">{{ .Node.Key }}</div>
				</div>
				<div class="row">
					<div class="col-md-3">Value<span style="float: right">=</span></div>
					<div class="col-md-9">{{ .Node.Value }}</div>
				</div>
				<div class="row">
					<div class="col-md-3">Expiration<span style="float: right">=</span></div>
					<div class="col-md-9">{{ .Node.Expiration }}</div>
				</div>
				<div class="row">
					<div class="col-md-3">TTL<span style="float: right">=</span></div>
					<div class="col-md-9">{{ .Node.TTL }}</div>
				</div>
				<div class="row">
					<div class="col-md-3">ModifiedIndex<span style="float: right">=</span></div>
					<div class="col-md-9">{{ .Node.ModifiedIndex }}</div>
				</div>
				<div class="row">
					<div class="col-md-3">CreatedIndex<span style="float: right">=</span></div>
					<div class="col-md-9">{{ .Node.CreatedIndex }}</div>
				</div>
			</div>
			{{ end }}
			<div id="tree" style="margin-bottom: 20px;">
				{{ if hasParent . }}<div class="row">
					<div class="col-md-6"><span class="glyphicon glyphicon-level-up"></span> <a href="{{ .Node.Key | dir }}">..</a></div>
				</div>{{ end }}
				{{ if .Node.Dir }}{{ range .Node.Nodes }}
				<div class="row">
					{{ if .Dir }}{{ template "folder" . }}{{ else }}{{ template "entry" . }}{{ end }}
				</div>{{ end }}{{ end }}
			</div>
		</div>

		<footer class="footer">
			<div class="container">
				<p class="text-muted"><a href="https://github.com/freman/go-etcd-www">etcd-www</a> version {{ version }} - Copyright 2015 Shannon Wynter - All rights reserved</p>
			</div>
		</footer>
	</body>
</html>
{{define "folder"}}<div class="col-md-6"><span class="glyphicon glyphicon-folder-close"></span> <a href="{{ .Key }}">{{ .Key | basename }}</a></div>{{ end }}
{{define "entry"}}<div class="col-md-6"><span class="glyphicon glyphicon-file"></span> <a href="{{ .Key }}">{{ .Key | basename }}</a><span style="float: right">=</span></div>
				<div class="col-md-6">{{ .Value }}</div>{{end}}
`
