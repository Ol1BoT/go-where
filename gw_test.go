package gw

import (
	"testing"
)

type GetParams struct {
	Limit     *int32  `form:"limit,omitempty" json:"limit,omitempty"`
	Offset    *int32  `form:"offset,omitempty" json:"offset,omitempty"`
	LastName  *string `form:"last_name,omitempty" json:"last_name,omitempty"`
	FirstName string  `form:"first_name,omitempty" json:"first_name,omitempty"`
	OrderBy   string  `form:"order_by,omitempty" json:"order_by,omitempty"`
}

type UpdateParams struct {
	LastName  *string `form:"last_name,omitempty" json:"last_name,omitempty"`
	FirstName string  `form:"first_name,omitempty" json:"first_name,omitempty"`
}

type UpdateWhereParams struct {
	PersonId int `json:"person_id"`
}

func TestOne(t *testing.T) {

	params := &GetParams{
		Limit:     Ptr(int32(50)),
		Offset:    Ptr(int32(1)),
		LastName:  nil,
		FirstName: "Oli",
		OrderBy:   "last_name",
	}

	cfg := &Config{
		QueryType: "pgx",
		WhereType: "AND",
		Tag:       "json",
	}

	b := NewBuilder(cfg)

	query := "SELECT * FROM public.person"

	rv, err := b.MakeSelectQuery(query, params)
	if err != nil {
		t.Error(err)
	}

	t.Log(rv.Params)
	t.Log(rv.StringQuery)

	up := &UpdateParams{
		FirstName: "Ol1",
		LastName:  Ptr("BoT"),
	}

	uw := &UpdateWhereParams{
		PersonId: 1,
	}

	update, err := b.MakeUpdateQuery("public.person", up, uw)
	if err != nil {
		t.Log(err)
	}

	t.Log(update.Params...)
	t.Log(update.StringQuery)

}

func TestTwo(t *testing.T) {

	params := &GetParams{
		Limit:     Ptr(int32(50)),
		Offset:    Ptr(int32(1)),
		LastName:  nil,
		FirstName: "Oli",
		OrderBy:   "last_name",
	}

	cfg := &Config{
		QueryType: "std",
		WhereType: "OR",
		Tag:       "json",
	}

	b := NewBuilder(cfg)

	query := "SELECT * FROM public.person"

	rv, err := b.MakeSelectQuery(query, params)
	if err != nil {
		t.Error(err)
	}

	t.Log(rv.Params)
	t.Log(rv.StringQuery)

	up := &UpdateParams{
		FirstName: "Ol1",
		LastName:  Ptr("BoT"),
	}

	uw := &UpdateWhereParams{
		PersonId: 1,
	}

	update, err := b.MakeUpdateQuery("public.person", up, uw)
	if err != nil {
		t.Log(err)
	}

	t.Log(update.Params...)
	t.Log(update.StringQuery)

}

func Ptr[T any](v T) *T {
	return &v
}
