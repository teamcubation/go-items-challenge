package item

type Response struct {
	TotalPages int    `json:"totalPages"`
	Data       []Item `json:"data"`
}
