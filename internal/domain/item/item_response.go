package item

type ItemResponse struct {
	TotalPages int    `json:"totalPages"`
	Data       []Item `json:"data"`
}
