package collector

import (
	"github.com/yash492/statusy/internal/adapter/collector/atlassian/circleci"
	"github.com/yash492/statusy/internal/adapter/collector/atlassian/plivo"
	"github.com/yash492/statusy/internal/adapter/collector/registry"
	"github.com/yash492/statusy/internal/domain/statuspage"
)

func RegisterAll() map[string]statuspage.StatusPageProvider {
	circleci.Register()
	plivo.Register()
	return registry.GetAll()
}
