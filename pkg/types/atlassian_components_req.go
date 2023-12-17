package types

import "time"

type AtlassianComponentsReq struct {
	Components []atlassianComponent `json:"components"`
}

type atlassianComponent struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	Status             string    `json:"status"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	Position           int       `json:"position"`
	Description        *string   `json:"description"`
	Showcase           bool      `json:"showcase"`
	GroupID            string    `json:"group_id"`
	PageID             string    `json:"page_id"`
	Group              bool      `json:"group"`
	OnlyShowIfDegraded bool      `json:"only_show_if_degraded"`
}
