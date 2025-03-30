package token

type ServiceKey struct {
	ID      int    `json:"service_key_id" sql:"service_key_id"`
	Key     string `json:"service_key" sql:"service_key"`
	Comment string `json:"service_key_comment" sql:"service_key_comment"`
}
