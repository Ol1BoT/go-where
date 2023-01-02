package postgres

import (
	"testing"
)

type GetCommunityParams struct {
	// Limit How many items to return at one time (max 100)
	Limit *int32 `form:"limit,omitempty" json:"limit,omitempty"`

	// Offset How to order the items returned, defaults to order by id
	Offset *int32 `form:"offset,omitempty" json:"offset,omitempty"`

	// LastName search by last name
	LastName *string `form:"last_name,omitempty" json:"last_name,omitempty"`

	// FirstName search by first name
	FirstName string `form:"first_name,omitempty" json:"first_name,omitempty"`
}

func TestXxx(t *testing.T) {

	test := &GetCommunityParams{
		Limit:     Ptr(int32(50)),
		Offset:    Ptr(int32(1)),
		LastName:  Ptr("BoT"),
		FirstName: "Oli",
	}

	sql := `SELECT * FROM PERSONS`

	rv, err := ConstructAndQuery(sql, "json", test)
	if err != nil {
		t.Error(err)
	}

	t.Log(rv.Params)
	t.Log(rv.StringQuery)

}

func Ptr[T any](v T) *T {
	return &v
}
