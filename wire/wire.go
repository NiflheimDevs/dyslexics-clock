//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/NiflheimDevs/dyslexics-clock/bootstrap"
	"github.com/NiflheimDevs/dyslexics-clock/internal/application/service"
	serviceimpl "github.com/NiflheimDevs/dyslexics-clock/internal/application/service/implement"
	"github.com/NiflheimDevs/dyslexics-clock/internal/delivery/handler"
	midauth "github.com/NiflheimDevs/dyslexics-clock/internal/delivery/middleware/authentication"
	"github.com/NiflheimDevs/dyslexics-clock/internal/delivery/middleware/panicwall"
	domainmessagebroker "github.com/NiflheimDevs/dyslexics-clock/internal/domain/message_broker" // New Import
	domainpkg "github.com/NiflheimDevs/dyslexics-clock/internal/domain/pkg"
	repository "github.com/NiflheimDevs/dyslexics-clock/internal/domain/repository"
	"github.com/NiflheimDevs/dyslexics-clock/internal/infra/database/driver"
	repositoryimpl "github.com/NiflheimDevs/dyslexics-clock/internal/infra/database/postgres"
	messagebrokerdriver "github.com/NiflheimDevs/dyslexics-clock/internal/infra/message_broker/driver" // New Import
	mosquittoimpl "github.com/NiflheimDevs/dyslexics-clock/internal/infra/message_broker/mosquitto" // New Import
	"github.com/NiflheimDevs/dyslexics-clock/internal/pkg"
	"github.com/google/wire"
)

var DatabaseProviderSet = wire.NewSet(
	driver.ConnectSQL,
)

var RepositoryProviderSet = wire.NewSet(
	repositoryimpl.NewDeviceRepo,
	repositoryimpl.NewAlarmRepo,

	wire.Bind(new(repository.AlarmRepo), new(*repositoryimpl.AlarmRepo)),
	wire.Bind(new(repository.DeviceRepo), new(*repositoryimpl.DeviceRepo)),
)

var PkgProviderSet = wire.NewSet(
	pkg.NewValidatorWrapper,
	pkg.NewSecretSauce,

	wire.Bind(new(domainpkg.SecretSauce), new(*pkg.SecretSauce)),
	wire.Bind(new(domainpkg.Validator), new(*pkg.ValidatorWrapper)),
)
var ServiceProviderSet = wire.NewSet(
	serviceimpl.NewDeviceService,
	serviceimpl.NewAlarmService,
	serviceimpl.NewJWT,
	serviceimpl.NewDeviceEventService, // New Service

	wire.Bind(new(service.DeviceService), new(*serviceimpl.DeviceService)),
	wire.Bind(new(service.AlarmService), new(*serviceimpl.AlarmService)),
	wire.Bind(new(service.JWT), new(*serviceimpl.JWT)),
	wire.Bind(new(service.DeviceEventService), new(*serviceimpl.DeviceEventService)), // New Service Binding
)

var MessageBrokerProviderSet = wire.NewSet( // New Provider Set
	messagebrokerdriver.ConnectMosquitto, // Provide MQTT Client
	mosquittoimpl.NewRequestPublisher,
	mosquittoimpl.NewDeviceEventSubscriber,
	mosquittoimpl.NewAlarmEventPublisher, // This was already present, but good to group
	mosquittoimpl.NewTimePublisher, // This was already present, but good to group

	wire.Bind(new(domainmessagebroker.RequestPublisher), new(*mosquittoimpl.RequestPublisher)),
	wire.Bind(new(domainmessagebroker.DeviceEventSubscriber), new(*mosquittoimpl.DeviceEventSubscriber)),
	wire.Bind(new(domainmessagebroker.AlarmEventPublisher), new(*mosquittoimpl.AlarmEventPublisher)),
	wire.Bind(new(domainmessagebroker.TimePublisher), new(*mosquittoimpl.TimePublisher)),
)

var HandlerProviderSet = wire.NewSet(
	handler.NewAlarmHandler,
	handler.NewDeviceHandler,
	wire.Struct(new(Handlers), "*"),
)

var MiddlewareProviderSet = wire.NewSet(
	panicwall.NewPanicWall,
	midauth.NewAuth,
	wire.Struct(new(Middlewares), "*"),
)

var ProviderSet = wire.NewSet(
	PkgProviderSet,
	ProvideConstants,
	ProvideEnv,
	DatabaseProviderSet,
	RepositoryProviderSet,
	ServiceProviderSet,
	MessageBrokerProviderSet, // Add new provider set
	HandlerProviderSet,
	MiddlewareProviderSet,
)

func ProvideConstants(container *bootstrap.Di) *bootstrap.Constants {
	return container.Const
}

func ProvideEnv(container *bootstrap.Di) *bootstrap.Env {
	return container.Env
}

type Handlers struct {
	AlarmHandler  *handler.AlarmHandler
	DeviceHandler *handler.DeviceHandler
}

type Middlewares struct {
	PanicWall *panicwall.PanicWall
	Auth      *midauth.Authentication
}

type App struct {
	Handlers    *Handlers
	Middlewares *Middlewares
	DeviceEventSubscriber domainmessagebroker.DeviceEventSubscriber
}

func InitApp(di *bootstrap.Di) (*App, error) {
	wire.Build(
		ProviderSet,
		wire.Struct(new(App), "*"),
	)
	return &App{}, nil
}
