package postgres

import (
	"testing"
)

type GetParams struct {
	Limit     *int32  `form:"limit,omitempty" json:"limit,omitempty"`
	Offset    *int32  `form:"offset,omitempty" json:"offset,omitempty"`
	LastName  *string `form:"last_name,omitempty" json:"last_name,omitempty"`
	FirstName string  `form:"first_name,omitempty" json:"first_name,omitempty"`
}

func TestOne(t *testing.T) {

	test := &GetParams{
		Limit:     Ptr(int32(50)),
		Offset:    Ptr(int32(1)),
		LastName:  nil,
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
