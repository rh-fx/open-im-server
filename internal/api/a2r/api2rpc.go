package a2r

import (
	"OpenIM/internal/apiresp"
	"OpenIM/pkg/common/log"
	"OpenIM/pkg/errs"
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func Call[A, B, C any](
	rpc func(client C, ctx context.Context, req *A, options ...grpc.CallOption) (*B, error),
	client func() (C, error),
	c *gin.Context,
) {
	var req A
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrArgs.Wrap(err.Error())) // 参数错误
		return
	}
	if check, ok := any(&req).(interface{ Check() error }); ok {
		if err := check.Check(); err != nil {
			apiresp.GinError(c, errs.ErrArgs.Wrap(err.Error())) // 参数校验失败
			return
		}
	}
	cli, err := client()
	if err != nil {
		apiresp.GinError(c, errs.ErrInternalServer.Wrap(err.Error())) // 获取RPC连接失败
		log.Error("0", "get rpc client conn err:", err.Error())
		return
	}
	data, err := rpc(cli, c, &req)
	if err != nil {
		log.Error("0", "rpc call err:", err.Error())
		apiresp.GinError(c, err) // RPC调用失败
		return
	}
	apiresp.GinSuccess(c, data) // 成功
}
