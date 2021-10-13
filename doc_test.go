package doc_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/uccu/go-doc"
)

func TestVV(t *testing.T) {

	ssdoc := doc.NewSSDoc(doc.SSDocInfo{
		Version:     "0.1.1",
		Title:       "本地接口文档标题",
		Description: "本地接口文档描述",
	}, map[doc.SSDocServerId]*doc.SSDocServer{
		"http": {
			Url:         "http://127.0.0.1:8080",
			Description: "本地接口文档",
		},
	})

	err := ssdoc.AddPacakges("github.com/uccu/go-doc/test").Export("doc")
	if err != nil {
		t.Error(err)
	}
}

func TestJson(t *testing.T) {

	ssdoc := doc.NewSSDoc(doc.SSDocInfo{
		Version:     "0.1.1",
		Title:       "本地接口文档标题",
		Description: "本地接口文档描述",
	}, map[doc.SSDocServerId]*doc.SSDocServer{
		"http": {
			Url:         "http://127.0.0.1:8080",
			Description: "本地接口文档",
		},
	})
	ssdoc.AddPacakges("github.com/uccu/go-doc/test")

	str, _ := json.Marshal(ssdoc)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Write([]byte(str))
	})
	http.ListenAndServe(":7000", nil)
}
