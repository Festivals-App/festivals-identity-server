package token

type APIKey struct {
	ID      int    `json:"api_key_id" sql:"api_key_id"`
	Key     string `json:"api_key" sql:"api_key"`
	Comment string `json:"api_key_comment" sql:"api_key_comment"`
}
