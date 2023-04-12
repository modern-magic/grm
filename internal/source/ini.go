package source

import (
	"github.com/nonzzz/ini"
	"github.com/nonzzz/ini/pkg/ast"
)

type GrmIniVisitor struct {
	key string
	val string
	ini.IniVisitor
}

func (v *GrmIniVisitor) Expression(node *ast.Expression) {
	if node.Key == v.key {
		if v.val == "" {
			v.val = node.Value
			return
		}
		node.Value = v.val
	}
}

type GrmIni struct {
	Path string
	ini  *ini.Ini
}

func NewGrmIniParse(conf *GrmConfig) *GrmIni {
	p := &GrmIni{
		ini: ini.New(),
	}
	p.ini.LoadFile(conf.confPath)
	return p
}

func (i *GrmIni) Set(k, v string) bool {
	i.ini.Accept(&GrmIniVisitor{
		key: k,
		val: v,
	})
	return i.ini.Err() != nil
}

func (i *GrmIni) Get(k string) string {
	v := &GrmIniVisitor{
		key: k,
	}
	i.ini.Accept(v)
	i.Path = v.val
	return v.val
}

func (i *GrmIni) ToString() string {
	return i.ini.String()
}
