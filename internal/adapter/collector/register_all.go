package collector

import (
	"github.com/yash492/statusy/internal/adapter/collector/atlassian/circleci"
	"github.com/yash492/statusy/internal/adapter/collector/atlassian/cloudflare"
	"github.com/yash492/statusy/internal/adapter/collector/atlassian/datadog"
	"github.com/yash492/statusy/internal/adapter/collector/atlassian/digitalocean"
	"github.com/yash492/statusy/internal/adapter/collector/atlassian/discord"
	"github.com/yash492/statusy/internal/adapter/collector/atlassian/github"
	"github.com/yash492/statusy/internal/adapter/collector/atlassian/newrelic"
	"github.com/yash492/statusy/internal/adapter/collector/atlassian/plivo"
	"github.com/yash492/statusy/internal/adapter/collector/atlassian/solarwindsobservability"
	"github.com/yash492/statusy/internal/adapter/collector/atlassian/twilio"
	"github.com/yash492/statusy/internal/adapter/collector/atlassian/zoom"
	"github.com/yash492/statusy/internal/adapter/collector/registry"
	"github.com/yash492/statusy/internal/domain/statuspage"
)

func RegisterAll() map[string]statuspage.StatusPageProvider {
	circleci.Register()
	plivo.Register()
	solarwindsobservability.Register()
	digitalocean.Register()
	zoom.Register()
	twilio.Register()
	github.Register()
	cloudflare.Register()
	discord.Register()
	newrelic.Register()
	datadog.Register()
	return registry.GetAll()
}
