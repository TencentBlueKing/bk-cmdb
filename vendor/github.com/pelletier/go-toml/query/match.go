package query

import (
	"fmt"
	"reflect"

	"github.com/pelletier/go-toml"
)

// base match
type matchBase struct {
	next pathFn
}

func (f *matchBase) setNext(next pathFn) {
	f.next = next
}

// terminating functor - gathers results
type terminatingFn struct {
	// empty
}

func newTerminatingFn() *terminatingFn {
	return &terminatingFn{}
}

func (f *terminatingFn) setNext(next pathFn) {
	// do nothing
}

func (f *terminatingFn) call(node interface{}, ctx *queryContext) {
	ctx.result.appendResult(node, ctx.lastPosition)
}

// match single key
type matchKeyFn struct {
	matchBase
	Name string
}

func newMatchKeyFn(name string) *matchKeyFn {
	return &matchKeyFn{Name: name}
}

func (f *matchKeyFn) call(node interface{}, ctx *queryContext) {
	if array, ok := node.([]*toml.Tree); ok {
		for _, tree := range array {
			item := tree.GetPath([]string{f.Name})
			if item != nil {
				ctx.lastPosition = tree.GetPositionPath([]string{f.Name})
				f.next.call(item, ctx)
			}
		}
	} else if tree, ok := node.(*toml.Tree); ok {
		item := tree.GetPath([]string{f.Name})
		if item != nil {
			ctx.lastPosition = tree.GetPositionPath([]string{f.Name})
			f.next.call(item, ctx)
		}
	}
}

// match single index
type matchIndexFn struct {
	matchBase
	Idx int
}

func newMatchIndexFn(idx int) *matchIndexFn {
	return &matchIndexFn{Idx: idx}
}

func (f *matchIndexFn) call(node interface{}, ctx *queryContext) {
	v := reflect.ValueOf(node)
	if v.Kind() == reflect.Slice {
		if v.Len() == 0 {
			return
		}

		// Manage negative values
		idx := f.Idx
		if idx < 0 {
			idx += v.Len()
		}
		if 0 <= idx && idx < v.Len() {
			callNextIndexSlice(f.next, node, ctx, v.Index(idx).Interface())
		}
	}
}

func callNextIndexSlice(next pathFn, node interface{}, ctx *queryContext, value interface{}) {
	if treesArray, ok := node.([]*toml.Tree); ok {
		ctx.lastPosition = treesArray[0].Position()
	}
	next.call(value, ctx)
}

// filter by slicing
type matchSliceFn struct {
	matchBase
	Start, End, Step *int
}

func newMatchSliceFn() *matchSliceFn {
	return &matchSliceFn{}
}

func (f *matchSliceFn) setStart(start int) *matchSliceFn {
	f.Start = &start
	return f
}

func (f *matchSliceFn) setEnd(end int) *matchSliceFn {
	f.End = &end
	return f
}

func (f *matchSliceFn) setStep(step int) *matchSliceFn {
	f.Step = &step
	return f
}

func (f *matchSliceFn) call(node interface{}, ctx *queryContext) {
	v := reflect.ValueOf(node)
	if v.Kind() == reflect.Slice {
		if v.Len() == 0 {
			return
		}

		var start, end, step int

		// Initialize step
		if f.Step != nil {
			step = *f.Step
		} else {
			step = 1
		}

		// Initialize start
		if f.Start != nil {
			start = *f.Start
			// Manage negative values
			if start < 0 {
				start += v.Len()
			}
			// Manage out of range values
			start = max(start, 0)
			start = min(start, v.Len()-1)
		} else if step > 0 {
			start = 0
		} else {
			start = v.Len() - 1
		}

		// Initialize end
		if f.End != nil {
			end = *f.End
			// Manage negative values
			if end < 0 {
				end += v.Len()
			}
			// Manage out of range values
			end = max(end, -1)
			end = min(end, v.Len())
		} else if step > 0 {
			end = v.Len()
		} else {
			end = -1
		}

		// Loop on values
		if step > 0 {
			for idx := start; idx < end; idx += step {
				callNextIndexSlice(f.next, node, ctx, v.Index(idx).Interface())
			}
		} else {
			for idx := start; idx > end; idx += step {
				callNextIndexSlice(f.next, node, ctx, v.Index(idx).Interface())
			}
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// match anything
type matchAnyFn struct {
	matchBase
}

func newMatchAnyFn() *matchAnyFn {
	return &matchAnyFn{}
}

func (f *matchAnyFn) call(node interface{}, ctx *queryContext) {
	if tree, ok := node.(*toml.Tree); ok {
		for _, k := range tree.Keys() {
			v := tree.GetPath([]string{k})
			ctx.lastPosition = tree.GetPositionPath([]string{k})
			f.next.call(v, ctx)
		}
	}
}

// filter through union
type matchUnionFn struct {
	Union []pathFn
}

func (f *matchUnionFn) setNext(next pathFn) {
	for _, fn := range f.Union {
		fn.setNext(next)
	}
}

func (f *matchUnionFn) call(node interface{}, ctx *queryContext) {
	for _, fn := range f.Union {
		fn.call(node, ctx)
	}
}

// match every single last node in the tree
type matchRecursiveFn struct {
	matchBase
}

func newMatchRecursiveFn() *matchRecursiveFn {
	return &matchRecursiveFn{}
}

func (f *matchRecursiveFn) call(node interface{}, ctx *queryContext) {
	originalPosition := ctx.lastPosition
	if tree, ok := node.(*toml.Tree); ok {
		var visit func(tree *toml.Tree)
		visit = func(tree *toml.Tree) {
			for _, k := range tree.Keys() {
				v := tree.GetPath([]string{k})
				ctx.lastPosition = tree.GetPositionPath([]string{k})
				f.next.call(v, ctx)
				switch node := v.(type) {
				case *toml.Tree:
					visit(node)
				case []*toml.Tree:
					for _, subtree := range node {
						visit(subtree)
					}
				}
			}
		}
		ctx.lastPosition = originalPosition
		f.next.call(tree, ctx)
		visit(tree)
	}
}

// match based on an externally provided functional filter
type matchFilterFn struct {
	matchBase
	Pos  toml.Position
	Name string
}

func newMatchFilterFn(name string, pos toml.Position) *matchFilterFn {
	return &matchFilterFn{Name: name, Pos: pos}
}

func (f *matchFilterFn) call(node interface{}, ctx *queryContext) {
	fn, ok := (*ctx.filters)[f.Name]
	if !ok {
		panic(fmt.Sprintf("%s: query context does not have filter '%s'",
			f.Pos.String(), f.Name))
	}
	switch castNode := node.(type) {
	case *toml.Tree:
		for _, k := range castNode.Keys() {
			v := castNode.GetPath([]string{k})
			if fn(v) {
				ctx.lastPosition = castNode.GetPositionPath([]string{k})
				f.next.call(v, ctx)
			}
		}
	case []*toml.Tree:
		for _, v := range castNode {
			if fn(v) {
				if len(castNode) > 0 {
					ctx.lastPosition = castNode[0].Position()
				}
				f.next.call(v, ctx)
			}
		}
	case []interface{}:
		for _, v := range castNode {
			if fn(v) {
				f.next.call(v, ctx)
			}
		}
	}
}
