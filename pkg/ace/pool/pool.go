// Package pool 包提供了一个对象池实现，用于管理对象的重用。
package pool

import (
	"sync"
)

// objPool 是对象池的结构体，它管理对象的创建、存储和过期。
type objPool[T PoolModeler] struct {
	pool *sync.Pool // 底层的对象池，用于存储和管理对象。
}

// New 创建并返回一个新的对象池。fn 是一个函数，用于创建新对象。
func New[T PoolModeler](fn func() any) *objPool[T] {
	p := &objPool[T]{
		pool: &sync.Pool{
			New: func() any {
				return fn()
			},
		},
	}

	return p
}

// Get 从对象池中获取一个对象。如果对象池中没有可用对象，则创建一个新的。
func (p *objPool[T]) Get() T {
	// 从sync.Pool中获取一个对象。
	obj, ok := p.pool.Get().(T)
	if !ok {
		// 类型不匹配，创建新对象
		if creator := p.pool.New; creator != nil {
			return creator().(T)
		}
		var zero T
		return zero
	}

	// 调用对象的get4Pool方法
	obj.get4Pool()
	return obj
}

// Put 将对象放回对象池中。如果对象已存在（基于UUID），则不放入。
func (p *objPool[T]) Put(obj T) {
	var _obj = any(obj)
	// 如果是 nil 或零值，不放入池中
	if _obj == nil {
		return
	}

	// 额外检查零值（如果 T 不是指针）
	if _obj == any(*new(T)) {
		return
	}

	if obj.put2Pool() {
		return
	}
	// 重置对象状态，并将其放回对象池中。
	obj.Reset()
	p.pool.Put(obj)
}
