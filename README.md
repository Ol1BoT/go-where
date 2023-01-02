# GO-WHERE

A basic SQL Query builder where you type out your entire query before you hit the `WHERE` clause

The idea is that when a user providers query params via REST, those query params are used to create a `WHERE x AND y AND z` statement.

These params are usually provided in a struct by your router/framework of choice, so they can just be passed into the query constructor.

by default it uses json tags to dictate what the field names in the where clause should be called, however you can change the tags.

OFFSET and LIMIT and reserved keywords, so if they are passed into the query constructor, they will be applied at the end of the query and LIMIT and OFFSET, one can also not live without the other.


```go
type GetParams struct {
	Limit     *int32  `form:"limit,omitempty" json:"limit,omitempty"`
	Offset    *int32  `form:"offset,omitempty" json:"offset,omitempty"`
	LastName  *string `form:"last_name,omitempty" json:"last_name,omitempty"`
	FirstName string  `form:"first_name,omitempty" json:"first_name,omitempty"`
}

func main() {

	test := &GetParams{
		Limit:     Ptr(int32(50)),
		Offset:    Ptr(int32(1)),
		LastName:  Ptr("BoT"),
		FirstName: "Ol1",
	}

	sql := `SELECT * FROM PERSONS`

	rv, err := ConstructAndQuery(sql, "json", test)
	if err != nil {
		panic(err)
	}

	fmt.Println(rv.Params) // [BoT Ol1]
	fmt.Println(rv.StringQuery) // SELECT * FROM PERSONS WHERE last_name = $1 AND first_name = $2 LIMIT 50 OFFSET 1 

}

func Ptr[T any](v T) *T {
	return &v
}
```