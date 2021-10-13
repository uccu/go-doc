package doc_test

import (
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

	err := ssdoc.Export("doc", "github.com/uccu/go-doc/test")
	if err != nil {
		t.Error(err)
	}
}
