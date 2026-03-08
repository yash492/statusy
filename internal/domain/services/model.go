package services

type ServiceSlug string

type Service struct {
	ID   uint
	Name string
}

func (s ServiceSlug) String() string {
	return string(s)
}
