package types

type AtlassianComponentsReq struct {
	Components []AtlassianComponent `json:"components"`
}

type AtlassianComponent struct {
	Name string `json:"name"`
}
