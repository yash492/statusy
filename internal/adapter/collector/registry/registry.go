package registry

import (
	"sync"

	"github.com/yash492/statusy/internal/domain/services"
	"github.com/yash492/statusy/internal/domain/statuspage"
)

type ProviderBuilderFunc func(serviceResult services.ServiceResult) statuspage.StatusPageProvider

var (
	providerBuilder map[string]ProviderBuilderFunc = map[string]ProviderBuilderFunc{}
	mu              *sync.Mutex                    = &sync.Mutex{}
)

func Register(serviceSlug string, serviceInitFunc ProviderBuilderFunc) {
	mu.Lock()
	defer mu.Unlock()
	providerBuilder[serviceSlug] = serviceInitFunc
}

func GetAll() map[string]ProviderBuilderFunc {
	return providerBuilder
}
