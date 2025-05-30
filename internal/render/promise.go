package render

import (
	"fmt"

	"github.com/xcnt/go-asyncapi/internal/common"
	"github.com/dave/jennifer/jen"
	"github.com/samber/lo"
)

func NewPromise[T any](ref string, origin common.PromiseOrigin) *Promise[T] {
	return &Promise[T]{ref: ref, origin: origin}
}

func NewCbPromise[T any](findCb func(item common.Renderer, path []string) bool, origin common.PromiseOrigin) *Promise[T] {
	return &Promise[T]{findCb: findCb, origin: origin}
}

type Promise[T any] struct {
	AssignErrorNote string // Optional error message additional note to be shown when assignment fails
	ref             string
	origin          common.PromiseOrigin
	findCb          func(item common.Renderer, path []string) bool

	target   T
	assigned bool
}

func (r *Promise[T]) Assign(obj any) {
	t, ok := obj.(T)
	if !ok {
		panic(fmt.Sprintf("Cannot assign an object %+v to a promise of type %T. %s", obj, r.target, r.AssignErrorNote))
	}
	r.target = t
	r.assigned = true
}

func (r *Promise[T]) Assigned() bool {
	return r.assigned
}

func (r *Promise[T]) FindCallback() func(item common.Renderer, path []string) bool {
	return r.findCb
}

func (r *Promise[T]) Target() T {
	return r.target
}

func (r *Promise[T]) Ref() string {
	return r.ref
}

func (r *Promise[T]) Origin() common.PromiseOrigin {
	return r.origin
}

func (r *Promise[T]) WrappedGolangType() (common.GolangType, bool) {
	if !r.assigned {
		return nil, false
	}
	v, ok := any(r.target).(common.GolangType)
	return v, ok
}

func (r *Promise[T]) IsPointer() bool {
	if !r.assigned {
		return false
	}
	if v, ok := any(r.target).(golangPointerType); ok {
		return v.IsPointer()
	}
	return false
}

func (r *Promise[T]) IsStruct() bool {
	if !r.assigned {
		return false
	}
	if v, ok := any(r.target).(golangStructType); ok {
		return v.IsStruct()
	}
	return false
}

// List links can only be PromiseOriginInternal, no way to set a callback in spec
func NewListCbPromise[T any](findCb func(item common.Renderer, path []string) bool) *ListPromise[T] {
	return &ListPromise[T]{findCb: findCb}
}

type ListPromise[T any] struct {
	AssignErrorNote string // Optional error message additional note to be shown when assignment fails
	findCb          func(item common.Renderer, path []string) bool

	targets  []T
	assigned bool
}

func (r *ListPromise[T]) AssignList(objs []any) {
	var ok bool
	r.targets, ok = lo.FromAnySlice[T](objs)
	if !ok {
		panic(fmt.Sprintf("Cannot assign slice of %+v to a promise of type %T. %s", objs, r.targets, r.AssignErrorNote))
	}
	r.assigned = true
}

func (r *ListPromise[T]) Assigned() bool {
	return r.assigned
}

func (r *ListPromise[T]) FindCallback() func(item common.Renderer, path []string) bool {
	return r.findCb
}

func (r *ListPromise[T]) Targets() []T {
	return r.targets
}

func NewRendererPromise(ref string, origin common.PromiseOrigin) *RendererPromise {
	return &RendererPromise{
		Promise: *NewPromise[common.Renderer](ref, origin),
	}
}

type RendererPromise struct {
	Promise[common.Renderer]
	// DirectRender marks the promise to be rendered directly, even if object it points to not marked to do so.
	// Be careful, in order to avoid duplicated object appearing in the output, this flag should be set only for
	// objects which are not marked to be rendered directly
	DirectRender bool
}

func (r *RendererPromise) RenderDefinition(ctx *common.RenderContext) []*jen.Statement {
	return r.target.RenderDefinition(ctx)
}

func (r *RendererPromise) RenderUsage(ctx *common.RenderContext) []*jen.Statement {
	return r.target.RenderUsage(ctx)
}

func (r *RendererPromise) DirectRendering() bool {
	return r.DirectRender // Prevent rendering the object we're point to for several times
}

func (r *RendererPromise) ID() string {
	if r.Assigned() {
		return r.target.ID()
	}
	return ""
}

func (r *RendererPromise) String() string {
	return "RendererPromise -> " + r.ref
}

func NewGolangTypePromise(ref string, origin common.PromiseOrigin) *GolangTypePromise {
	return &GolangTypePromise{
		Promise: *NewPromise[common.GolangType](ref, origin),
	}
}

type GolangTypePromise struct {
	Promise[common.GolangType]
	// DirectRender marks the promise to be rendered directly, even if object it points to not marked to do so.
	// Be careful, in order to avoid duplicated object appearing in the output, this flag should be set only for
	// objects which are not marked to be rendered directly
	DirectRender bool
}

func (r *GolangTypePromise) TypeName() string {
	return r.target.TypeName()
}

func (r *GolangTypePromise) DirectRendering() bool {
	return r.DirectRender // Prevent rendering the object we're point to for several times
}

func (r *GolangTypePromise) RenderDefinition(ctx *common.RenderContext) []*jen.Statement {
	return r.target.RenderDefinition(ctx)
}

func (r *GolangTypePromise) RenderUsage(ctx *common.RenderContext) []*jen.Statement {
	return r.target.RenderUsage(ctx)
}

func (r *GolangTypePromise) ID() string {
	return "GolangTypePromise"
}

func (r *GolangTypePromise) String() string {
	return "GolangTypePromise -> " + r.ref
}
