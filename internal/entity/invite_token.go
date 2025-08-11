package entity

import (
	"github.com/uptrace/bun"
)

type InviteToken struct {
	bun.BaseModel `bun:"table:invite_tokens"`
	Id            int64  `bun:"id,pk,autoincrement"`
	Token         string `bun:"token"`
}
