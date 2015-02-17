package main

import (
	"os"
	"log"
	"flag"
	"strings"
	"net/http"
	"html/template"
	"github.com/coreos/go-etcd/etcd"
)

var etcdPeers = newFlagStrs([]string{"http://localhost:4001"})
var listen = ":4747"

var templates = template.Must(template.New("base").Funcs(templateFunctions).Parse(baseTemplateHTML))
var etcdClient *etcd.Client

type requestData struct {
	Etcd *etcd.Response
	Action string
	ReadOnly bool
}

func httpHandle (w http.ResponseWriter, r *http.Request) {
	result, err := etcdClient.Get(r.URL.Path, true, false)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	action := "view"
	switch r.URL.Query().Get("a") {
		case "createDirectory":
			action = "createDirectory"
		case "createValue":
			action = "createValue"
		case "editValue":
			action = "editValue"
		case "delete":
			action = "delete"
	}

	err = templates.ExecuteTemplate(w, "base", requestData{result, action, false})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func init() {
	if envpeers := os.Getenv("ETCD_PEERS"); envpeers != "" {
		etcdPeers = newFlagStrs(strings.Split(envpeers, ","))
	}

	if envlisten := os.Getenv("LISTEN"); envlisten != "" {
		listen = envlisten
	}

	flag.Var(etcdPeers, "etcd-peer", "etcd peers, repeat to list more than one, alternatively env ETCD_PEERS")
	flag.StringVar(&listen, "port", listen, "port to listen on")
	flag.Parse()

	etcdClient = etcd.NewClient(etcdPeers.Values)
}

func main() {
	http.HandleFunc("/", httpHandle)

	err := http.ListenAndServe(listen, nil)
	if err != nil {
		log.Fatal(err)
	}
}
