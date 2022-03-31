package internal

type Colors struct {
	Log    string
	Danger string
	Tip    string
}

type ColorName int

const (
	LogColor ColorName = iota
	DangerColor
	TipColor
)

func (cs Colors) Color(name ColorName) string {
	switch name {
	case LogColor:
		return cs.Log
	case DangerColor:
		return cs.Danger
	case TipColor:
		return cs.Tip
	}
	return ""
}

var AnsiColor = Colors{
	Log:    "\x1b[37m%s",
	Danger: "\x1b[31m%s",
	Tip:    "\x1b[36m%s",
}
