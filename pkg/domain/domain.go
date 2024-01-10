package domain

import "github.com/yash492/statusy/pkg/store"

var Incident store.IncidentStore
var Component store.ComponentStore
var Service store.ServiceStore
var Subscription store.SubscriptionStore
var ChatopsExtension store.ChatopsExtensionStore
var SquadcastExtension store.SquadcastExtensionStore
var PagerdutyExtension store.PagerdutyExtensionStore
var WebhookExtension store.WebhookExtensionStore
