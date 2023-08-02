package luadoc

import (
	lua "github.com/yuin/gopher-lua"
)

func newTable(L *lua.LState, fields map[string]lua.LValue, funcs map[string]lua.LGFunction) *lua.LTable {
	t := L.NewTable()
	for k, v := range fields {
		L.SetField(t, k, v)
	}

	if funcs != nil {
		L.SetFuncs(t, funcs)
	}

	return t
}
