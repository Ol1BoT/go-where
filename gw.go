package gw

import (
	"fmt"
	"reflect"
	"strings"
)

type Config struct {
	Tag       string //will default to use JSON tag if blank
	QueryType string //Should be "std" or "pgx", will default to "std" if blank
	WhereType string //will default to AND if blank
}

type QueryStatement struct {
	StringQuery string
	Params      []any
}

type Builder struct {
	QueryType string
	WhereType string
	Tag       string
}

const (
	AND     = "AND"
	OR      = "OR"
	OFFSET  = "OFFSET"
	LIMIT   = "LIMIT"
	ORDERBY = "ORDERBY"
	JSON    = "json"
	PGX     = "pgx"
	STD     = "std"
)

var (
	PgxConfig = Config{
		Tag:       JSON,
		QueryType: "pgx",
	}

	StdConfig = Config{
		Tag:       JSON,
		QueryType: "std",
	}
)

func NewBuilder(c *Config) *Builder {
	return &Builder{
		WhereType: c.WhereType,
		QueryType: c.QueryType,
		Tag:       c.Tag,
	}
}

func (b *Builder) MakeUpdateQuery(table string, update any, where any) (*QueryStatement, error) {
	if b.Tag == "" {
		b.Tag = JSON
	}

	if b.WhereType == "" {
		b.WhereType = AND
	}

	if b.WhereType != AND && b.WhereType != OR {
		return nil, fmt.Errorf("invalid where type, must be AND or OR, you provided: %s", b.WhereType)
	}

	qs := &QueryStatement{}

	sb := strings.Builder{}

	uType := reflect.TypeOf(update)
	if uType == nil {
		return nil, fmt.Errorf("No fields to update")
	}

	sb.WriteString("UPDATE " + table + " SET")

	uType, updateElement := getElement(uType, update)

	if uType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("update must be a struct, or pointer to a struct")
	}

	updateVals := make([]any, 0)

	for i := 0; i < uType.NumField(); i++ {
		var val reflect.Value
		field := reflect.ValueOf(updateElement).Field(i)

		if field.Kind() == reflect.Pointer {
			if field.IsNil() {
				continue
			}
			val = field.Elem()

		} else {
			val = field
		}
		if val.Kind() == reflect.Pointer {
			val = val.Elem()
		}

		tg := uType.Field(i).Tag.Get(b.Tag)
		if tg == "" {
			return nil, fmt.Errorf("%s does not exist on %s", b.Tag, uType.Field(i).Name)
		}

		if uType.NumField()-1 == i {
			if b.QueryType == PGX {
				sb.WriteString(fmt.Sprintf(" %s = $%d", strings.Split(tg, ",")[0], len(updateVals)+1))
			} else {
				sb.WriteString(fmt.Sprintf(" %s = ?", strings.Split(tg, ",")[0]))
			}
			updateVals = append(updateVals, val)
			continue
		}

		if b.QueryType == PGX {
			sb.WriteString(fmt.Sprintf(" %s = $%d,", strings.Split(tg, ",")[0], len(updateVals)+1))
		} else {
			sb.WriteString(fmt.Sprintf(" %s = ?,", strings.Split(tg, ",")[0]))
		}
		updateVals = append(updateVals, val)
	}

	wType := reflect.TypeOf(where)
	if wType == nil {
		return &QueryStatement{
			StringQuery: sb.String(),
			Params:      updateVals,
		}, nil
	}

	wType, whereElement := getElement(wType, where)

	if wType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("where must be a struct, or pointer to a struct")
	}

	var limit string
	var offset string
	var orderBy string
	whereVals := make([]any, 0)

	for i := 0; i < wType.NumField(); i++ {
		var val reflect.Value
		field := reflect.ValueOf(whereElement).Field(i)

		if field.Kind() == reflect.Pointer {
			if field.IsNil() {
				continue
			}
			val = field.Elem()

		} else {
			val = field
		}
		if val.Kind() == reflect.Pointer {
			val = val.Elem()
		}

		if strings.ToUpper(wType.Field(i).Name) == ORDERBY {
			if val.CanInt() {
				orderBy = fmt.Sprintf("%d", val.Int())
			} else {
				orderBy = fmt.Sprintf("%s", val)
			}
			continue
		}

		if strings.ToUpper(wType.Field(i).Name) == OFFSET {
			if val.CanInt() {
				offset = fmt.Sprintf("%d", val.Int())
			} else {
				offset = fmt.Sprintf("%s", val)
			}
			continue
		}

		if strings.ToUpper(wType.Field(i).Name) == LIMIT {
			if val.CanInt() {
				limit = fmt.Sprintf("%d", val.Int())
			} else {
				limit = fmt.Sprintf("%s", val)
			}
			continue
		}

		if i == 0 {
			sb.WriteString(" WHERE")
			tg := wType.Field(i).Tag.Get(b.Tag)
			if b.QueryType == PGX {
				sb.WriteString(fmt.Sprintf(" %s = $%d", strings.Split(tg, ",")[0], len(whereVals)+len(updateVals)+1))
			} else {
				sb.WriteString(fmt.Sprintf(" %s = ?", strings.Split(tg, ",")[0]))
			}
			whereVals = append(whereVals, val)
			continue
		}

		tg := wType.Field(i).Tag.Get(b.Tag)
		if tg == "" {
			return nil, fmt.Errorf("%s does not exist on %s", b.Tag, wType.Field(i).Name)
		}
		if b.QueryType == PGX {
			sb.WriteString(fmt.Sprintf(" %s %s = $%d", b.WhereType, strings.Split(tg, ",")[0], len(whereVals)+len(updateVals)+1))
		} else {
			sb.WriteString(fmt.Sprintf(" %s %s = ?", b.WhereType, strings.Split(tg, ",")[0]))
		}
		updateVals = append(whereVals, val)
	}

	if orderBy != "" {
		ob := fmt.Sprintf(" ORDER BY %s", orderBy)
		sb.WriteString(ob)
	}

	if offset != "" && limit != "" {
		ender := fmt.Sprintf(" LIMIT %s OFFSET %s", limit, offset)
		sb.WriteString(ender)
	}

	var vals []any

	vals = append(vals, updateVals...)
	vals = append(vals, whereVals...)

	qs.Params = vals
	qs.StringQuery = sb.String()

	return qs, nil

}

func (b *Builder) MakeSelectQuery(query string, params any) (*QueryStatement, error) {

	if b.Tag == "" {
		b.Tag = JSON
	}

	if b.WhereType == "" {
		b.WhereType = AND
	}

	if b.WhereType != AND && b.WhereType != OR {
		return nil, fmt.Errorf("invalid where type, must be AND or OR, you provided: %s", b.WhereType)
	}

	qs := &QueryStatement{}

	sb := strings.Builder{}

	sb.WriteString(query)

	tp := reflect.TypeOf(params)
	if tp == nil {
		qs.StringQuery = sb.String()
		return qs, nil
	}

	tp, element := getElement(tp, params)

	vals := make([]any, 0)
	var limit string
	var offset string
	var orderBy string

	for i := 0; i < tp.NumField(); i++ {
		var val reflect.Value
		field := reflect.ValueOf(element).Field(i)

		if field.Kind() == reflect.Pointer {
			if field.IsNil() {
				continue
			}
			val = field.Elem()

		} else {
			val = field
		}
		if val.Kind() == reflect.Pointer {
			val = val.Elem()
		}

		if strings.ToUpper(tp.Field(i).Name) == ORDERBY {
			if val.CanInt() {
				orderBy = fmt.Sprintf("%d", val.Int())
			} else {
				orderBy = fmt.Sprintf("%s", val)
			}
			continue
		}

		if strings.ToUpper(tp.Field(i).Name) == OFFSET {
			if val.CanInt() {
				offset = fmt.Sprintf("%d", val.Int())
			} else {
				offset = fmt.Sprintf("%s", val)
			}
			continue
		}

		if strings.ToUpper(tp.Field(i).Name) == LIMIT {
			if val.CanInt() {
				limit = fmt.Sprintf("%d", val.Int())
			} else {
				limit = fmt.Sprintf("%s", val)
			}
			continue
		}

		if len(vals) == 0 {
			sb.WriteString(" WHERE")
			tg := tp.Field(i).Tag.Get(b.Tag)
			if b.QueryType == PGX {
				sb.WriteString(fmt.Sprintf(" %s = $%d", strings.Split(tg, ",")[0], len(vals)+1))
			} else {
				sb.WriteString(fmt.Sprintf(" %s = ?", strings.Split(tg, ",")[0]))
			}
			vals = append(vals, val)
			continue
		}

		tg := tp.Field(i).Tag.Get(b.Tag)
		if tg == "" {
			return nil, fmt.Errorf("%s does not exist on %s", b.Tag, tp.Field(i).Name)
		}

		if b.QueryType == PGX {
			sb.WriteString(fmt.Sprintf(" %s %s = $%d", b.WhereType, strings.Split(tg, ",")[0], len(vals)+1))
		} else {
			sb.WriteString(fmt.Sprintf(" %s %s = ?", b.WhereType, strings.Split(tg, ",")[0]))
		}
		vals = append(vals, val)
	}

	if orderBy != "" {
		ob := fmt.Sprintf(" ORDER BY %s", orderBy)
		sb.WriteString(ob)
	}

	if offset != "" && limit != "" {
		ender := fmt.Sprintf(" LIMIT %s OFFSET %s", limit, offset)
		sb.WriteString(ender)
	}

	qs.Params = vals
	qs.StringQuery = sb.String()

	return qs, nil
}

func getElement(rt reflect.Type, target any) (reflect.Type, any) {

	var element any

	if rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
		element = reflect.ValueOf(target).Elem().Interface()
	} else {
		element = target
	}

	return rt, element

}
