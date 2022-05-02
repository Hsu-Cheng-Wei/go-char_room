package models

type UserClaim struct {
	ID   string
	Name string
	Iat  uint32 `json:"iat"`
	Exp  uint32 `json:"exp"`
}
