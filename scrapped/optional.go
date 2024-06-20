package container

import "errors"

type Optional[V any] interface {
	Exists() bool
	Get() (*V, error)
	MustGet() V
}

type PointerOf[V any] struct {
	val *V
}

func NewEmptyOf[V any]() Optional[V] {
	return PointerOf[V]{val: nil}
}

func NewPointerOf[V any](val V) Optional[V] {
	return PointerOf[V]{val: &val}
}

func (p PointerOf[V]) Exists() bool {
	return p.val != nil
}

func (p PointerOf[V]) Get() (*V, error) {
	if !p.Exists() {
		return nil, errors.New("optional is empty. cannot get")
	}

	return p.val, nil
}

func (p PointerOf[V]) MustGet() V {
	return *p.val
}
