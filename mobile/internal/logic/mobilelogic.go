// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"mobile/internal/svc"
	"mobile/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MobileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMobileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MobileLogic {
	return &MobileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MobileLogic) Mobile(req *types.VehicleControlReq) (resp *types.VehicleControlResp, err error) {
	// todo: add your logic here and delete this line

	return
}
