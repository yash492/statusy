import type {
	GetAllSubscription,
	SaveSubscription,
	ServiceForSubscription,
	SubscriptionWithComponents
} from '$lib/types/subscriptions';
import axios from 'axios';

export class SubscriptionAPI {
	async ServicesList(query: string) {
		return axios.get<{ data: ServiceForSubscription[] }>('/subscriptions/services', {
			params: {
				query: query
			}
		});
	}

	async Add(subscription: SaveSubscription) {
		return axios.post<{ data: { msg: string } }>('/subscriptions', subscription);
	}

	async GetByID(subscriptionID: string) {
		return axios.get<{ data: SubscriptionWithComponents }>(`/subscriptions/${subscriptionID}`);
	}

	async Update(subscription: SaveSubscription, subscriptionID: string) {
		return axios.put<{ data: { msg: string } }>(`/subscriptions/${subscriptionID}`, subscription);
	}

	async GetAll(serviceName: string, pageNumber: number, pageLimit: number) {
		return axios.get<{
			data: GetAllSubscription[];
			meta: { total_count: number; page_number: number; page_limit: number };
		}>('/dashboard', {
			params: {
				service_name: serviceName,
				page_number: pageNumber,
				pageLimit: pageLimit
			}
		});
	}
}
