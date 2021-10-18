package doc

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"
	"runtime"
	"text/template"
)

type doc struct {
	ssdoc *SSDoc
	j     []byte
	def   string
}

func (d *doc) Json(w http.ResponseWriter) *doc {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Write(d.j)
	return d
}

func (d *doc) Html(w http.ResponseWriter) *doc {
	_, file, _, ok := runtime.Caller(1)
	if ok {
		index, _ := ioutil.ReadFile(path.Dir(file) + "/index.html")
		t, _ := template.New("index.html").ParseFiles(string(index))
		t.Execute(w, d.def)
	}
	return d
}

func New(i SSDocInfo, m map[SSDocServerId]*SSDocServer, pkgs []string, def string) *doc {
	doc := &doc{
		ssdoc: NewSSDoc(i, m),
		j:     make([]byte, 0),
		def:   def,
	}
	for _, v := range pkgs {
		doc.ssdoc.AddPacakges(v)
	}
	doc.j, _ = json.Marshal(doc.ssdoc)
	return doc
}
