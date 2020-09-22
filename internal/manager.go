package internal

type Manager interface {
	Create(creds Creds) (*string, error)
	Remove(id Id) error
	Get(id Id) (*CredsWithId, error)
	GetAll() []CredsWithId
	ObserveIfSupported(service ManagerService)
	IsObserveSupported() bool
	Update(id Id, creds Creds) error
}
