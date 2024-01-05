import type { AddSubscription, ServiceForSubscription } from '$lib/types/subscriptions';
import axios from 'axios';

export class SubscriptionAPI {
	async ServicesList(query: string) {
		return axios.get<{ data: ServiceForSubscription[] }>('/subscriptions/services', {
			params: {
				query: query
			}
		});
	}

	async Add(subscription: AddSubscription) {
		return axios.post<{ data: { msg: string } }>('/subscriptions/add-service', subscription);
	}
}
