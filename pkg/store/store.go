package store

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/yash492/statusy/pkg/schema"
)

type ServiceStore interface {
	Create(services []schema.Service) error
	GetAll() ([]schema.Service, error)
}

type ComponentStore interface {
	Create(components []schema.Component) ([]schema.Component, error)
	GetAllByServiceID(serviceID uint) ([]schema.Component, error)
}

type IncidentStore interface {
	Create([]schema.Incident) ([]schema.Incident, error)
	GetByProviderID(providerID string) (schema.Incident, error)
	CreateIncidentUpdates(incidentUpdates []schema.IncidentUpdate) ([]schema.IncidentUpdate, error)
	GetLastIncidentUpdatesTimeByService(serviceID uint, incidentIDs []uint) ([]schema.LastIncidentUpdateForIncident, error)
	CreateIncidentComponents(incidentComponents []schema.IncidentComponent) error
}

type SubscriptionStore interface {
	GetAllServicesForSubscriptions(serviceName string) ([]schema.ServicesForSubsciptions, error)
	Create(serviceID uint, componentIDs []uint, isAllComponents bool) error
	GetWithComponents(subscriptionID uuid.UUID) ([]schema.SubscriptionWithComponents, error)
	Update(subscriptionID uuid.UUID, componentIDs []uint, isAllComponents bool) error
	GetForIncidentUpdates(incidentUpdateID uint) ([]schema.SubscriptionForIncidentUpdates, error)
}

type SquadcastExtensionStore interface {
	Get() (schema.SquadcastExtension, error)
	Save(webhookURL string, uuid uuid.UUID) error
	Delete(uuid uuid.UUID) error
}

type PagerdutyExtensionStore interface {
	Get() (schema.PagerdutyExtension, error)
	Save(routingKey string, uuid uuid.UUID) error
	Delete(uuid uuid.UUID) error
}

type ChatopsExtensionStore interface {
	Get() ([]schema.ChatopsExtension, error)
	GetByType(chatopType string) (schema.ChatopsExtension, error)
	Save(chatopsType string, webhookURL string, uuid uuid.UUID) error
	Delete(uuid uuid.UUID) error
}

type WebhookExtensionStore interface {
	Get() (schema.WebhookExtension, error)
	Save(webhookURL string, secret sql.NullString, uuid uuid.UUID) error
	Delete(uuid uuid.UUID) error
}
