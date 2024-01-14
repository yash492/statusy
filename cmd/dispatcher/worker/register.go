package worker

import (
	"github.com/yash492/statusy/cmd/dispatcher/webhook"
	"github.com/yash492/statusy/pkg/types"
)

func New() WorkerMap {
	//Slack
	// registerWorker(types.SlackWorker, types.IncidentTriggeredEventType, slack.IncidentOpenWorker{})
	// registerWorker(types.SlackWorker, types.IncidentInProgressEventType, slack.IncidentInProgressWorker{})
	// registerWorker(types.SlackWorker, types.IncidentResolvedEventType, slack.IncidentClosedWorker{})

	// //MSTeams
	// registerWorker(types.MsTeamsWorker, types.IncidentTriggeredEventType, msteams.IncidentOpenWorker{})
	// registerWorker(types.MsTeamsWorker, types.IncidentInProgressEventType, msteams.IncidentInProgressWorker{})
	// registerWorker(types.MsTeamsWorker, types.IncidentResolvedEventType, msteams.IncidentClosedWorker{})

	// //Discord
	// registerWorker(types.DiscordWorker, types.IncidentTriggeredEventType, discord.IncidentOpenWorker{})
	// registerWorker(types.DiscordWorker, types.IncidentInProgressEventType, discord.IncidentInProgressWorker{})
	// registerWorker(types.DiscordWorker, types.IncidentResolvedEventType, discord.IncidentClosedWorker{})

	// //Squadcast
	// registerWorker(types.SquadcastWorker, types.IncidentTriggeredEventType, squadcast.IncidentOpenWorker{})
	// registerWorker(types.SquadcastWorker, types.IncidentResolvedEventType, squadcast.IncidentClosedWorker{})

	// //Pagerduty
	// registerWorker(types.PagerdutyWorker, types.IncidentTriggeredEventType, pagerduty.IncidentOpenWorker{})
	// registerWorker(types.PagerdutyWorker, types.IncidentResolvedEventType, pagerduty.IncidentClosedWorker{})

	// Discord
	registerWorker(types.WebhookWorker, types.IncidentTriggeredEventType, webhook.IncidentOpenWorker{})
	registerWorker(types.WebhookWorker, types.IncidentInProgressEventType, webhook.IncidentInProgressWorker{})
	registerWorker(types.WebhookWorker, types.IncidentResolvedEventType, webhook.IncidentClosedWorker{})

	return dispatchWorker
}

func registerWorker(workerName string, eventType string, worker Worker) {
	eventWorkerMap, ok := dispatchWorker[workerName]
	if !ok {
		eventWorkerMap = make(map[string]Worker, 0)
	}
	eventWorkerMap[eventType] = worker
	dispatchWorker[workerName] = eventWorkerMap
}
