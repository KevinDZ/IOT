package logic

import (
	"context"
	"encoding/json"
	"fmt"

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
	// 1. 组装 MQTT 消息体
	command := map[string]string{"action": in.Action}
	payload, err := json.Marshal(command)
	if err != nil {
		fmt.Println("序列化失败: ", err)
		return &vehicle.ControlResp{Code: 500, Msg: "序列化失败"}, nil
	}
	fmt.Println("接收云端下发的控制指令: ", in.VehicleID, in.Action, string(payload))
	// 2. 通过 MQTT 下发指令给真实的物理车载硬件
	topic := "/vehicle/" + in.VehicleID + "/control"
	token := l.svcCtx.MQTTClient.Publish(topic, 1, false, payload)
	token.Wait()

	if token.Error() != nil {
		return &vehicle.ControlResp{Code: 500, Msg: "MQTT下发失败"}, nil
	}

	// 3. 返回成功状态
	return &vehicle.ControlResp{Code: 200, Msg: "指令已发送至车载硬件"}, nil
}
