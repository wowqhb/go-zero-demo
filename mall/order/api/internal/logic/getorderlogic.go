package logic

import (
	"context"

	"go-zero-demo/mall/order/api/internal/svc"
	"go-zero-demo/mall/order/api/internal/types"
	"go-zero-demo/mall/user/rpc/userclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderLogic {
	return &GetOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrderLogic) GetOrder(req *types.OrderReq) (resp *types.OrderReply, err error) {
	out, err := l.svcCtx.UserRpc.GetUser(l.ctx, &userclient.IdRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}

	return &types.OrderReply{
		Id:     out.Id,
		Name:   out.Name,
		Gender: out.Gender,
	}, nil
}
