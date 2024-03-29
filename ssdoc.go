package doc

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"
)

type SSDoc struct {
	Version string                          `json:"version"` // ssdoc api版本
	Info    SSDocInfo                       `json:"info"`    // 文档信息
	Servers map[SSDocServerId]*SSDocServer  `json:"servers"` // 服务信息
	Apis    map[SSDocCategoryId][]*SSDocApi `json:"apis"`    // 接口信息
}

type SSDocCategoryId string
type SSDocServerId string

type SSDocInfo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

type SSDocServer struct {
	Url         string `json:"url"`                   // 服务地址前缀
	Description string `json:"description,omitempty"` // 描述
}

type SSDocApi struct {
	Name        string            `json:"name"`                  // 名称
	Description string            `json:"description,omitempty"` // 描述
	Path        string            `json:"path"`                  // 路径
	Method      []string          `json:"method,omitempty"`      // 请求方式
	Type        string            `json:"type"`                  // http/ws
	Category    SSDocCategoryId   `json:"category"`              // 分类
	Server      SSDocServerId     `json:"server,omitempty"`      // 指定服务
	Tag         []string          `json:"tag,omitempty"`         // 标签
	Accept      []string          `json:"accept,omitempty"`      // 请求返回的类型,json/xml等
	Header      []*SSDocHeader    `json:"header,omitempty"`      // 请求头
	Rest        *SSDocTypeWithKey `json:"rest,omitempty"`        // REST参数
	Body        *SSDocTypeWithKey `json:"body,omitempty"`        // 请求体参数
	Success     []*SSDocRet       `json:"success,omitempty"`     // 成功返回内容
	Fail        []*SSDocRet       `json:"fail,omitempty"`        // 失败返回内容
}

type SSDocHeader struct {
	Name        string `json:"name"`                  // 名称
	Description string `json:"description,omitempty"` // 描述
	Required    bool   `json:"required,omitempty"`    // 是否必须
}

type SSDocTypeWithKey struct {
	Key      string  `json:"key"`
	Default  *string `json:"default,omitempty"`  // 默认值
	Json     *string `json:"json,omitempty"`     // json key
	Required bool    `json:"required,omitempty"` // 是否必须
	*SSDocType
}

type SSDocType struct {
	Name        string              `json:"name"`                  // 类型名字
	Description string              `json:"description,omitempty"` // 描述
	Type        Type                `json:"type"`                  // 类型
	TypeName    string              `json:"typeName"`              // 类型名字
	Value       []*SSDocTypeWithKey `json:"value,omitempty"`       // 值
}

type SSDocRet struct {
	Code  int16             `json:"code"`
	Key   string            `json:"key"`
	Value *SSDocTypeWithKey `json:"value"`
}

func NewSSDoc(info SSDocInfo, servers map[SSDocServerId]*SSDocServer) *SSDoc {

	ssdoc := &SSDoc{
		Info:    info,
		Servers: servers,
		Apis:    make(map[SSDocCategoryId][]*SSDocApi),
	}

	_, file, _, ok := runtime.Caller(1)
	if ok {
		version, _ := ioutil.ReadFile(path.Dir(file) + "/version")
		ssdoc.Version = string(version)
	}

	return ssdoc
}

func (doc *SSDoc) AddApi(i *DocApi) *SSDoc {

	api := &SSDocApi{
		Name:        i.Summary,
		Description: i.Description,
		Path:        i.Router,
		Category:    SSDocCategoryId(i.Category),
		Method:      i.Method,
		Type:        i.Type,
		Server:      SSDocServerId(i.Server),
		Tag:         i.Tag,
		Accept:      i.Accept,
	}

	if api.Method == nil {
		api.Method = []string{"get"}
	}

	if api.Category == "" {
		api.Category = "default"
	}

	if i.Header != nil {
		api.Header = make([]*SSDocHeader, 0)
		for _, h := range i.Header {
			header := &SSDocHeader{
				Name:        h.Name,
				Description: h.Remark,
				Required:    h.Required,
			}
			api.Header = append(api.Header, header)
		}
	}

	if i.Rest != nil {
		api.Rest = parseType(i.Rest)
	}

	if i.Body != nil {
		api.Body = parseType(i.Body)
	}
	if i.Success != nil {
		api.Success = make([]*SSDocRet, 0)
		for _, r := range i.Success {
			ret := &SSDocRet{
				Code:  r.Code,
				Key:   r.Key,
				Value: parseType(r.Value),
			}
			api.Success = append(api.Success, ret)
		}
	}
	if i.Fail != nil {
		api.Fail = make([]*SSDocRet, 0)
		for _, r := range i.Fail {
			ret := &SSDocRet{
				Code:  r.Code,
				Key:   r.Key,
				Value: parseType(r.Value),
			}
			api.Fail = append(api.Fail, ret)
		}
	}

	_, ok := doc.Apis[api.Category]
	if !ok {
		doc.Apis[api.Category] = make([]*SSDocApi, 0)
	}

	doc.Apis[api.Category] = append(doc.Apis[api.Category], api)

	return doc
}

func (doc *SSDoc) AddPacakges(pacakges ...string) *SSDoc {
	apis := GetApis(pacakges...)
	for _, api := range apis {
		doc.AddApi(api)
	}
	return doc
}

func (doc *SSDoc) Export(dir string) error {
	dir = strings.TrimRight(dir, "/\\")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	js, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	dir += "/doc.json"

	os.WriteFile(dir, js, os.ModePerm)
	return nil
}

func parseType(t *TypeSpecWithKey) *SSDocTypeWithKey {

	typ := &SSDocTypeWithKey{
		Key: t.Key,
		SSDocType: &SSDocType{
			Name:        t.Name,
			Type:        t.Type,
			TypeName:    t.TypeName,
			Description: t.Comment,
		},
	}

	if t.Doc != nil {
		typ.Description = t.Doc[0]
	}

	if t.Value != nil {
		typ.Value = make([]*SSDocTypeWithKey, 0)
		for _, t := range t.Value {
			a := parseType(t)
			if t.Tags != nil {
				if tag, _ := t.Tags.Get("binding"); tag != nil {
					opt := append(tag.Options, tag.Name)
					for _, o := range opt {
						if o == "required" {
							a.Required = true
						}
					}
				}
				if tag, _ := t.Tags.Get("default"); tag != nil {
					a.Default = &tag.Name
				}
				if tag, _ := t.Tags.Get("json"); tag != nil {
					if tag.Name == "-" {
						continue
					}
					a.Json = &tag.Name
				}
			}
			typ.Value = append(typ.Value, a)
		}
	} else if t.Type == TypeType {
		v := parseTypeType(t.TypeName, t.pkg, t.file)
		if v != nil {
			typ.Value = []*SSDocTypeWithKey{parseType(v)}
		} else {
			typ.Type = CustomType
		}

	}
	return typ
}
