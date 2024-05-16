package service

// PackRequest represents a request to pack items.
type PackRequest struct {
	Items uint `json:"items" format:"uint" example:"543"`
}

// Pack represents a pack of items.
type Pack struct {
	Box      uint `json:"box" format:"uint" example:"50"`
	Quantity uint `json:"quantity" format:"uint" example:"3"`
}

// PackResponse represents a response to a pack request.
type PackResponse struct {
	Packs []Pack `json:"packs,omitempty"`
}

// HTTPError represents an HTTP error.
type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Bad request"`
}
