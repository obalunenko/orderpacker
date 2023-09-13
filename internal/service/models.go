package service

type PackRequest struct {
	Items uint `json:"items"`
}

type PackResponse struct {
	Boxes []uint `json:"boxes"`
}
