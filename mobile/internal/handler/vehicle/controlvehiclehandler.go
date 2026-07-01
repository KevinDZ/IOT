// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package vehicle

import (
	"fmt"
	"net/http"

	"mobile/internal/logic/vehicle"
	"mobile/internal/svc"
	"mobile/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ControlVehicleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.VehicleControlReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		fmt.Println("发起的控制指令: ", req.VehicleID, req.Action)
		l := vehicle.NewControlVehicleLogic(r.Context(), svcCtx)
		resp, err := l.ControlVehicle(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
