// Package pool 包提供了一个对象池实现，用于管理对象的重用。
package pool

import (
	"context"
	"sync"
)

// objPool 是对象池的结构体，它管理对象的创建、存储和过期。
type objPool struct {
	pool *sync.Pool // 底层的对象池，用于存储和管理对象。
}

// New 创建并返回一个新的对象池。fn 是一个函数，用于创建新对象。
func New(ctx context.Context, fn func() any) *objPool {
	p := &objPool{
		pool: &sync.Pool{New: fn},
	}

	return p
}

// Get 从对象池中获取一个对象。如果对象池中没有可用对象，则创建一个新的。
func (p *objPool) Get() any {
	// 从sync.Pool中获取一个对象。
	obj := p.pool.Get()
	// 如果对象为nil，则调用pool.New来创建一个新的对象。
	if obj == nil {
		return p.pool.New()
	}
	// 尝试将对象断言为types.Modeler类型。
	if m, ok := obj.(PoolModeler); ok {
		m.Get4Pool()
		return m
	}
	// 如果类型断言失败，创建一个新的对象。
	return p.pool.New()
}

// Put 将对象放回对象池中。如果对象已存在（基于UUID），则不放入。
func (p *objPool) Put(obj PoolModeler) {
	// 忽略nil对象。
	if obj == nil {
		return
	}

	if obj.Put2Pool() {
		return
	}
	// 重置对象状态，并将其放回对象池中。
	obj.Reset()
	p.pool.Put(obj)
}
