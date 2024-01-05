export type ServiceForSubscription = {
	name: string;
	id: number;
};

export type AddSubscription = {
	service_id: number;
	is_all_components: boolean;
	custom_components: number[];
};
