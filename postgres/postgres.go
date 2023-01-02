package postgres

import (
	"fmt"
	"reflect"
	"strings"
)

type QueryStatement struct {
	StringQuery string
	Params      []any
}

const (
	AND = "AND"
	OR  = "OR"
)

func ConstructAndQuery(query string, tag string, params any) (*QueryStatement, error) {
	return ConstructQuery(query, tag, AND, params)
}

func ConstructOrQuery(query string, tag string, params any) (*QueryStatement, error) {
	return ConstructQuery(query, tag, OR, params)
}

func ConstructQuery(query string, tag string, whereType string, params any) (*QueryStatement, error) {
	if tag == "" {
		tag = "json"
	}

	if whereType != AND && whereType != OR {
		return nil, fmt.Errorf("invalid where type, must be AND or OR, you provided: %s", whereType)
	}

	qs := &QueryStatement{}

	sb := strings.Builder{}

	sb.WriteString(query)

	tp := reflect.TypeOf(params)
	if tp == nil {
		qs.StringQuery = sb.String()
		return qs, nil
	}

	var element any

	if tp.Kind() == reflect.Pointer {
		tp = tp.Elem()
		element = reflect.ValueOf(params).Elem().Interface()
	} else {
		element = params
	}

	vals := make([]any, 0)
	var limit string
	var offset string

	for i := 0; i < tp.NumField(); i++ {
		var val reflect.Value
		field := reflect.ValueOf(element).Field(i)

		if field.Kind() == reflect.Pointer {
			val = reflect.ValueOf(element).Field(i).Elem()
		} else {
			val = reflect.ValueOf(element).Field(i)
		}
		if val.Kind() == reflect.Pointer {
			val = val.Elem()
		}

		if strings.ToUpper(tp.Field(i).Name) == "OFFSET" {
			if val.CanInt() {
				offset = fmt.Sprintf("%d", val.Int())
			} else {
				offset = fmt.Sprintf("%s", val)
			}
			continue
		}

		if strings.ToUpper(tp.Field(i).Name) == "LIMIT" {
			if val.CanInt() {
				limit = fmt.Sprintf("%d", val.Int())
			} else {
				limit = fmt.Sprintf("%s", val)
			}
			continue
		}

		if len(vals) == 0 {
			sb.WriteString(" WHERE")
			tg := tp.Field(i).Tag.Get(tag)
			sb.WriteString(fmt.Sprintf(" %s = $%d", strings.Split(tg, ",")[0], len(vals)+1))
			vals = append(vals, val)
			continue
		}

		tg := tp.Field(i).Tag.Get(tag)
		sb.WriteString(fmt.Sprintf(" %s %s = $%d", whereType, strings.Split(tg, ",")[0], len(vals)+1))
		vals = append(vals, val)
	}

	if offset != "" && limit != "" {

		ender := fmt.Sprintf(" LIMIT %s OFFSET %s", limit, offset)
		sb.WriteString(ender)

	}

	qs.Params = vals
	qs.StringQuery = sb.String()

	return qs, nil
}
