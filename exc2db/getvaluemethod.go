package exc2db

import (
	"time"
)

func NewGetValueMethod(method, format string) *GetValueMethod {
	g := new(GetValueMethod)
	g.Method = method
	g.Format = format
	return g
}

func (g *GetValueMethod) GetValue(ef *ExcelFile) (string, interface{}) {
	switch g.Method {
	case "now":
		now := time.Now()
		return now.Format(g.Format), now
	case "order":
		return "", nil
	case "translate":
		return "", nil
	default:
		return "", nil
	}
}
