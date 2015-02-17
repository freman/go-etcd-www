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
		parent := path.Dir(r.Node.Key)
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
{{ if .Etcd.Node.Key }}{{ range paths .Etcd.Node.Key }}{{ if . }}
				<li><a href="{{ . }}">{{ . | basename }}</a></li>
{{ end }}{{ end }}{{ end }}
			</ol>
{{ if eq .Action "view"}}
	{{ template "view" . }}
{{ else if .ReadOnly }}
			<b>Access denied</b>
{{ else if eq .Action "createDirectory" }}
	{{ template "createDirectory" . }}
{{ else if eq .Action "createValue" }}
	{{ template "createValue" . }}
{{ else if eq .Action "editValue" }}
	{{ template "editValue" . }}
{{ else if eq .Action "delete" }}
	{{ template "delete" . }}
{{ end }}
		</div>

		<footer class="footer">
			<div class="container">
				<p class="text-muted"><a href="https://github.com/freman/go-etcd-www">etcd-www</a> version {{ version }} - Copyright 2015 Shannon Wynter - All rights reserved</p>
			</div>
		</footer>
	</body>
</html>
{{define "view"}}
{{ if not .ReadOnly }}
			<div>
{{ if .Etcd.Node.Dir }}				<a class="btn btn-default btn-xs" href="?a=createValue"><span class="glyphicon glyphicon-file"></span> Create Value</a>
				<a class="btn btn-default btn-xs" href="?a=createDirectory"><span class="glyphicon glyphicon-folder-close"></span> Create Directory</a>{{ else }}
				<a class="btn btn-default btn-xs" href="?a=editValue"><span class="glyphicon glyphicon-pencil"></span> Edit Value</a>{{ end }}
{{ if hasParent .Etcd }}				<a class="btn btn-danger btn-xs" href="?a=delete"><span class="glyphicon glyphicon-trash"></span> Delete</a>{{ end }}
			</div>
{{ end }}
{{ if .Etcd.Node.Key }}
			<div id="properties" style="margin-bottom: 20px;">
				<div class="row">
					<div class="col-md-3">Key<span style="float: right">=</span></div>
					<div class="col-md-9">{{ .Etcd.Node.Key }}</div>
				</div>
{{ if not .Etcd.Node.Dir }}				<div class="row">
					<div class="col-md-3">Value<span style="float: right">=</span></div>
					<div class="col-md-9">{{ .Etcd.Node.Value }}</div>
				</div>{{ end }}
				<div class="row">
					<div class="col-md-3">Expiration<span style="float: right">=</span></div>
					<div class="col-md-9">{{ .Etcd.Node.Expiration }}</div>
				</div>
				<div class="row">
					<div class="col-md-3">TTL<span style="float: right">=</span></div>
					<div class="col-md-9">{{ .Etcd.Node.TTL }}</div>
				</div>
				<div class="row">
					<div class="col-md-3">ModifiedIndex<span style="float: right">=</span></div>
					<div class="col-md-9">{{ .Etcd.Node.ModifiedIndex }}</div>
				</div>
				<div class="row">
					<div class="col-md-3">CreatedIndex<span style="float: right">=</span></div>
					<div class="col-md-9">{{ .Etcd.Node.CreatedIndex }}</div>
				</div>
			</div>
			{{ end }}
			<div id="tree" style="margin-bottom: 20px;">
				{{ if hasParent .Etcd }}<div class="row">
					<div class="col-md-6"><span class="glyphicon glyphicon-level-up"></span> <a href="{{ .Etcd.Node.Key | dir }}">..</a></div>
				</div>{{ end }}
				{{ if .Etcd.Node.Dir }}{{ range .Etcd.Node.Nodes }}
				<div class="row">
					{{ if .Dir }}{{ template "folder" . }}{{ else }}{{ template "entry" . }}{{ end }}
				</div>{{ end }}{{ end }}
			</div>
{{ end }}
{{ define "createDirectory" }}
			<form method="post">
				<div class="form-group">
					<label for="dirName">Directory name</label>
					<input type="text" class="form-control" name="dirName" id="dirName" placeholder="Enter a node/directory name">
				</div>
				<button type="submit" class="btn btn-default">Create</button>
			</form>
{{ end }}
{{ define "createValue" }}
			<form method="post">
				<div class="form-group">
					<label for="valueName">Value name</label>
					<input type="text" class="form-control" name="valueName" id="valueName" placeholder="Enter a name">
				</div>
				<div class="form-group">
					<label for="valueValue">Value Value</label>
					<input type="text" class="form-control" name="valueValue" id="valueValue" placeholder="Enter a value">
				</div>
				<button type="submit" class="btn btn-default">Create</button>
			</form>
{{ end }}
{{ define "editValue" }}
			<form method="post">
				<div class="form-group">
					<label for="valueValue">Value Value</label>
					<input type="hidden" name="oldValue" id="oldValue" value="{{ .Etcd.Node.Value }}" placeholder="Enter a value">
					<input type="text" class="form-control" name="newValue" id="newValue" value="{{ .Etcd.Node.Value }}" placeholder="Enter a value">
				</div>
				<button type="submit" class="btn btn-default">Save</button>
			</form>
{{ end }}
{{ define "delete" }}
			<form method="post">
				<p>Are you sure you want to delete the {{ if .Etcd.Node.Dir }}directory{{ else }}file{{ end }}
				"{{ .Etcd.Node.Key }}"?</p><p>This action is non-reverasble and will result in the
				{{ if .Etcd.Node.Dir }}recursive loss of all the data held below this
				directory{{ else }}lost of the contents of this key{{ end }}.</p>
				<input name="confirm" value="Cancel" type="submit" class="btn btn-default">
				<input name="confirm" value="Delete" type="submit" class="btn btn-danger">
			</form>
{{ end }}
{{ define "folder" }}<div class="col-md-6"><span class="glyphicon glyphicon-folder-close"></span> <a href="{{ .Key }}">{{ .Key | basename }}</a></div>{{ end }}
{{ define "entry" }}<div class="col-md-6"><span class="glyphicon glyphicon-file"></span> <a href="{{ .Key }}">{{ .Key | basename }}</a><span style="float: right">=</span></div>
				<div class="col-md-6">{{ .Value }}</div>{{ end }}
`
