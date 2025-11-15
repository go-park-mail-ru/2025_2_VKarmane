package support

type CreateRequestInput struct {
	Category string `json:"category" validate:"required"`
	Message  string `json:"message" validate:"required"`
}

type CreateRequestOutput struct {
	ID       int    `json:"id"`
	Status   string `json:"status"`
	Category string `json:"category"`
	Message  string `json:"message"`
}

type UserRequestsOutput struct {
	Requests []SupportItem `json:"requests"`
}

type SupportItem struct {
	ID       int    `json:"id"`
	Category string `json:"category"`
	Status   string `json:"status"`
	Message  string `json:"message"`
}

type UpdateStatusInput struct {
	Status string `json:"status" validate:"required"`
}

type StatsItem struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

type StatsOutput struct {
	Items []StatsItem `json:"items"`
}
