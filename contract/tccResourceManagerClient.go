package contract

import (
	"context"
	"github.com/cuigh/auxo/net/rpc"
)

type TCCResourceManagerClient struct {
	rpc.LazyClient
}

func (s *TCCResourceManagerClient) Confirm(ctx context.Context, req *RMConfirmTransRequest) (*RMConfirmTransResponse, error) {
	c, err := s.Try()
	if err != nil {
		return nil, err
	}

	resp := new(RMConfirmTransResponse)
	err = c.Call(ctx, "TCCResourceManagerService", "Confirm", []interface{}{req}, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *TCCResourceManagerClient) Cancel(ctx context.Context, req *RMCancelTransRequest) (*RMCancelTransResponse, error) {
	c, err := s.Try()
	if err != nil {
		return nil, err
	}

	resp := new(RMCancelTransResponse)
	err = c.Call(ctx, "TCCResourceManagerService", "Cancel", []interface{}{req}, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
