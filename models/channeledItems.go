package models

type ChanneledItems struct {
	Items   int `json:"items"`
	Channel Ch  `json:"channels"`
}
