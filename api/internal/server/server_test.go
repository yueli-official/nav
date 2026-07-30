package server_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"

	"github.com/yueli-official/nav/api/internal/runtime"
	"github.com/yueli-official/nav/api/internal/server"
)

func TestReadinessSuccess(t *testing.T) {
	status, body, _ := readinessResponse(t, func(context.Context) error { return nil })
	if status != 200 || gjson.New(body).Get("status").String() != "ready" {
		t.Fatalf("status=%d body=%s", status, body)
	}
}

func TestReadinessFailure(t *testing.T) {
	status, body, _ := readinessResponse(t, func(context.Context) error { return errors.New("database unavailable") })
	if status != 503 || gjson.New(body).Get("code").String() != "common.not_ready" {
		t.Fatalf("status=%d body=%s", status, body)
	}
}

func TestReadinessTimeout(t *testing.T) {
	status, body, elapsed := readinessResponse(t, func(ctx context.Context) error {
		<-ctx.Done()
		return ctx.Err()
	})
	if status != 503 || elapsed < 2500*time.Millisecond || elapsed > 4500*time.Millisecond {
		t.Fatalf("status=%d elapsed=%s body=%s", status, elapsed, body)
	}
}

func readinessResponse(t *testing.T, check runtime.ReadinessCheck) (int, string, time.Duration) {
	t.Helper()
	httpServer := g.Server(t.Name())
	httpServer.SetAddr("127.0.0.1:0")
	httpServer.SetDumpRouterMap(false)
	server.Configure(httpServer, server.Deps{ReadyChecks: map[string]runtime.ReadinessCheck{"database": check}})
	httpServer.Start()
	t.Cleanup(func() { _ = httpServer.Shutdown() })

	client := g.Client()
	client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", httpServer.GetListenedPort()))
	started := time.Now()
	response, err := client.Get(context.Background(), "/readyz")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Close()
	return response.StatusCode, response.ReadAllString(), time.Since(started)
}
