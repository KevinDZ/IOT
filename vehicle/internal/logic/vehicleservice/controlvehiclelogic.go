package vehicleservicelogic

import (
	"context"

	"vehicle/internal/svc"
	"vehicle/pb/vehicle"

	"github.com/zeromicro/go-zero/core/logx"
)

type ControlVehicleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewControlVehicleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ControlVehicleLogic {
	return &ControlVehicleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 接收云端下发的控制指令
func (l *ControlVehicleLogic) ControlVehicle(in *vehicle.ControlReq) (*vehicle.ControlResp, error) {
	// todo: add your logic here and delete this line

	return &vehicle.ControlResp{}, nil
}
