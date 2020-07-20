package internal

type Manager interface {
	Create(lp LoginPasswordAcls) error
	Remove(login Login) error
	GetAll() []LoginPasswordAcls
	ObserveIfSupported(service ManagerService)
	IsObserveSupported() bool
}
