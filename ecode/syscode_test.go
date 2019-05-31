package ecode

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"regexp"
	"strings"
	"testing"
	"text/template"
)

const src = `//
// !!!! 此文件为生成,不要编辑
//
package ecode

var (
	OK						= add(0)		// "OK"

	SysBegin				= add(-1)		// "系统错误开始"
	ParamInvalid			= add(-2)		// "无效参数"
	ServerErr				= add(-3)		// "服务器错误"
	RequestErr				= add(-4)		// "请求错误"
)

`

const xref = `{{ .CodeSrc  }}
var err_msg_map = map[Code]string { {{ range $tag, $values := .TypeAnnotations }}
		{{ $tag }}:				"{{ $values }}",{{end}}
}
`

const file_path = "./syscode.go"

func TestGenErrMsg(t *testing.T) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		fmt.Println(fmt.Sprintf("解析文件错误: %s", err.Error()))
		return
	}

	items := make(map[string]string)

	add := func(name, cat string) {
		items[name] = cat
	}

	tag := regexp.MustCompile("\"[^\"]+\"")

	parse := func(name, comment string) {
		defs := tag.FindAllString(comment, -1)
		for _, def := range defs {
			def = strings.Trim(def, "\"")
			add(name, def)
		}
	}

	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.ValueSpec:
			for _, ident := range x.Names {
				for _, cmt := range x.Comment.List {
					fmt.Println(ident.Name, cmt.Text)
					parse(ident.Name, cmt.Text)
				}
			}
		}
		return true
	})
	fmt.Println(items)

	temp, err := template.New("").Parse(xref)
	if err != nil {
		fmt.Println(fmt.Sprintf("解析模板错误: %s", err.Error()))
		return
	}

	os.Remove(file_path)
	fd, create_err := os.Create(file_path)
	if create_err != nil {
		fmt.Println("创建错误码文件失败")
		return
	}
	defer fd.Close()

	temp.Execute(fd, map[string]interface{}{
		"PackageName":     f.Name.Name,
		"TypeAnnotations": items,
		"CodeSrc":         src,
	})
}
