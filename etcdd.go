package main

import (
	"github.com/coreos/go-etcd/etcd"
	"log"
	"html/template"
	"net/http"
	"os"
)

var templates = template.Must(template.New("base").Funcs(templateFunctions).Parse(baseTemplateHTML))

func initTempalates() {
	templates = template.Must(templates.New("folder").Parse(folderTemplateHTML))
	templates = template.Must(templates.New("entry").Parse(entryTemplateHTML))
}

func getPort() string {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "4747"
		log.Println("INFO: No PORT environment variable set, using default " + port)
	}
	return ":" + port
}

func getEtcd() []string {
	var etcd = os.Getenv("ETCD_PEERS")
	if etcd == "" {
		etcd = "http://localhost:4001"
		log.Println("INFO: No ETCD_PEERS environment variable set, using default " + etcd)
	}

	return []string{etcd}
}

func main() {
	initTempalates();
	http.HandleFunc("/", queryEtcd)

	log.Println("Starting...")

	err := http.ListenAndServe(getPort(), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func queryEtcd (w http.ResponseWriter, r *http.Request) {
	e := etcd.NewClient(getEtcd())
	result, err := e.Get(r.URL.Path, true, false)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return;
	}

	log.Printf("%#v", result.Node.Nodes[0]);

	err = templates.ExecuteTemplate(w, "base", result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return;
	}
}
