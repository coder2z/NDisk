package authority

import (
	"github.com/BurntSushi/toml"
	xapp "github.com/myxy99/component"
	"github.com/myxy99/component/pkg/xconsole"
	"github.com/myxy99/component/pkg/xdefer"
	"github.com/myxy99/component/pkg/xflag"
	"github.com/myxy99/component/pkg/xvalidator"
	"github.com/myxy99/component/xcfg"
	"github.com/myxy99/component/xcfg/datasource/manager"
	"github.com/myxy99/component/xgovern"
	"github.com/myxy99/component/xinvoker"
	"github.com/myxy99/component/xmonitor"
	myValidator "github.com/myxy99/ndisk/internal/authority/validator"
	"github.com/myxy99/ndisk/pkg/rpc"
	"google.golang.org/grpc"
	"net"
	"sync"
)

type Server struct {
	err error
	*sync.WaitGroup
}

func (s *Server) PrepareRun(stopCh <-chan struct{}) (err error) {
	s.initCfg()
	s.invoker()
	s.debug()
	s.initValidator()
	s.govern()
	return s.err
}

func (s *Server) debug() {
	xconsole.ResetDebug(xapp.Debug())
	xapp.PrintVersion()
}

func (s *Server) Run(stopCh <-chan struct{}) (err error) {
	go func() {
		<-stopCh
		s.Add(1)
		xdefer.Clean()
		s.Done()
	}()
	var (
		rpcCfg *rpc.Config
		lis    net.Listener
	)
	rpcCfg = xcfg.UnmarshalWithExpect("rpcError", rpc.DefaultConfig()).(*rpc.Config)
	s.err = rpc.DefaultRegistryEtcd(rpcCfg)
	if s.err != nil {
		return
	}
	lis, s.err = net.Listen("tcp", rpcCfg.Addr())
	if s.err != nil {
		return
	}
	serve := grpc.NewServer(rpc.DefaultOption(rpcCfg)...)

	xdefer.Register(func() error {
		serve.Stop()
		xconsole.Red("grpc server shutdown success ")
		return nil
	})
	//NFilePb.RegisterNFileServiceServer(serve, new(rpcError.Server))
	xconsole.Greenf("grpc server start up success:", rpcCfg.Addr())
	s.err = serve.Serve(lis)
	s.Wait()
	return s.err
}

func (s *Server) initCfg() {
	if s.err != nil {
		return
	}
	var data xcfg.DataSource
	data, s.err = manager.NewDataSource(xflag.NString("run", "xcfg"))
	if s.err != nil {
		return
	}
	s.err = xcfg.LoadFromDataSource(data, toml.Unmarshal)
}

func (s *Server) invoker() {
	if s.err != nil {
		return
	}
	xdefer.Register(func() error {
		return xinvoker.Close()
	})
	xinvoker.Register(
		//xgorm.Register("mysql"),
		//xredis.Register("redis"),
	)
	s.err = xinvoker.Init()
}

func (s *Server) initValidator() {
	if s.err != nil {
		return
	}
	s.err = xvalidator.Init(xcfg.GetString("server.locale"), myValidator.RegisterValidation)
}

func (s *Server) govern() {
	if s.err != nil {
		return
	}
	xmonitor.Run()
	go xgovern.Run()
}
