// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package repo

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type ShortenedUrl struct {
	ID               pgtype.UUID        `json:"id"`
	OriginalUrl      string             `json:"original_url"`
	Hash             string             `json:"hash"`
	NumberOfAccesses int32              `json:"number_of_accesses"`
	CreatedAt        pgtype.Timestamptz `json:"created_at"`
	UpdatedAt        pgtype.Timestamptz `json:"updated_at"`
}
