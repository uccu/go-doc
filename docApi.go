package doc

import (
	"regexp"
	"strings"

	"github.com/uccu/go-stringify"
)

// @Summary	获取用户信息
// @Desc/Description	获取用户信息
// @Category 用户分类
// @Router 	路径
// @Type 	http/ws
// @Method 	POST[,GET]
// @Tag 	用户
// @Accept 	json
// @Header... 	Auth [true] 备注
// @Rest 	Struct
// @Body 	Struct
// @Success... [200] KEY Struct
// @FAIL... 	[400] KEY Struct
type DocApi struct {
	Summary     string       `json:"summary"`
	Description string       `json:"description"`
	Category    string       `json:"category"`
	Router      string       `json:"router"`
	Type        string       `json:"type"`
	Method      []string     `json:"method"`
	Tag         []string     `json:"tag"`
	Accept      []string     `json:"accept"`
	Header      []*DocHeader `json:"header"`
	Rest        *DocStruct   `json:"rest"`
	Body        *DocStruct   `json:"body"`
	Success     []*DocRet    `json:"success"`
	Fail        []*DocRet    `json:"fail"`
	pkg         *Pkg         `json:"-"`
}

type DocHeader struct {
	Name    string
	Require bool
	Remark  string
}

type DocRet struct {
	Code  int16
	Key   string
	Value *DocStruct
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
	case "Description":
		return doc.ParseDescription(commentPieces)
	case "Desc":
		return doc.ParseDescription(commentPieces)
	case "Category":
		return doc.ParseCategory(commentPieces)
	case "Router":
		return doc.ParseRouter(commentPieces)
	case "Type":
		return doc.ParseType(commentPieces)
	case "Method":
		return doc.ParseMethod(commentPieces)
	case "Tag":
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

func (doc *DocApi) parseDocStruct(s string) *DocStruct {

	strus := stringify.ToStringSlice(s, ".")
	var stru Struct
	for k, name := range strus {
		if k == 0 {
			stru = doc.pkg.GetStru(name)
			if stru == nil {
				stru = doc.pkg.GetPkg(name)
				if stru == nil {
					return nil
				}
				continue
			}
		}
		stru = stru.GetStru(name)
		if stru == nil {
			return nil
		}
	}

	if docStruct, ok := stru.(*DocStruct); ok {
		return docStruct
	}
	return nil
}

func (doc *DocApi) ParseSuccess(s []string) bool {
	if len(s) < 2 {
		return false
	}

	if len(s) == 2 {
		s[2] = s[1]
		s[1] = s[0]
		s[0] = "200"
	}

	stru := doc.parseDocStruct(s[2])
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
	if len(s) < 2 {
		return false
	}

	if len(s) == 2 {
		s[2] = s[1]
		s[1] = s[0]
		s[0] = "200"
	}

	stru := doc.parseDocStruct(s[2])
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
	doc.Body = doc.parseDocStruct(s[0])
	if doc.Body == nil {
		return false
	}
	return true
}

func (doc *DocApi) ParseRest(s []string) bool {
	if len(s) == 0 {
		return false
	}
	doc.Rest = doc.parseDocStruct(s[0])
	if doc.Rest == nil {
		return false
	}
	return true
}

func (doc *DocApi) ParseHeader(s []string) bool {
	if len(s) < 2 {
		return false
	}

	if len(s) == 2 {
		s[2] = s[1]
		s[1] = "false"
	}

	docHeader := &DocHeader{
		Name:   s[0],
		Remark: s[2],
	}

	if s[1] == "true" {
		docHeader.Require = true
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

	if doc.Accept == nil {
		doc.Accept = make([]string, 0)
	}
	for _, s := range s {
		for _, s := range stringify.ToStringSlice(s) {
			doc.Accept = append(doc.Accept, s)
		}
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
		for _, s := range stringify.ToStringSlice(s) {
			doc.Tag = append(doc.Tag, s)
		}
	}
	return true
}

func (doc *DocApi) ParseMethod(s []string) bool {
	if len(s) == 0 {
		return false
	}

	if doc.Method == nil {
		doc.Method = make([]string, 0)
	}

	for _, s := range s {
		for _, s := range stringify.ToStringSlice(s) {
			doc.Method = append(doc.Method, s)
		}
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

func NewDocApi(comments []string, pkg *Pkg) *DocApi {
	doc := &DocApi{
		Accept: []string{"json"},
		Method: []string{"post"},
		Type:   "http",
		pkg:    pkg,
	}

	for _, v := range comments {
		doc.ParseComment(v)
	}

	return doc
}
