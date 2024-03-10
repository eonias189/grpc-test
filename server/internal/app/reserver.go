package app

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/eonias189/grpc-test/gen/go/proto"
)

type Reserver struct {
	proto.UnimplementedReverserServer
	logger  *slog.Logger
	history []string
	wg      sync.WaitGroup
}

func (r *Reserver) logHistory(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ctx.Done():
			r.logger.Info("stop logging history")
			return
		case <-ticker.C:
			r.logger.Info(fmt.Sprintf("history: %v", r.history))
		}
	}
}

func (r *Reserver) Close() {
	r.wg.Wait()
}

func (r *Reserver) Run(ctx context.Context) {
	r.logger.Info("reserver started")
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		r.logHistory(ctx)
	}()
}

func Reverse(s string) string {
	sRunes := []rune(s)
	resRunes := make([]rune, len(sRunes))
	for i := 0; i < len(sRunes); i++ {
		resRunes[len(resRunes)-i-1] = sRunes[i]
	}
	return string(resRunes)
}

func (r *Reserver) Reverse(ctx context.Context, req *proto.Request) (*proto.Response, error) {
	r.history = append(r.history, req.GetValue())
	return &proto.Response{Res: Reverse(req.GetValue())}, nil
}

func NewReserver(logger *slog.Logger) *Reserver {
	r := &Reserver{logger: logger}
	return r
}
