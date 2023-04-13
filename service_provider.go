package micro

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/logs"
	"go-micro.dev/v4"
)

type ServiceProvider struct {
	app             contracts.Application
	ServiceRegister func(service micro.Service) error
	service         micro.Service
}

func (provider *ServiceProvider) Register(application contracts.Application) {
	provider.app = application
	application.Singleton("micro", func(config contracts.Config) micro.Service {
		var (
			microConfig = config.Get("micro").(Config)
			service     = micro.NewService()
			options     = append(microConfig.CustomOptions, micro.HandleSignal(microConfig.Signal))
		)
		if microConfig.Registry != nil {
			options = append(options, micro.Registry(microConfig.Registry))
		}

		if microConfig.Auth != nil {
			options = append(options, micro.Auth(microConfig.Auth))
		}

		if microConfig.Broker != nil {
			options = append(options, micro.Broker(microConfig.Broker))
		}

		if microConfig.Cmd != nil {
			options = append(options, micro.Cmd(microConfig.Cmd))
		}

		if microConfig.Config != nil {
			options = append(options, micro.Config(microConfig.Config))
		}

		if microConfig.Client != nil {
			options = append(options, micro.Client(microConfig.Client))
		}

		if microConfig.Server != nil {
			options = append(options, micro.Server(microConfig.Server))
		}

		if microConfig.Store != nil {
			options = append(options, micro.Store(microConfig.Store))
		}

		if microConfig.Client != nil {
			options = append(options, micro.Client(microConfig.Client))
		}

		if microConfig.Runtime != nil {
			options = append(options, micro.Runtime(microConfig.Runtime))
		}

		if microConfig.Transport != nil {
			options = append(options, micro.Transport(microConfig.Transport))
		}

		if microConfig.Profile != nil {
			options = append(options, micro.Profile(microConfig.Profile))
		}
		if microConfig.Context != nil {
			options = append(options, micro.Context(microConfig.Context))
		}

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

func (provider *ServiceProvider) Start() (err error) {
	provider.app.Call(func(service micro.Service) {
		provider.service = service
		if err = provider.ServiceRegister(service); err != nil {
			return
		}

		err = service.Server().Start()
	})

	if err != nil {
		logs.WithError(err).Error("micro.ServiceProvider.Start: micro server start failed")
		go func() { provider.app.Stop() }()
	}

	return err
}

func (provider *ServiceProvider) Stop() {
	if provider.service != nil {
		err := provider.service.Server().Stop()
		if err != nil {
			logs.WithError(err).Error("micro service closed")
		}
	}
}
