package micro

import (
	"github.com/goal-web/contracts"
	"github.com/micro/go-micro/v2"
)

type ServiceProvider struct {
	app contracts.Application
}

func (s *ServiceProvider) Register(application contracts.Application) {
	s.app = application
	application.Singleton("micro", func(config contracts.Config) micro.Service {
		var (
			microConfig = config.Get("micro").(Config)
			service     = micro.NewService(
				micro.Registry(microConfig.Registry),
				micro.Auth(microConfig.Auth),
				micro.Broker(microConfig.Broker),
				micro.Cmd(microConfig.Cmd),
				micro.Config(microConfig.Config),
				micro.Client(microConfig.Client),
				micro.Server(microConfig.Server),
				micro.Store(microConfig.Store),
				micro.Runtime(microConfig.Runtime),
				micro.Transport(microConfig.Transport),
				micro.Profile(microConfig.Profile),
				micro.Context(microConfig.Context),
				micro.HandleSignal(microConfig.Signal),
			)
			options = make([]micro.Option, 0)
		)

		for _, handler := range microConfig.BeforeStart {
			options = append(options, micro.BeforeStart(handler))
		}
		for _, handler := range microConfig.BeforeStop {
			options = append(options, micro.BeforeStop(handler))
		}

		service.Init(options...)

		return service
	})
}

func (s *ServiceProvider) Start() error {
	return s.app.Call(func(service micro.Service) error {
		return service.Run()
	})[0].(error)
}

func (s *ServiceProvider) Stop() {

}
