package doc

import (
	"go/ast"
	"regexp"
	"strings"

	"github.com/uccu/go-stringify"
)

// @Summary				名字
// @Desc/Description	描述
// @Category			分类
// @Router				路径
// @Type				http/ws
// @Method				post[,get]
// @Tag					tag[,tag]
// @Accept				json
// @Header...			KEY required 备注
// @Rest				Struct
// @Body				Struct
// @Success...			code KEY Struct
// @FAIL...				code KEY Struct
type DocApi struct {
	Summary     string           `json:"summary"`
	Description string           `json:"description"`
	Category    string           `json:"category"`
	Router      string           `json:"router"`
	Type        string           `json:"type"`
	Server      string           `json:"server"`
	Method      []string         `json:"method,omitempty"`
	Tag         []string         `json:"tag,omitempty"`
	Accept      []string         `json:"accept,omitempty"`
	Header      []*DocHeader     `json:"header,omitempty"`
	Rest        *TypeSpecWithKey `json:"rest,omitempty"`
	Body        *TypeSpecWithKey `json:"body,omitempty"`
	Success     []*DocRet        `json:"success,omitempty"`
	Fail        []*DocRet        `json:"fail,omitempty"`
	pkg         *Pkg             `json:"-"`
	file        string           `json:"-"`
}

type DocHeader struct {
	Name     string
	Required bool
	Remark   string
}

type DocRet struct {
	Code  int16
	Key   string
	Value *TypeSpecWithKey
}

func (doc *DocApi) ParseComment(comment string) bool {

	comment = strings.ReplaceAll(comment, "\t", " ")
	r, _ := regexp.Compile("// *@([a-zA-Z]+) *")

	types := r.FindSubmatch([]byte(comment))
	if len(types) < 2 {
		return false
	}

	typ := string(types[1])
	comment = strings.ReplaceAll(comment, string(types[0]), "")
	r, _ = regexp.Compile(" +")
	comment = strings.Trim(r.ReplaceAllString(comment, " "), " ")
	commentPieces := stringify.ToStringSlice(comment, " ")

	switch typ {
	case "Summary":
		return doc.ParseSummary(commentPieces)
	case "Desc", "Description":
		return doc.ParseDescription(commentPieces)
	case "Category":
		return doc.ParseCategory(commentPieces)
	case "Router":
		return doc.ParseRouter(commentPieces)
	case "Type":
		return doc.ParseType(commentPieces)
	case "Server":
		return doc.ParseServer(commentPieces)
	case "Method":
		return doc.ParseMethod(commentPieces)
	case "Tag", "Tags":
		return doc.ParseTag(commentPieces)
	case "Accept":
		return doc.ParseAccept(commentPieces)
	case "Header":
		return doc.ParseHeader(commentPieces)
	case "Rest":
		return doc.ParseRest(commentPieces)
	case "Body":
		return doc.ParseBody(commentPieces)
	case "Success":
		return doc.ParseSuccess(commentPieces)
	case "Fail":
		return doc.ParseFail(commentPieces)
	}

	return false
}

func (doc *DocApi) ParseSuccess(s []string) bool {
	if len(s) < 3 {
		return false
	}

	stru := parseTypeType(s[2], doc.pkg, doc.file)
	if stru == nil {
		return false
	}

	docRet := &DocRet{
		Code:  int16(stringify.ToInt(s[0])),
		Key:   s[1],
		Value: stru,
	}

	if doc.Success == nil {
		doc.Success = make([]*DocRet, 0)
	}
	doc.Success = append(doc.Success, docRet)
	return true
}

func (doc *DocApi) ParseFail(s []string) bool {
	if len(s) < 3 {
		return false
	}

	stru := parseTypeType(s[2], doc.pkg, doc.file)
	if stru == nil {
		return false
	}

	docRet := &DocRet{
		Code:  int16(stringify.ToInt(s[0])),
		Key:   s[1],
		Value: stru,
	}

	if doc.Fail == nil {
		doc.Fail = make([]*DocRet, 0)
	}
	doc.Fail = append(doc.Fail, docRet)
	return true
}

func (doc *DocApi) ParseBody(s []string) bool {
	if len(s) == 0 {
		return false
	}
	doc.Body = parseTypeType(s[0], doc.pkg, doc.file)
	return doc.Body != nil
}

func (doc *DocApi) ParseRest(s []string) bool {
	if len(s) == 0 {
		return false
	}
	doc.Rest = parseTypeType(s[0], doc.pkg, doc.file)
	return doc.Rest != nil
}

func (doc *DocApi) ParseHeader(s []string) bool {
	if len(s) < 3 {
		return false
	}

	docHeader := &DocHeader{
		Name:   s[0],
		Remark: s[2],
	}

	if s[1] == "true" {
		docHeader.Required = true
	}

	if doc.Header == nil {
		doc.Header = make([]*DocHeader, 0)
	}
	doc.Header = append(doc.Header, docHeader)

	return true
}

func (doc *DocApi) ParseAccept(s []string) bool {
	if len(s) == 0 {
		return false
	}

	doc.Accept = make([]string, 0)
	for _, s := range s {
		doc.Accept = append(doc.Accept, stringify.ToStringSlice(s)...)
	}

	return true
}

func (doc *DocApi) ParseTag(s []string) bool {
	if len(s) == 0 {
		return false
	}

	if doc.Tag == nil {
		doc.Tag = make([]string, 0)
	}

	for _, s := range s {
		doc.Tag = append(doc.Tag, stringify.ToStringSlice(s)...)
	}
	return true
}

func (doc *DocApi) ParseMethod(s []string) bool {
	if len(s) == 0 {
		return false
	}

	doc.Method = make([]string, 0)
	for _, s := range s {
		doc.Method = append(doc.Method, stringify.ToStringSlice(s)...)
	}
	return true
}

func (doc *DocApi) ParseRouter(s []string) bool {
	if len(s) > 0 {
		doc.Router = s[0]
		return true
	}
	return false
}

func (doc *DocApi) ParseType(s []string) bool {
	if len(s) > 0 {
		doc.Type = s[0]
		return true
	}
	return false
}
func (doc *DocApi) ParseServer(s []string) bool {
	if len(s) > 0 {
		doc.Server = s[0]
		return true
	}
	return false
}

func (doc *DocApi) ParseSummary(s []string) bool {
	if len(s) > 0 {
		doc.Summary = s[0]
		return true
	}
	return false
}

func (doc *DocApi) ParseDescription(s []string) bool {
	if len(s) > 0 {
		doc.Description = s[0]
		return true
	}
	return false
}

func (doc *DocApi) ParseCategory(s []string) bool {
	if len(s) > 0 {
		doc.Category = s[0]
		return true
	}
	return false
}

func NewDocApi(comments *ast.CommentGroup, pkg *Pkg, file string) *DocApi {
	if comments == nil {
		return nil
	}

	doc := &DocApi{
		Accept: []string{"json"},
		Method: []string{"post"},
		Type:   "http",
		pkg:    pkg,
		file:   file,
	}

	for _, v := range comments.List {
		doc.ParseComment(v.Text)
	}
	return doc
}
