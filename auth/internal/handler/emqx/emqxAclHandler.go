// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package emqx

import (
	"net/http"

	"auth/internal/logic/emqx"
	"auth/internal/svc"
	"auth/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func EmqxAclHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.EmqxAclReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := emqx.NewEmqxAclLogic(r.Context(), svcCtx)
		resp, err := l.EmqxAcl(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
