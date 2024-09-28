package pool

import (
	"context"
	"github.com/linbaozhong/gentity/pkg/ace/types"
	"sync"
	"time"
)

type objPool struct {
	ctx context.Context

	pool *sync.Pool
	// uuid -> obj 只保留存在于pool中的对象的uuid,防止重复对象被放入
	keys *sync.Map

	expire   time.Duration
	interval time.Duration

	cleanTimer *time.Timer
}

type opt func(*objPool)

// New 创建一个对象池 ,fn 为创建新对象的函数
func New(ctx context.Context, fn func() any, opts ...opt) *objPool {
	p := &objPool{
		ctx:      ctx,
		pool:     &sync.Pool{New: fn},
		keys:     &sync.Map{},
		expire:   2 * time.Minute,
		interval: time.Minute,
	}

	for _, opt := range opts {
		opt(p)
	}

	p.cleanTimer = time.NewTimer(p.interval)

	go p.cleanup()

	return p
}

// WithExpire 设置对象在池中的生命时长
func WithExpire(d time.Duration) opt {
	return func(p *objPool) {
		p.expire = d
	}
}

// WithInterval 设置清理间隔
func WithInterval(d time.Duration) opt {
	return func(p *objPool) {
		p.interval = d
	}
}

// Get 从池中获取一个对象
func (p *objPool) Get() any {
	a := p.pool.Get()
	if m, ok := a.(types.AceModeler); ok {
		p.keys.Delete(m.UUID())
		return m
	}
	return p.pool.New()
}

// Put
func (p *objPool) Put(obj types.AceModeler) {
	uuid := obj.UUID()
	// 如果已经存在，则丢弃，防止重复放入
	if _, ok := p.keys.Load(uuid); ok {
		return
	}
	p.pool.Put(obj)
	p.keys.Store(uuid, time.Now())
}

// cleanup
func (p *objPool) cleanup() {
	for {
		select {
		case <-p.ctx.Done():
			p.cleanTimer.Stop()
			p.keys = nil
			p.pool = nil
			return
		case <-p.cleanTimer.C:
			expired := time.Now().Add(-p.expire)
			p.keys.Range(func(key, value any) bool {
				if t, ok := value.(time.Time); ok && t.Before(expired) {
					p.keys.Delete(key)
				}
				return true
			})
			p.cleanTimer.Reset(p.interval)
		}
	}
}
