# GO-WHERE

A basic SQL Query builder helper where you type out your entire query before you hit the `WHERE` clause

The idea is that when a user providers query params via REST, those query params are used to create a `WHERE x AND y AND z` statement.

These params are usually provided in a struct by your router/framework of choice, so they can just be passed into the query constructor.

by default it uses `json` struct tags to dictate what the field names in the `WHERE` clause should be called, however you can change the tags.

`Offset`, `Limit` and `OrderBy` are reserved keywords in the struct, so if they are passed into the query constructor, they will be applied at the end of the query. Concerning `limit` and `Offset` one can also not live without the other, thus would be ignored if either one is missing.


```go
package main

import (
	"fmt"

	"github.com/Ol1BoT/go-where"
)

type GetParams struct {
	Limit     *int32  `form:"limit,omitempty" json:"limit,omitempty"`
	Offset    *int32  `form:"offset,omitempty" json:"offset,omitempty"`
	LastName  *string `form:"last_name,omitempty" json:"last_name,omitempty"`
	FirstName string  `form:"first_name,omitempty" json:"first_name,omitempty"`
}

type UpdateParams struct {
	LastName  *string `form:"last_name,omitempty" json:"last_name,omitempty"`
	FirstName string  `form:"first_name,omitempty" json:"first_name,omitempty"`
}

type UpdateWhereParams struct {
	PersonId int `json:"person_id"`
}

func main() {

	params := &GetParams{
		Limit:     Ptr(int32(50)),
		Offset:    Ptr(int32(1)),
		LastName:  Ptr("BoT"),
		FirstName: "Ol1",
	}

	cfg := &Config{
		QueryType: "pgx",
		WhereType: "AND",
		Tag:       "json",
	}

	b := gw.NewBuilder(cfg)

	query := `SELECT * FROM public.person`

	rv, err := b.MakeSelectQuery(query,params)
	if err != nil {
		panic(err)
	}

	fmt.Println(rv.Params) // [BoT Ol1]
	fmt.Println(rv.StringQuery) // SELECT * FROM public.person WHERE last_name = $1 AND first_name = $2 LIMIT 50 OFFSET 1 

	up := &UpdateParams{
		FirstName: "Ol1",
		LastName:  Ptr("BoT"),
	}

	uw := &UpdateWhereParams{
		PersonId: 1,
	}

	update, err := b.MakeUpdateQuery("public.person", up, uw)
	if err != nil {
		panic(err)
	}

	fmt.Println(update.Params...) // [BoT Ol1 1]
	fmt.Println(update.StringQuery) //UPDATE public.person SET last_name = $1, first_name = $2 WHERE person_id = $3

}

func Ptr[T any](v T) *T {
	return &v
}
```