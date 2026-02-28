package collector

import (
	"github.com/yash492/statusy/internal/adapter/collector/atlassian/circleci"
	"github.com/yash492/statusy/internal/adapter/collector/registry"
)

func RegisterAll() map[string]registry.ProviderBuilderFunc {
	circleci.Register()
	return registry.GetAll()
}
