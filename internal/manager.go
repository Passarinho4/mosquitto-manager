package internal

type Manager interface {
	Create(lp LoginPasswordAcls) (*string, error)
	Remove(id Id) error
	Get(id Id) (*LoginPasswordAcls, error)
	GetAll() []LoginPasswordAcls
	ObserveIfSupported(service ManagerService)
	IsObserveSupported() bool
	Update(id Id, lp LoginPasswordAcls) error
}
