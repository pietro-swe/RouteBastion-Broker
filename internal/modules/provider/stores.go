package provider

type ProvidersStore interface {
	Create(provider *Provider) error
	GetByID(id string) (*Provider, error)
	GetByName(name string) (*Provider, error)
	GetAllAvailable() ([]*Provider, error)
	List() ([]*Provider, error)
	Update(provider *Provider) error
	Delete(id string) error
}
