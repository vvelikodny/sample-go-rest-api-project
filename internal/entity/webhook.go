package entity

// Webhook represents an webhook record.
type Webhook struct {
	ID          int    `json:"id"`
	CityID      int    `json:"city_id"`
	CallbackURL string `json:"callback_url"`
}
