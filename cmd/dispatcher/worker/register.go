package worker

import (
	"github.com/yash492/statusy/cmd/dispatcher/discord"
	"github.com/yash492/statusy/cmd/dispatcher/msteams"
	"github.com/yash492/statusy/cmd/dispatcher/slack"
	"github.com/yash492/statusy/pkg/types"
)

func New() WorkerMap {
	//Slack
	registerWorker(types.SlackWorker, types.IncidentTriggeredEventType, slack.IncidentOpenWorker{})
	registerWorker(types.SlackWorker, types.IncidentInProgressEventType, slack.IncidentInProgressWorker{})
	registerWorker(types.SlackWorker, types.IncidentResolvedEventType, slack.IncidentClosedWorker{})

	//MSTeams
	registerWorker(types.MsTeamsWorker, types.IncidentTriggeredEventType, msteams.IncidentOpenWorker{})
	registerWorker(types.MsTeamsWorker, types.IncidentInProgressEventType, msteams.IncidentInProgressWorker{})
	registerWorker(types.MsTeamsWorker, types.IncidentResolvedEventType, msteams.IncidentClosedWorker{})

	//Discord
	registerWorker(types.DiscordWorker, types.IncidentTriggeredEventType, discord.IncidentOpenWorker{})
	registerWorker(types.DiscordWorker, types.IncidentInProgressEventType, discord.IncidentInProgressWorker{})
	registerWorker(types.DiscordWorker, types.IncidentResolvedEventType, discord.IncidentClosedWorker{})
	
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
