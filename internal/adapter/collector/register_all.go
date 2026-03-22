package collector

import (
	"github.com/yash492/statusy/internal/adapter/collector/atlassian/circleci"
	"github.com/yash492/statusy/internal/adapter/collector/atlassian/claude"
	"github.com/yash492/statusy/internal/adapter/collector/registry"
	"github.com/yash492/statusy/internal/domain/statuspage"
	"resty.dev/v3"
)

func RegisterAll(client *resty.Client) map[string]statuspage.StatusPageProvider {
	claude.Register(client)
	circleci.Register(client)
	// plivo.Register(client)
	// solarwindsobservability.Register(client)
	// digitalocean.Register(client)
	// zoom.Register(client)
	// twilio.Register(client)
	// github.Register(client)
	// cloudflare.Register(client)
	// discord.Register(client)
	// dropbox.Register(client)
	// newrelic.Register(client)
	// datadog.Register(client)
	// cursor.Register(client)
	return registry.GetAll()
}
