package dto

type ErrorMessage struct {
	Error Error `json:"error"`
}
type Error struct {
	Description string `json:"description"`
	Code        string `json:"code,omitempty"`
}
