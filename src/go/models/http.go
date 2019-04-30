package models

// ActionRequest struct to define the structure of an Action Request
type ActionRequest struct {
	Name      string            `json:"name"`
	ExtraText map[string]string `json:"extraText"`
}

// ReadingRequest struct to define the structure od a Reading Request
type ReadingRequest struct {
	Name string `json:"name"`
}

// ReadingResponse struct to define the structure od a Reading Response
type ReadingResponse struct {
	Data    string  `json:"data"`
	Reading Reading `json:"reading"`
}
