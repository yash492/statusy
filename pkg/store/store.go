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
	GetByProviderIDs(providerIDs []string) ([]schema.Incident, error)
	CreateIncidentUpdates(incidentUpdates []schema.IncidentUpdate) ([]schema.IncidentUpdate, error)
	CreateIncidentComponents(incidentComponents []schema.IncidentComponent) error
	GetIncidentUpdatesByProviderIDs(providerIDs []string) ([]schema.IncidentUpdate, error)
}

type SubscriptionStore interface {
	GetAllServicesForSubscriptions(serviceName string) ([]schema.ServicesForSubsciption, error)
	Create(serviceID uint, componentIDs []uint, isAllComponents bool) error
	GetByID(subscriptionID uuid.UUID) (schema.SubscriptionWithService, error)
	GetWithComponents(subscriptionID uuid.UUID) ([]schema.SubscriptionWithComponent, error)
	Update(subscriptionID uuid.UUID, componentIDs []uint, isAllComponents bool) error
	GetForIncidentUpdates(incidentUpdateID uint) ([]schema.SubscriptionForIncidentUpdate, error)
	DashboardSubscription(serviceName string, offset, limit uint) ([]schema.DashboardSubscription, error)
	GetIncidentsForSubscription(subscriptionUUID uuid.UUID, offset, limit uint) ([]schema.SubscriptionIncident, error)
	Delete(subscriptionID uuid.UUID) error
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
