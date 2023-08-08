package source

import (
	"github.com/nonzzz/ini"
)

type GrmIni struct {
	Path     string
	ini      *ini.Ini
	selector ini.Selector
}

func NewGrmIniParse(conf *GrmConfig) *GrmIni {
	p := &GrmIni{}
	i, _ := ini.New().LoadFile(conf.ConfPath)
	p.ini = i
	p.selector = ini.NewSelector(p.ini)
	return p
}

func (i *GrmIni) Set(k, v string) bool {
	op := i.selector.Query(k, ini.ExpressionKind)
	_, err := op.Get()
	if err != nil {
		return false
	}
	return op.Set(ini.AttributeBindings{
		Text: v,
	})
}

func (i *GrmIni) Get(k string) {
	op := i.selector.Query(k, ini.ExpressionKind)
	expr, err := op.Get()
	if err != nil {
		return
	}
	i.Path = expr.Attribute().Value
}

func (i *GrmIni) ToString() string {
	s, _ := i.ini.Printer()
	return s
}
