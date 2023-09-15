package service

type PackRequest struct {
	Items uint `json:"items"`
}

type Pack struct {
	Box      uint `json:"box"`
	Quantity uint `json:"quantity"`
}
type PackResponse struct {
	Packs []Pack `json:"packs"`
}
