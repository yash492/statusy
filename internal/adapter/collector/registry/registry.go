package registry

import (
	"sync"

	"github.com/yash492/statusy/internal/domain/statuspage"
)

var (
	providerBuilder map[string]statuspage.StatusPageProvider = map[string]statuspage.StatusPageProvider{}
	mu              *sync.Mutex                              = &sync.Mutex{}
)

func Register(serviceSlug string, service statuspage.StatusPageProvider) {
	mu.Lock()
	defer mu.Unlock()
	providerBuilder[serviceSlug] = service
}

func GetAll() map[string]statuspage.StatusPageProvider {
	return providerBuilder
}
