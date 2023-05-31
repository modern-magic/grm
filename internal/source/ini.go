package source

import (
	"fmt"
	"strings"

	"github.com/nonzzz/ini"
	"github.com/nonzzz/ini/pkg/ast"
)

type GrmIni struct {
	Path string
	ini  *ini.Ini
}

func NewGrmIniParse(conf *GrmConfig) *GrmIni {
	p := &GrmIni{}
	i, _ := ini.New().LoadFile(conf.ConfPath)
	p.ini = i
	return p
}

func (i *GrmIni) Set(k, v string) bool {
	i.ini.Walk(func(node, _ ast.Node) {
		switch t := node.(type) {
		case *ast.ExpressionNode:
			if t.Key == k {
				t.Text = fmt.Sprintf("%s = %s", k, v)
			}
		}
	})
	return true
}

func (i *GrmIni) Get(k string) (val string) {
	i.ini.Walk(func(node, _ ast.Node) {
		switch t := node.(type) {
		case *ast.ExpressionNode:
			if t.Key == k {
				val = strings.TrimSpace(t.Value)
			}
		}
	})
	i.Path = val
	return val
}

func (i *GrmIni) ToString() string {
	s, _ := i.ini.Printer()
	return s
}
