export type SaveChatOps = {
	uuid?: string;
	type: 'slack' | 'discord' | 'msteams';
	webhook_url: string;
};

export type SaveSquadcast = {
	uuid?: string;
	webhook_url: string;
};

export type SavePagerduty = {
	uuid?: string;
	routing_key: string;
};

export type SaveWebhook = {
	uuid?: string;
	webhook_url: string;
	secret?: string;
};

export type GetChatOps = {
	slack: {
		webhook_url: string;
		uuid: string;
		is_configured: boolean;
	};
	msteams: {
		webhook_url: string;
		uuid: string;
		is_configured: boolean;
	};
	discord: {
		webhook_url: string;
		uuid: string;
		is_configured: boolean;
	};
};

export type GetIncidentManagement = {
	squadcast: {
		is_configured: boolean;
		webhook_url: string;
		uuid: string;
	};
	pagerduty: {
		is_configured: boolean;
		routing_key: string;
		uuid: string;
	};
};

export type GetWebhook = {
	is_configured: boolean;
	uuid: string;
	webhook_url: string;
	secret: string;
};

export type DeleteChatopsData = {
	uuid: string;
	type: string;
};
