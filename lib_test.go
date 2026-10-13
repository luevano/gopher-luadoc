package luadoc

import (
	"strings"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func testLib(L *lua.LState) *Lib {
	return &Lib{
		Name:        "test",
		Description: "Test library",
		Vars: []*Var{
			{Name: "pattern", Value: L.NewUserData(), Type: "test_pattern"},
			{Name: "encoding", Value: L.NewUserData(), Type: "test_encoding"},
			{Name: "other", Value: L.NewUserData(), Type: "test_other"},
			{Name: "name", Value: lua.LString("x"), Type: String},
		},
		Funcs: []*Func{
			{
				Name:  "split",
				Value: func(*lua.LState) int { return 0 },
				Params: []*Param{
					{Name: "fn", Type: Func{Returns: []*Param{{Type: String}, {Type: Number, Optional: true}}}.AsType()},
				},
			},
		},
		Classes: []*Class{
			{Name: "test_pattern", Methods: []*Method{{Name: "match", Value: func(*lua.LState) int { return 0 }}}},
		},
	}
}

func TestLuaDocIsStable(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	first := testLib(L).LuaDoc()
	for range 20 {
		if doc := testLib(L).LuaDoc(); doc != first {
			t.Fatalf("the doc changed between runs:\n%s\n---\n%s", first, doc)
		}
	}
}

func TestLuaDocAliases(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	doc := testLib(L).LuaDoc()

	for _, want := range []string{"--- @alias test_encoding userdata\n--- @alias test_other userdata"} {
		if !strings.Contains(doc, want) {
			t.Errorf("missing %q in:\n%s", want, doc)
		}
	}
	// classes and builtin types already exist
	for _, unwanted := range []string{"@alias test_pattern", "@alias string"} {
		if strings.Contains(doc, unwanted) {
			t.Errorf("unexpected %q in:\n%s", unwanted, doc)
		}
	}
}

func TestFuncAsTypeHasAllReturns(t *testing.T) {
	got := Func{
		Params:  []*Param{{Name: "a", Type: String}, {Name: "b", Type: Number, Optional: true}},
		Returns: []*Param{{Type: String}, {Type: Boolean, Optional: true}},
	}.AsType()
	if want := "fun(a: string, b: number?): string, boolean?"; got != want {
		t.Errorf("wanted %q, got %q", want, got)
	}
}

func TestLoaderInstallsRootMembers(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	sub := &Lib{Name: "sub", Vars: []*Var{{Name: "x", Value: lua.LNumber(1)}}}
	root := &Lib{
		Name:  "root",
		Vars:  []*Var{{Name: "version", Value: lua.LString("1.0")}},
		Funcs: []*Func{{Name: "hi", Value: func(L *lua.LState) int { L.Push(lua.LString("hi")); return 1 }}},
		Libs:  []*Lib{sub},
	}
	L.PreloadModule(root.Name, root.Loader())

	if err := L.DoString(`local r = require("root"); result = r.version .. r.hi() .. r.sub.x`); err != nil {
		t.Fatal(err)
	}
	if got := L.GetGlobal("result").String(); got != "1.0hi1" {
		t.Errorf("got %q", got)
	}
}

func TestLuaDocInlineDescriptions(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	lib := &Lib{
		Name:        "multi",
		Description: "First line.\nSecond line.",
		Vars:        []*Var{{Name: "v", Value: lua.LNumber(1), Description: "A var\n  over two lines."}},
		Funcs: []*Func{{
			Name:    "f",
			Value:   func(*lua.LState) int { return 0 },
			Params:  []*Param{{Name: "p", Type: String, Description: "A param\nover two lines."}},
			Returns: []*Param{{Name: "r", Type: String, Description: "A return\n\nover two lines."}},
		}},
	}
	doc := lib.LuaDoc()

	for _, want := range []string{
		"--- First line.\n--- Second line.",
		"--- @field public v number A var over two lines.\n",
		"--- @param p string A param over two lines.\n",
		"--- @return string r A return over two lines.\n",
	} {
		if !strings.Contains(doc, want) {
			t.Errorf("missing %q in:\n%s", want, doc)
		}
	}
}
