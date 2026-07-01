package logic

import (
	"context"

	"vehicle/internal/svc"
	"vehicle/vehicle"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetVehicleStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetVehicleStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVehicleStatusLogic {
	return &GetVehicleStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 供其他服务查询车辆实时状态
func (l *GetVehicleStatusLogic) GetVehicleStatus(in *vehicle.ControlReq) (*vehicle.VehicleStatusResp, error) {
	// todo: add your logic here and delete this line

	return &vehicle.VehicleStatusResp{}, nil
}
