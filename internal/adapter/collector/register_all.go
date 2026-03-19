package collector

import (
	"github.com/yash492/statusy/internal/adapter/collector/atlassian/circleci"
	"github.com/yash492/statusy/internal/adapter/collector/atlassian/claude"
	"github.com/yash492/statusy/internal/adapter/collector/atlassian/cloudflare"
	"github.com/yash492/statusy/internal/adapter/collector/atlassian/cursor"
	"github.com/yash492/statusy/internal/adapter/collector/atlassian/datadog"
	"github.com/yash492/statusy/internal/adapter/collector/atlassian/digitalocean"
	"github.com/yash492/statusy/internal/adapter/collector/atlassian/discord"
	"github.com/yash492/statusy/internal/adapter/collector/atlassian/dropbox"
	"github.com/yash492/statusy/internal/adapter/collector/atlassian/github"
	"github.com/yash492/statusy/internal/adapter/collector/atlassian/newrelic"
	"github.com/yash492/statusy/internal/adapter/collector/atlassian/plivo"
	"github.com/yash492/statusy/internal/adapter/collector/atlassian/solarwindsobservability"
	"github.com/yash492/statusy/internal/adapter/collector/atlassian/twilio"
	"github.com/yash492/statusy/internal/adapter/collector/atlassian/zoom"
	"github.com/yash492/statusy/internal/adapter/collector/registry"
	"github.com/yash492/statusy/internal/domain/statuspage"
	"resty.dev/v3"
)

func RegisterAll(client *resty.Client) map[string]statuspage.StatusPageProvider {
	claude.Register(client)
	circleci.Register(client)
	plivo.Register(client)
	solarwindsobservability.Register(client)
	digitalocean.Register(client)
	zoom.Register(client)
	twilio.Register(client)
	github.Register(client)
	cloudflare.Register(client)
	discord.Register(client)
	dropbox.Register(client)
	newrelic.Register(client)
	datadog.Register(client)
	cursor.Register(client)
	return registry.GetAll()
}
