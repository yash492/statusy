package services

type ProviderType string
type ServiceSlug string

type Service struct {
	ID   uint
	Name string
}

func (s ServiceSlug) String() string {
	return string(s)
}

func (s ProviderType) String() string {
	return string(s)
}
