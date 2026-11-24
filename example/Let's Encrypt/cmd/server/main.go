package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"lets/internal/acme"
	"lets/internal/api"
	"lets/internal/cas"
	"lets/internal/config"
	"lets/internal/deploy"
	"lets/internal/dns"
	"lets/internal/scheduler"
	"lets/internal/service"
	"lets/internal/storage"
)

func main() {
	configPath := flag.String("config", "config.yaml", "config file path")
	flag.Parse()

	// 初始化日志
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	// 加载配置
	cfg, err := config.Load(*configPath)
	if err != nil {
		slog.Error("load config failed", "error", err)
		os.Exit(1)
	}
	slog.Info("config loaded", "ca", cfg.ACME.CADir, "db", cfg.Database.Path)

	// 初始化存储
	store, err := storage.New(cfg.Database.Path)
	if err != nil {
		slog.Error("init storage failed", "error", err)
		os.Exit(1)
	}
	defer store.Close()

	// 初始化 ACME 客户端
	acmeClient, err := acme.NewClient(cfg)
	if err != nil {
		slog.Error("init ACME client failed", "error", err)
		os.Exit(1)
	}

	// 初始化 CAS 客户端
	casClient, err := cas.NewClient(cfg)
	if err != nil {
		slog.Error("init CAS client failed", "error", err)
		os.Exit(1)
	}

	// 初始化 DNS 工厂
	dnsFactory := dns.NewFactory(cfg)

	// 初始化部署管理器
	deployManager, err := deploy.NewManager(cfg)
	if err != nil {
		slog.Error("init deploy manager failed", "error", err)
		os.Exit(1)
	}

	// 初始化服务
	svc := service.NewCertService(cfg, store, acmeClient, casClient, dnsFactory, deployManager)

	// 启动定时续期
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sched := scheduler.New(cfg, svc)
	sched.Start(ctx)
	defer sched.Stop()

	// 启动 HTTP 服务
	webDir := "web"
	if _, err = os.Stat(webDir + "/index.html"); os.IsNotExist(err) {
		webDir = "" // 管理界面不存在时不提供静态资源
	}
	router := api.NewRouter(api.NewHandler(svc), webDir)

	srv := &http.Server{
		Addr:    cfg.Server.Addr,
		Handler: router,
	}

	// 优雅关闭
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		slog.Info("shutting down...")
		cancel()
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		_ = srv.Shutdown(shutdownCtx)
	}()

	slog.Info("server started", "addr", cfg.Server.Addr)
	if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}
