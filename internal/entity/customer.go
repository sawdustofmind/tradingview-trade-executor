package entity

import (
	"github.com/uptrace/bun"
)

type Customer struct {
	bun.BaseModel `bun:"table:customers"`

	Id                 int64   `bun:"id,pk,autoincrement"`
	Username           string  `bun:"username"`
	Password           string  `bun:"password"`
	LegalName          string  `bun:"legal_name"`
	Gender             string  `bun:"gender"`
	Country            string  `bun:"country"`
	PhoneNumber        string  `bun:"phone_number"`
	ImageBase64        string  `bun:"image_base64"`
	BybitApiKey        *string `bun:"bybit_api_key"`
	BybitTestApiKey    *string `bun:"bybit_test_api_key"`
	BybitApiSecret     *string `bun:"bybit_api_secret" json:"-"`
	BybitTestApiSecret *string `bun:"bybit_test_api_secret"`
}
