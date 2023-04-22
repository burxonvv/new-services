package models

// Error ...
type Error struct {
	Message string `json:"message"`
}

// StandardErrorModel ...
type StandardErrorModel struct {
	Error Error `json:"error"`
}

type Success struct {
	Message string `json:"message"`
}