import type { PaginationMeta } from '$lib/types/api';
import type {
	GetAllSubscription,
	IncidentsForSubscription,
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

	async Delete(subscriptionID: string) {
		return axios.delete<{ data: { msg: string } }>(`/subscriptions/${subscriptionID}`);
	}

	async GetAll(serviceName: string, pageNumber: number, pageLimit: number) {
		return axios.get<{
			data: GetAllSubscription[];
			meta: PaginationMeta;
		}>('/dashboard', {
			params: {
				service_name: serviceName,
				page_number: pageNumber,
				pageLimit: pageLimit
			}
		});
	}

	async GetAllIncidents(subscriptionUUID: string, pageNumber: number, pageLimit: number) {
		return axios.get<{ data: IncidentsForSubscription; meta: PaginationMeta }>(
			`subscriptions/${subscriptionUUID}/incidents`,
			{
				params: {
					page_number: pageNumber,
					page_limit: pageLimit
				}
			}
		);
	}
}
