import type {
	GetChatOps,
	GetIncidentManagement,
	GetWebhook,
	SaveChatOps,
	SavePagerduty,
	SaveSquadcast,
	SaveWebhook
} from '$lib/types/integrations';
import axios from 'axios';

export class IntegrationsAPI {
	async GetChatOps() {
		return axios.get<{ data: GetChatOps }>('integrations/chatops');
	}

	async SaveChatOps(data: SaveChatOps) {
		return axios.put('integrations/chatops', data);
	}

	async GetIncidentManagement() {
		return axios.get<{ data: GetIncidentManagement }>('/integrations/incident-management');
	}

	async SaveSquadcast(data: SaveSquadcast) {
		return axios.put('/integrations/incident-management/squadcast', data);
	}

	async SavePagerduty(data: SavePagerduty) {
		return axios.put('/integrations/incident-management/pagerduty', data);
	}

	async SaveWebhook(data: SaveWebhook) {
		return axios.put('/integrations/webhook', data);
	}

	async GetWebhook() {
		return axios.get<{ data: GetWebhook }>('integrations/webhook');
	}

	async DeleteSquadcast(uuid: string) {
		return axios.delete(`/integrations/incident-management/squadcast/${uuid}`);
	}

	async DeletePagerduty(uuid: string) {
		return axios.delete(`/integrations/incident-management/pagerduty/${uuid}`);
	}
	async DeleteChatOps(uuid: string, type: string) {
		return axios.delete(`/integrations/chatops/${uuid}?type=${type}`);
	}
	async DeleteWebhook(uuid: string) {
		return axios.delete(`/integrations/webhook/${uuid}`);
	}
}
