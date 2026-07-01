// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package vehicle

import (
	"context"
	"encoding/json"

	"mobile/internal/svc"
	"mobile/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ControlVehicleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewControlVehicleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ControlVehicleLogic {
	return &ControlVehicleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ControlVehicleLogic) ControlVehicle(req *types.VehicleControlReq) (resp *types.VehicleControlResp, err error) {
	// 1. 组装 MQTT 消息体
	command := map[string]string{"action": req.Action, "vehicleID": req.VehicleID}
	payload, err := json.Marshal(command)
	if err != nil {
		return nil, err
	}

	// 2. 通过 svcCtx 获取 MQTT Client 并发布消息
	topic := "/vehicle/" + req.VehicleID + "/control"
	token := l.svcCtx.MQTTClient.Publish(topic, 1, false, payload)
	token.Wait()

	// 3. 处理发布结果
	if token.Error() != nil {
		l.Logger.Errorf("MQTT 指令下发失败: %v", token.Error())
		return &types.VehicleControlResp{
			Code: 500,
			Msg:  "指令下发失败，请稍后重试",
		}, nil
	}

	return &types.VehicleControlResp{
		Code: 200,
		Msg:  "指令下发成功",
	}, nil

}
