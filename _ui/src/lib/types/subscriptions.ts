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
