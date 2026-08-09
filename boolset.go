// Package boolset provides an analysis.Analyzer that finds sets
// dressed as map[T]bool.
//
// Only the unambiguous shape is reported: a map created inside a
// function — by make or a composite literal, since only there is
// every store in view — that only ever stores literal true and is
// read for membership. Storing any other value — false, a variable, a
// computed bool, whether by assignment or as a composite-literal
// element — marks the value as meaningful state and stands the
// check down, as does a read that consumes the value, such as a
// comma-ok with a used value or a range over values. Deleting,
// clearing, and taking len mean the same thing for a set, so they
// are fine. A map that arrives from elsewhere — a parameter, a
// struct field, a package variable, an aliased value — or escapes
// into a call is not examined: its stores are not all in view.
// Each variable is reported once, at its declaration.
package boolset

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "boolset",
	Doc:  "reports map[T]bool used as a set",
	Run:  run,
}

type candidate struct {
	def *ast.Ident
	bad bool
}

func run(pass *analysis.Pass) (any, error) {
	for _, f := range pass.Files {
		checkFile(pass, f)
	}
	return nil, nil
}

func checkFile(pass *analysis.Pass, f *ast.File) {
	var (
		stack []ast.Node
		cands []*candidate

		info  = pass.TypesInfo
		byObj = make(map[types.Object]*candidate)
	)
	ast.Inspect(f, func(n ast.Node) bool {
		if n == nil {
			stack = stack[:len(stack)-1]
			return true
		}
		stack = append(stack, n)
		id, ok := n.(*ast.Ident)
		if !ok {
			return true
		}
		if obj, isDef := info.Defs[id].(*types.Var); isDef {
			if !boolMap(obj.Type()) || obj.Parent() == pass.Pkg.Scope() {
				return true
			}
			// Only a declaration with its initializer in view is a
			// candidate: a parameter, result, struct field, or range
			// variable arrives from elsewhere.
			c := &candidate{def: id}
			switch parent := stack[len(stack)-2].(type) {
			case *ast.ValueSpec:
				for i, name := range parent.Names {
					if name == id && i < len(parent.Values) {
						c.bad = !newSet(info, parent.Values[i])
					}
				}
			case *ast.AssignStmt:
				for i, lhs := range parent.Lhs {
					if lhs != id {
						continue
					}
					if i >= len(parent.Rhs) || !newSet(info, parent.Rhs[i]) {
						c.bad = true
					}
				}
			default:
				return true
			}
			byObj[obj] = c
			cands = append(cands, c)
			return true
		}
		obj := info.Uses[id]
		c := byObj[obj]
		if c == nil {
			return true
		}
		classify(info, c, stack)
		return true
	})
	for _, c := range cands {
		if c.bad {
			continue
		}
		pass.Report(analysis.Diagnostic{
			Pos:     c.def.Pos(),
			End:     c.def.End(),
			Message: c.def.Name + " is a set and can be map[T]struct{}",
		})
	}
}

func boolMap(t types.Type) bool {
	m, ok := t.Underlying().(*types.Map)
	if !ok {
		return false
	}
	b, ok := m.Elem().Underlying().(*types.Basic)
	return ok && b.Kind() == types.Bool
}

// classify buckets one use of the candidate by its syntactic
// context, disqualifying the candidate where the shape proves
// meaningful state or an escape.
func classify(info *types.Info, c *candidate, stack []ast.Node) {
	switch parent := stack[len(stack)-2].(type) {
	case *ast.IndexExpr:
		classifyIndex(info, c, stack)
	case *ast.RangeStmt:
		if parent.Value != nil {
			c.bad = true // the bool values are read
		}
	case *ast.CallExpr:
		// Anything but a set-shaped builtin escapes into a call.
		fun, _ := parent.Fun.(*ast.Ident)
		switch b, _ := info.Uses[fun].(*types.Builtin); {
		case b == nil:
			c.bad = true
		case b.Name() == "len", b.Name() == "delete",
			b.Name() == "clear":
		default:
			c.bad = true
		}
	case *ast.AssignStmt:
		// On the left, a whole-map reassignment: the new source
		// must construct a set here. On the right, the map flows
		// into another variable.
		id := stack[len(stack)-1]
		for i, lhs := range parent.Lhs {
			if lhs != id {
				continue
			}
			if i >= len(parent.Rhs) || !newSet(info, parent.Rhs[i]) {
				c.bad = true
			}
			return
		}
		for _, rhs := range parent.Rhs {
			if rhs == id {
				c.bad = true
			}
		}
	default:
		c.bad = true
	}
}

// classifyIndex handles m[k] uses: a store of literal true is the
// set shape and any other store is meaningful state; a comma-ok
// read with a blank value is a membership test and one that
// consumes the value proves a real bool map. Any other read is a
// membership read and is fine.
func classifyIndex(info *types.Info, c *candidate, stack []ast.Node) {
	index := stack[len(stack)-2].(*ast.IndexExpr)
	switch container := stack[len(stack)-3].(type) {
	case *ast.ValueSpec:
		if len(container.Names) != 2 {
			return
		}
		if container.Names[0].Name != "_" {
			c.bad = true
		}
	case *ast.AssignStmt:
		assign := container
		for i, lhs := range assign.Lhs {
			if lhs != index {
				continue
			}
			if len(assign.Lhs) != 1 || !isTrue(info, assign.Rhs[i]) {
				c.bad = true
			}
			return
		}
		commaOK := len(assign.Lhs) == 2 && len(assign.Rhs) == 1 &&
			assign.Rhs[0] == index
		if commaOK {
			id, isID := assign.Lhs[0].(*ast.Ident)
			if !isID || id.Name != "_" {
				c.bad = true
			}
		}
	}
}

// isTrue reports whether e is the predeclared true.
func isTrue(info *types.Info, e ast.Expr) bool {
	id, ok := e.(*ast.Ident)
	return ok && info.Uses[id] == types.Universe.Lookup("true")
}

// newSet reports whether e constructs a set here: a call of the
// predeclared make, or a composite literal whose every element
// stores the predeclared true.
func newSet(info *types.Info, e ast.Expr) bool {
	switch node := e.(type) {
	case *ast.CallExpr:
		id, ok := node.Fun.(*ast.Ident)
		if !ok {
			return false
		}
		b, isBuiltin := info.Uses[id].(*types.Builtin)
		return isBuiltin && b.Name() == "make"
	case *ast.CompositeLit:
		for _, el := range node.Elts {
			kv, ok := el.(*ast.KeyValueExpr)
			if !ok || !isTrue(info, kv.Value) {
				return false
			}
		}
		return true
	}
	return false
}
