package doc

import (
	"encoding/json"
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

type DocConf struct {
	SSDocInfo SSDocInfo
	Server    map[SSDocServerId]*SSDocServer
	Pkgs      []string
	Url       string
	Name      string
}

func (d *doc) Json(w http.ResponseWriter) *doc {
	w.Write(d.j)
	return d
}

func (d *doc) Html(w http.ResponseWriter) *doc {
	_, file, _, ok := runtime.Caller(0)
	if ok {
		t, _ := template.New("index.html").ParseFiles(path.Dir(file) + "/index.html")
		t.Execute(w, d.def)
	}
	return d
}

func New(c DocConf) *doc {
	doc := &doc{
		ssdoc: NewSSDoc(c.SSDocInfo, c.Server),
		j:     make([]byte, 0),
		def:   c.Url,
	}
	for _, v := range c.Pkgs {
		doc.ssdoc.AddPacakges(c.Name, v)
	}
	doc.j, _ = json.Marshal(doc.ssdoc)
	return doc
}
