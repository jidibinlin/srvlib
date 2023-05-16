package grpchandler

import (
	"context"
	"github.com/gzjjyz/srvlib/pb3/health"
)

type Health struct {
}

func (h *Health) Ping(_ context.Context, _ *health.PingReq) (*health.PingResp, error) {
	var resp health.PingResp
	return &resp, nil
}

var _ health.HealthServer = (*Health)(nil)
