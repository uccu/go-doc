package doc_test

import (
	"bufio"
	"encoding/json"
	"fmt"
	"go/parser"
	"go/token"
	"net/http"
	"os"
	"path"
	"runtime"
	"testing"
	"text/template"
	"time"

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

	err := ssdoc.AddPacakges("test").Export("doc")
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
	ssdoc.AddPacakges("github.com/uccu/go-doc/test/user")

	str, _ := json.Marshal(ssdoc)

	j := "http://127.0.0.1:7000/doc.json"

	http.HandleFunc("/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Write([]byte(str))
	})
	http.HandleFunc("/index.html", func(w http.ResponseWriter, r *http.Request) {

		t, _ := template.New("index.html").ParseFiles("index.html")
		t.Execute(w, j)
	})
	fmt.Println("doc addr : http://127.0.0.1:7000/index.html")
	fmt.Println("json addr : http://127.0.0.1:7000/doc.json")
	http.ListenAndServe(":7000", nil)
}

func TestTTT(t *testing.T) {

}

func file(f string) string {
	_, file, _, _ := runtime.Caller(0)
	return path.Dir(file) + "/" + f
}

func TestAst(t *testing.T) {

	fset := token.NewFileSet()
	// 解析src但在处理导入后停止。
	f, err := parser.ParseDir(fset, file("doc"), nil, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 从文件CAST打印导入。
	for {
		time.Sleep(time.Second)
		fmt.Printf("%p", f)
	}

}

func TestD(t *testing.T) {

	base, _ := os.Getwd()
	f, _ := os.Open(base + "/go.mod")
	rd := bufio.NewScanner(f)
	if rd.Scan() {
		str := []byte(rd.Text())
		js, _ := json.Marshal(string(str[7:]))
		fmt.Println(string(js))
	}
	f.Close()

}
