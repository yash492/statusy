export type ServiceForSubscription = {
	name: string;
	id: number;
};

export type SaveSubscription = {
	service_id: number;
	is_all_components: boolean;
	custom_components: number[];
};

export type SubscriptionWithComponents = {
	service_id: number;
	is_all_components: boolean;
	service_name: string;
	uuid: string;
	components: components[];
};

type components = {
	name: string;
	id: number;
	is_configured: boolean;
};

export type DashboardTable = {
	serviceName: string;
	isDown: boolean;
	incident: string;
	subscriptionUUID: string;
	incidentLink: string;
};

export type GetAllSubscription = {
	incident_id: number;
	service_id: number;
	service_name: string;
	subscription_uuid: string;
	incident_name: string;
	incident_link: string;
	incident_impact: string;
	is_down: boolean;
};

export type IncidentsForSubscription = {
	service_name: string;
	service_id: number;
	is_all_components_configured: boolean;
	components: { id: number; name: string }[]; // Adjust the type accordingly if components have a specific structure
	incidents: Incident[];
};

export type Incident = {
	id: number;
	last_updated_status_time: string;
	status: string;
	created_at: Date;
	name: string;
	link: string;
	normalised_status: string;
};

export type SubscriptionIncidentsTable = {
	name: string;
	status: string;
	createdAt: Date;
	link: string;
	normalisedStatus: string;
};
