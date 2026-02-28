package collector

import (
	"github.com/yash492/statusy/internal/adapter/collector/atlassian/circleci"
	"github.com/yash492/statusy/internal/adapter/collector/atlassian/plivo"
	"github.com/yash492/statusy/internal/adapter/collector/registry"
)

func RegisterAll() map[string]registry.ProviderBuilderFunc {
	circleci.Register()
	plivo.Register()
	return registry.GetAll()
}
