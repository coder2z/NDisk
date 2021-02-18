package xclient

import (
	"github.com/myxy99/component/pkg/xconsole"
	"github.com/myxy99/component/pkg/xdefer"
	"github.com/myxy99/component/xcfg"
	"google.golang.org/grpc"

	"github.com/myxy99/ndisk/pkg/constant"
	NUserPb "github.com/myxy99/ndisk/pkg/pb/nuser"
	xrpc "github.com/myxy99/ndisk/pkg/rpc"
)

var (
	userClient NUserPb.NUserServiceClient
	grpcCfg    *xrpc.GRPCConfig
)

func GetUserRpc() NUserPb.NUserServiceClient {
	if userClient == nil {
		userClient = NUserPb.NewNUserServiceClient(connection("ndisk_nuser"))
	}
	return userClient
}

func GetGRPCCfg() *xrpc.GRPCConfig {
	if grpcCfg == nil {
		grpcCfg = xcfg.UnmarshalWithExpect("rpc", xrpc.DefaultGRPCConfig()).(*xrpc.GRPCConfig)
	}
	return grpcCfg
}

func connection(servername string, op ...grpc.DialOption) *grpc.ClientConn {
	option := xrpc.DefaultClientOption(GetGRPCCfg())
	option = append(option, grpc.WithInsecure())
	option = append(option, op...)
	conn, err := grpc.Dial(constant.GRPCTargetEtcd.Format(constant.DefaultNamespaces, servername), option...)
	if err != nil {
		panic(err.Error())
	}
	xdefer.Register(func() error {
		xconsole.Redf("grpc conn close => server name:", servername)
		return conn.Close()
	})
	return conn
}
