package domain

import "github.com/yash492/statusy/pkg/store"

func New() {
	Incident = store.NewIncidentDBConn()
	Component = store.NewComponentDBConn()
	Service = store.NewServiceDBConn()
	Subscription = store.NewSubscriptionConn()
	SquadcastExtension = store.NewSquadcastExtensionConn()
	PagerdutyExtension = store.NewPagerdutyExtensionConn()
	ChatopsExtension = store.NewChatOpsExtensionConn()
	WebhookExtension = store.NewWebhookExtensionConn()
}
