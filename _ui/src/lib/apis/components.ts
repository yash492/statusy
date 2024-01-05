import type { ComponentsForService } from '$lib/types/components';
import axios from 'axios';

export class ComponentsAPI {
	async ComponentsForService(serviceID: number) {
		return axios.get<{ data: ComponentsForService[] }>(`/services/${serviceID}/components`);
	}
}
