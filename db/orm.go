package db

import (
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func sqlTableNameString(s interface{}) string {
	name := reflect.TypeOf(s).Name()
	return toSnakeCase(name + "s")
}

func sqlFieldNames(s interface{}, exclude_props ...string) string {
	t := reflect.TypeOf(s)
	sb := strings.Builder{}

OUTER:
	for i := 0; i < t.NumField(); i++ {
		props, ok := t.Field(i).Tag.Lookup("sql_props")
		if !ok {
			continue
		}

		for _, val := range exclude_props {
			if strings.Contains(props, val) {
				continue OUTER
			}
		}

		name, ok := t.Field(i).Tag.Lookup("sql_name")
		if !ok {
			continue
		}

		sb.WriteString(name)
		if i < t.NumField()-1 {
			sb.WriteString(", ")
		}
	}
	return sb.String()
}

func sqlFieldValues(s interface{}, exclude_props ...string) string {
	v := reflect.ValueOf(s)
	sb := strings.Builder{}

OUTER:
	for i := 0; i < v.NumField(); i++ {
		props, ok := v.Type().Field(i).Tag.Lookup("sql_props")

		if !ok {
			continue
		}

		for _, val := range exclude_props {
			if strings.Contains(props, val) {
				continue OUTER
			}
		}

		sb.WriteString(fmt.Sprintf("'%+v'", v.Field(i)))
		if i < v.NumField()-1 {
			sb.WriteString(", ")
		}
	}
	return sb.String()
}

func sqlFieldValue(s interface{}, val string) (string, error) {
	v := reflect.ValueOf(s)
	for i := 0; i < v.NumField(); i++ {
		if v.Type().Field(i).Tag.Get("sql_name") == val {
			return fmt.Sprintf("'%+v'", v.Field(i)), nil
		}
	}
	tableName := sqlTableNameString(s)
	return "", fmt.Errorf("field not found in -> %s", tableName)
}

func SchemaString(structs ...interface{}) ([]string, error) {
	final := []string{}
	sb := strings.Builder{}

	for _, s := range structs {
		t := reflect.TypeOf(s)
		fcount := t.NumField()
		struct_name := sqlTableNameString(s)

		// generate field strings
		sb.WriteString("CREATE TABLE IF NOT EXISTS " + struct_name + " (\n")

		for i := range fcount {
			v, ok := t.Field(i).Tag.Lookup("sql_name")
			if !ok {
				continue
			}
			sb.WriteString("\t" + v)
			end, ok := t.Field(i).Tag.Lookup("sql_props")
			if !ok {
				continue
			}
			sb.WriteString("\t\t\t" + end)

			// generate foreign key strings
			fullKey, ok := t.Field(i).Tag.Lookup("sql_fk")
			if !ok {
				if i < fcount-1 {
					sb.WriteString(",\n")
				}
				continue
			}

			split := strings.Split(fullKey, ".")
			if len(split) != 2 {
				return nil, fmt.Errorf("sql_fk tag has an invalid value")
			}

			fk := fmt.Sprintf(" REFERENCES %ss(%s)", split[0], split[1])
			sb.WriteString(fk)

			if i < fcount-1 {
				sb.WriteString(",\n")
			}
		}

		sb.WriteString("\n)")
		final = append(final, sb.String())
		sb.Reset()
	}

	return final, nil
}

func Insert(s interface{}) (sql.Result, error) {
	query := "INSERT INTO " + sqlTableNameString(s) + " (" + sqlFieldNames(s, "AUTOINCREMENT") + ") VALUES (" + sqlFieldValues(s, "AUTOINCREMENT") + ")"
	fmt.Printf("\nInsert Query: %s\n\n\n", query)
	return dbInstance.Exec(query)
}

func Remove(s interface{}) (sql.Result, error) {
	id, err := sqlFieldValue(s, "id")
	if err != nil {
		return nil, err
	}

	query := "DELETE FROM " + sqlTableNameString(s) + " WHERE id = " + id
	fmt.Printf("\nRemove Query: %s\n\n\n", query)
	return dbInstance.Exec(query)
}

func Update() {
	panic("SQL ORM Update not implemented")
}

func Query(entry interface{}, values ...string) (*sql.Rows, error) {
	table_name := sqlTableNameString(entry)
	query := strings.Builder{}
	query.WriteString("SELECT * FROM " + table_name + " WHERE ")
	for i, val := range values {
		first_letter := string(val[0])
		rest := val[1:]
		final := strings.ToUpper(first_letter) + rest
		query.WriteString(table_name + "." + val + " = " + fmt.Sprintf("%+v", reflect.ValueOf(entry).FieldByName(final)))
		if i < len(values)-1 {
			query.WriteString(" AND ")
		}
	}
	fmt.Printf("\nQuery: %s\n\n\n", query.String())
	return dbInstance.Query(query.String())
}

func Exec(s string, args ...any) (sql.Result, error) {
	return dbInstance.Exec(s)
}

func CloseDB() {
	dbInstance.Close()
}
