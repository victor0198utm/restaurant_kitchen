package models

type ChanneledOrder struct {
	Order_id     int              `json:"order_id"`
	Table_id     int              `json:"table_id"`
	Waiter_id    int              `json:"waiter_id"`
	CItems       []ChanneledItems `json:"channeled_items"`
	Priority     int              `json:"priority"`
	Max_wait     int              `json:"max_wait"`
	Pick_up_time int              `json:"pick_up_time"`
}
