# gopher-luadoc
[![Go Reference](https://pkg.go.dev/badge/github.com/luevano/gopher-luadoc.svg)](https://pkg.go.dev/github.com/luevano/gopher-luadoc)

Self-documenting libraries for [gopher-lua](https://github.com/yuin/gopher-lua): describe the
functions, variables and classes of a library once, then load it into a Lua state and generate its
[LuaLS](https://luals.github.io/) annotations (`@class`, `@param`, `@return`, ...) for editor completion.

> **Warning**
>
> The API is not stable and may change at any time.

**Note**: This is a fork of [mangalorg/gopher-luadoc](https://github.com/mangalorg/gopher-luadoc).
Used by [luaprovider](https://github.com/luevano/luaprovider) and [mangal](https://github.com/luevano/mangal).

## Usage

```go
package main

import (
	"fmt"
	"strings"

	luadoc "github.com/luevano/gopher-luadoc"
	lua "github.com/yuin/gopher-lua"
)

func main() {
	text := &luadoc.Lib{
		Name:        "text",
		Description: "Text helpers.",
		Funcs: []*luadoc.Func{{
			Name:        "upper",
			Description: "Uppercase a string.",
			Value: func(L *lua.LState) int {
				L.Push(lua.LString(strings.ToUpper(L.CheckString(1))))
				return 1
			},
			Params:  []*luadoc.Param{{Name: "s", Type: luadoc.String}},
			Returns: []*luadoc.Param{{Name: "upper", Type: luadoc.String}},
		}},
	}
	// the root lib can have its own members too: require("mylib").text.upper("hi")
	lib := &luadoc.Lib{Name: "mylib", Libs: []*luadoc.Lib{text}}

	L := lua.NewState()
	defer L.Close()
	L.PreloadModule(lib.Name, lib.Loader())

	// the annotations, e.g. to save as mylib.lua
	fmt.Println(lib.LuaDoc())
}
```
