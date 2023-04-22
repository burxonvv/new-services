package models

type Policy struct {
	User   string `json:"user"`
	Domain string `json:"domain"`
	Action string `json:"action"`
}

type RoleRequest struct {
	Id   string `json:"id"`
	Role string `json:"role"`
}
