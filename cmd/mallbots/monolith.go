package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/irononet/mallbots/internal/config"
	"github.com/irononet/mallbots/internal/monolith"
	"github.com/irononet/mallbots/internal/waiter"
)

type app struct {
	cfg     config.AppConfig
	db      *sql.DB
	logger  zerolog.Logger
	modules []monolith.Module
	mux     *chi.Mux
	rpc     *grpc.Server
	waiter  waiter.Waiter
	nc *nats.Conn
	js nats.JetStreamContext
}

func (a *app) Config() config.AppConfig {
	return a.cfg
}

func (a *app) DB() *sql.DB {
	return a.db
}

func (a *app) Logger() zerolog.Logger {
	return a.logger
}

func (a *app) Mux() *chi.Mux {
	return a.mux
}

func (a *app) RPC() *grpc.Server {
	return a.rpc
}

func (a *app) JS() nats.JetStreamContext{
	return a.js
}

func (a *app) Waiter() waiter.Waiter {
	return a.waiter
}

func (a *app) StartupModules() error {
	for _, module := range a.modules {
		if err := module.Startup(a.Waiter().Context(), a); err != nil {
			return err
		}
	}
	return nil
}

func (a *app) waitForWeb(ctx context.Context) error {
	webServer := http.Server{
		Addr:    a.cfg.Web.Address(),
		Handler: a.mux,
	}

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		fmt.Println("web server started")
		defer fmt.Println("web server shutdown")
		if err := webServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})
	group.Go(func() error {
		<-gCtx.Done()
		fmt.Println("web server to be shutdown")
		ctx, cancel := context.WithTimeout(context.Background(), a.cfg.ShutdownTimeout)
		defer cancel()
		if err := webServer.Shutdown(ctx); err != nil {
			return err
		}
		return nil
	})

	return group.Wait()
}

func (a *app) waitForRpc(ctx context.Context) error {
	listener, err := net.Listen("tcp", a.cfg.Rpc.Address())
	if err != nil {
		return err
	}

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		fmt.Println("rpc server started")
		defer fmt.Println("rpc server shutdown")
		if err := a.RPC().Serve(listener); err != nil && err != grpc.ErrServerStopped {
			return err
		}
		return nil
	})
	group.Go(func() error {
		<-gCtx.Done()
		fmt.Println("rpc server to be shutdown")
		stopped := make(chan struct{})
		go func() {
			a.RPC().GracefulStop()
			close(stopped)
		}()
		timeout := time.NewTimer(a.cfg.ShutdownTimeout)
		select {
		case <-timeout.C:
			a.RPC().Stop()
			return fmt.Errorf("rpc server failed to stop gracefully")
		case <-stopped:
			return nil
		}
	})

	return group.Wait()
}

func (a *app) waitForStream(ctx context.Context) error{
	closed := make(chan struct{})
	a.nc.SetClosedHandler(func(*nats.Conn){
		close(closed)
	})
	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error{
		fmt.Println("message stream started")
		defer fmt.Println("message stream stopped")
		<-closed
		return nil
	})
	group.Go(func() error{
		<-gCtx.Done()
		return a.nc.Drain()
	})
	return group.Wait()
}
