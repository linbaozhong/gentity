// Package pool 包提供了一个对象池实现，用于管理对象的重用。
package pool

import (
	"context"
	"github.com/linbaozhong/gentity/pkg/util"
	"github.com/orcaman/concurrent-map/v2"
	"strconv"
	"sync"
	"time"
)

// objPool 是对象池的结构体，它管理对象的创建、存储和过期。
type objPool struct {
	// ctx context.Context // 上下文，用于控制对象池的生命周期。

	pool *sync.Pool // 底层的对象池，用于存储和管理对象。
	// keys 用于跟踪对象的唯一标识符，防止重复对象被放入池中。
	keys cmap.ConcurrentMap[uint64, time.Time] // *sync.Map
	// expire 定义对象在池中的最大存活时间。
	expire time.Duration
	// interval 定义清理过程的运行间隔。
	interval time.Duration
	// // cleanTimer 定时器，用于定期执行清理任务。
	// cleanTimer *time.Timer
}

// opt 是一个函数，用于配置对象池。
type opt func(*objPool)

// New 创建并返回一个新的对象池。fn 是一个函数，用于创建新对象。
func New(ctx context.Context, fn func() any, opts ...opt) *objPool {
	p := &objPool{
		// ctx:      ctx,
		pool: &sync.Pool{New: fn},
		keys: cmap.NewWithCustomShardingFunction[uint64, time.Time](func(key uint64) uint32 {
			return uint32(util.MemHashString(strconv.FormatUint(key, 10)))
		}), // &sync.Map{},
		expire:   2 * time.Minute, // 默认对象过期时间为2分钟。
		interval: time.Minute,     // 默认清理间隔为1分钟。
	}

	// 应用可选配置。
	for _, opt := range opts {
		opt(p)
	}
	// 启动后台goroutine执行清理任务。
	go p.cleanup(ctx)

	return p
}

// WithExpire 返回一个opt函数，用于设置对象的过期时间。
func WithExpire(d time.Duration) opt {
	return func(p *objPool) {
		p.expire = d
	}
}

// WithInterval 返回一个opt函数，用于设置清理间隔。
func WithInterval(d time.Duration) opt {
	return func(p *objPool) {
		p.interval = d
	}
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
		// 如果对象类型正确，重置其状态，并从keys中删除对应的UUID。
		p.keys.Remove(m.UUID())
		m.Reset()
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

	uuid := obj.UUID()
	// 检查对象是否已经存在于池中。
	if _, ok := p.keys.Get(uuid); ok {
		return
	}

	// 重置对象状态，并将其放回对象池中。
	obj.Reset()
	p.pool.Put(obj)
	// 存储对象的UUID和最后使用时间。
	p.keys.Set(uuid, time.Now())
}

// cleanup 是一个定时运行的清理任务，用于删除过期对象。
func (p *objPool) cleanup(ctx context.Context) {
	// 创建定时器，用于定期清理过期对象。
	cleanTimer := time.NewTimer(p.interval)
	defer cleanTimer.Stop()

	for {
		select {
		case <-ctx.Done(): // 如果上下文被取消，退出并清理goroutine。
			p.keys.Clear()
			p.pool = nil
			return
		case <-cleanTimer.C:
			// 计算过期时间点。
			expired := time.Now().Add(-p.expire)
			p.keys.IterCb(func(key uint64, v time.Time) {

			})
			p.keys.IterCb(func(key uint64, value time.Time) {
				if value.Before(expired) {
					// 删除过期对象。
					p.keys.Remove(key)
				}
			})
			// 重置定时器。
			cleanTimer.Reset(p.interval)
		}
	}
}
