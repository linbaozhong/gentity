package scheduler

import (
	"context"
	"log/slog"
	"time"

	"lets/internal/config"
	"lets/internal/model"
	"lets/internal/service"
)

// Scheduler 定时续期调度器
type Scheduler struct {
	cfg     *config.Config
	service *service.CertService
	cancel  context.CancelFunc
}

func New(cfg *config.Config, svc *service.CertService) *Scheduler {
	return &Scheduler{cfg: cfg, service: svc}
}

// Start 启动定时任务
func (s *Scheduler) Start(ctx context.Context) {
	if !s.cfg.Scheduler.Enabled {
		slog.Info("scheduler disabled")
		return
	}
	ctx, s.cancel = context.WithCancel(ctx)
	interval := time.Duration(s.cfg.Scheduler.Interval) * time.Hour

	go func() {
		// 启动时先检查一次
		s.checkAll(ctx)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				slog.Info("scheduler stopped")
				return
			case <-ticker.C:
				s.checkAll(ctx)
			}
		}
	}()
	slog.Info("scheduler started", "interval_hours", s.cfg.Scheduler.Interval)
}

// Stop 停止定时任务
func (s *Scheduler) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
}

func (s *Scheduler) checkAll(ctx context.Context) {
	slog.Info("checking certificates for renewal")
	certs, err := s.service.List(ctx)
	if err != nil {
		slog.Error("list certificates failed", "error", err)
		return
	}
	for _, cert := range certs {
		if cert.Status != model.CertStatusActive {
			continue
		}
		meta, err := s.service.GetMetaForCheck(cert.Domain)
		if err != nil {
			continue
		}
		need, err := s.service.CheckRenewal(ctx, meta)
		if err != nil {
			slog.Error("check renewal failed", "domain", cert.Domain, "error", err)
			continue
		}
		if need {
			slog.Info("certificate needs renewal", "domain", cert.Domain, "days_left", cert.DaysLeft)
			if err = s.service.Renew(ctx, cert.Domain); err != nil {
				slog.Error("renew failed", "domain", cert.Domain, "error", err)
			}
		}
	}
}
