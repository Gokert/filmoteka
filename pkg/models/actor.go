package models

type ActorItem struct {
	Id       uint64 `json:"id"`
	Name     string `json:"name"`
	Gender   string `json:"gen"`
	Birthday string `json:"birthdate"`
}
