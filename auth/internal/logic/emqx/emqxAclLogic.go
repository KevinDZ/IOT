// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package emqx

import (
	"context"
	"strings"

	"auth/internal/svc"
	"auth/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type EmqxAclLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEmqxAclLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmqxAclLogic {
	return &EmqxAclLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EmqxAclLogic) EmqxAcl(req *types.EmqxAclReq) (resp *types.EmqxAclResp, err error) {
	// 1. 基础校验：如果客户端 ID 不符合规范，直接拒绝
	if !strings.HasPrefix(req.ClientId, "mobile_") &&
		!strings.HasPrefix(req.ClientId, "vehicle_") &&
		!strings.HasPrefix(req.ClientId, "cloud_") {
		l.Errorf("ACL DENIED: invalid clientId format, clientId=%s", req.ClientId)
		return &types.EmqxAclResp{Result: "deny"}, nil
	}

	// 2. 业务级鉴权：提取 Topic 中的 vehicle_id
	// 假设 Topic 格式为 /vehicle/1001/status
	parts := strings.Split(req.Topic, "/")
	if len(parts) >= 3 {
		vehicleID := parts[2]

		// 如果是手机端发起请求，校验该用户是否绑定了这辆车
		if strings.HasPrefix(req.ClientId, "mobile_") {
			userID := strings.TrimPrefix(req.ClientId, "mobile_")

			// // 调用 Model 层查询车辆绑定关系
			// isBound, err := l.svcCtx.UserModel.CheckVehicleBinding(l.ctx, userID, vehicleID)
			// if err != nil {
			// 	// 数据库查询出错，记录日志并拒绝（安全优先原则）
			// 	l.Errorf("CheckVehicleBinding failed: userId=%s, vehicleId=%s, err=%v", userID, vehicleID, err)
			// 	return &types.EmqxAclResp{Result: "deny"}, nil
			// }

			// if !isBound {
			// 	// 用户试图访问未绑定的车辆，拒绝
			// 	l.Infof("ACL DENIED: userId=%s tried to access unbound vehicleId=%s", userID, vehicleID)
			// 	return &types.EmqxAclResp{Result: "deny"}, nil
			// }

			// TODO 【Mock 数据】：假设只有 user1 绑定了 vehicle_1001
			if userID == "user1" && vehicleID == "1001" {
				l.Infof("ACL ALLOWED (Mock): userId=%s accessed vehicleId=%s", userID, vehicleID)
				return &types.EmqxAclResp{Result: "allow"}, nil
			}

			// 其他情况一律视为未绑定
			l.Infof("ACL DENIED (Mock): userId=%s tried to access unbound vehicleId=%s", userID, vehicleID)
			return &types.EmqxAclResp{Result: "deny"}, nil
		}
	}

	// 3. 校验通过，放行
	return &types.EmqxAclResp{Result: "allow"}, nil
}
