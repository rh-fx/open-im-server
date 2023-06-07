package api

import (
	"context"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/a2r"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mcontext"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/third"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func NewThird(discov discoveryregistry.SvcDiscoveryRegistry) *Third {
	// conn, err := discov.GetConn(context.Background(), config.Config.RpcRegisterName.OpenImThirdName)
	// if err != nil {
	// panic(err)
	// }
	return &Third{discov: discov}
}

type Third struct {
	conn   *grpc.ClientConn
	discov discoveryregistry.SvcDiscoveryRegistry
}

func (o *Third) client(ctx context.Context) (third.ThirdClient, error) {
	conn, err := o.discov.GetConn(ctx, config.Config.RpcRegisterName.OpenImThirdName)
	if err != nil {
		return nil, err
	}
	return third.NewThirdClient(conn), nil
}

func (o *Third) ApplyPut(c *gin.Context) {
	a2r.Call(third.ThirdClient.ApplyPut, o.client, c)
}

func (o *Third) GetPut(c *gin.Context) {
	a2r.Call(third.ThirdClient.GetPut, o.client, c)
}

func (o *Third) ConfirmPut(c *gin.Context) {
	a2r.Call(third.ThirdClient.ConfirmPut, o.client, c)
}

func (o *Third) GetHash(c *gin.Context) {
	a2r.Call(third.ThirdClient.GetHashInfo, o.client, c)
}

func (o *Third) GetSignalInvitationInfo(c *gin.Context) {
	a2r.Call(third.ThirdClient.GetSignalInvitationInfo, o.client, c)
}

func (o *Third) GetSignalInvitationInfoStartApp(c *gin.Context) {
	a2r.Call(third.ThirdClient.GetSignalInvitationInfoStartApp, o.client, c)
}

func (o *Third) FcmUpdateToken(c *gin.Context) {
	a2r.Call(third.ThirdClient.FcmUpdateToken, o.client, c)
}

func (o *Third) SetAppBadge(c *gin.Context) {
	a2r.Call(third.ThirdClient.SetAppBadge, o.client, c)
}

func (o *Third) GetURL(c *gin.Context) {
	if c.Request.Method == http.MethodPost {
		a2r.Call(third.ThirdClient.GetUrl, o.client, c)
		return
	}
	name := c.Query("name")
	if name == "" {
		c.String(http.StatusBadRequest, "name is empty")
		return
	}
	operationID := c.Query("operationID")
	if operationID == "" {
		operationID = "auto_" + strconv.Itoa(rand.Int())
	}
	expires, _ := strconv.ParseInt(c.Query("expires"), 10, 64)
	if expires <= 0 {
		expires = 3600 * 1000
	}
	attachment, _ := strconv.ParseBool(c.Query("attachment"))
	client, err := o.client(c)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Set(constant.OperationID, operationID)
	resp, err := client.GetUrl(mcontext.SetOperationID(c, operationID), &third.GetUrlReq{Name: name, Expires: expires, Attachment: attachment})
	if err != nil {
		if errs.ErrArgs.Is(err) {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		if errs.ErrRecordNotFound.Is(err) {
			c.String(http.StatusNotFound, err.Error())
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, resp.Url)
}
