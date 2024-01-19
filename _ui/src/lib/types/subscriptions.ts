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
