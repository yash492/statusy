package datadog

import "resty.dev/v3"

func Register(client *resty.Client) {
	registerEU(client)
	registerUS3(client)
	registerUS5(client)
	registerAP1(client)
	registerGov(client)
	registerAP2(client)
}
