/**
作者:guangbo
模块：sql生成模块
说明：
创建时间：2015-10-30
**/
package GxMisc

import (
	"reflect"
	"strconv"
)

func GenerateCreateSql(info interface{}, tableId string) string {
	tableName := "tb_" + getTableName(info) + tableId
	key := ""

	ret := "CREATE TABLE " + tableName + " ("

	dataStruct := reflect.Indirect(reflect.ValueOf(info))
	dataStructType := dataStruct.Type()
	first := true
	for i := 0; i < dataStructType.NumField(); i++ {
		fieldType := dataStructType.Field(i)

		fieldTag := fieldType.Tag
		if fieldTag.Get("pk") == "true" {
			if key == "" {
				key = fieldType.Name
			} else {
				key += ", " + fieldType.Name
			}
		}

		if !first {
			ret += ","
		}

		switch fieldType.Type.Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if fieldTag.Get("type") == "time" {
				ret += " " + fieldType.Name + " DATETIME NOT NULL"
			} else {
				ret += " " + fieldType.Name + " INT NOT NULL"
			}
		case reflect.Float32, reflect.Float64:
			ret += " " + fieldType.Name + " FLOAT NOT NULL"
		case reflect.String:
			lenStr := fieldTag.Get("len")
			if lenStr != "" {
				ret += " " + fieldType.Name + " VARCHAR(" + lenStr + ") NOT NULL"
			} else {
				ret += " " + fieldType.Name + " VARCHAR(32) NOT NULL"
			}
		case reflect.Bool:
			ret += " " + fieldType.Name + " INT NOT NULL"
		case reflect.Slice:
			ret += " " + fieldType.Name + " BLOB NOT NULL"
		}
		first = false
	}

	if key != "" {
		ret += ", PRIMARY KEY (" + key + "))"
	} else {
		ret += ")"
	}

	return ret
}

// CREATE INDEX mytable_categoryid
// 　ON mytable (category_id);

func GenerateIndexSql(info interface{}, tableId string, indexId string) string {
	tableName := "tb_" + getTableName(info) + tableId
	indexName := "index_" + getTableName(info) + tableId

	dataStruct := reflect.Indirect(reflect.ValueOf(info))
	dataStructType := dataStruct.Type()
	ret := ""
	for i := 0; i < dataStructType.NumField(); i++ {
		fieldType := dataStructType.Field(i)
		fieldTag := fieldType.Tag

		if fieldTag.Get("index") != indexId {
			continue
		}

		if ret == "" {
			ret = fieldType.Name
		} else {
			ret += ", " + fieldType.Name
		}
		indexName += "_" + fieldType.Name
	}
	return "CREATE INDEX " + indexName + " ON " + tableName + " (" + ret + ")"
}

func GenerateInsertSql(info interface{}, tableId string) string {
	tableName := "tb_" + getTableName(info) + tableId
	ret := "INSERT INTO " + tableName + " VALUES("

	dataStruct := reflect.Indirect(reflect.ValueOf(info))
	dataStructType := dataStruct.Type()
	first := true
	for i := 0; i < dataStructType.NumField(); i++ {
		fieldType := dataStructType.Field(i)
		fieldValue := dataStruct.Field(i)
		fieldTag := fieldType.Tag

		if !first {
			ret += ","
		}

		str := ""
		switch fieldType.Type.Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			if fieldTag.Get("type") == "time" {
				str = "FROM_UNIXTIME(" + strconv.FormatInt(fieldValue.Int(), 10) + ")"
			} else {
				str = strconv.FormatInt(fieldValue.Int(), 10)
			}
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if fieldTag.Get("type") == "time" {
				str = "FROM_UNIXTIME(" + strconv.FormatUint(fieldValue.Uint(), 10) + ")"
			} else {
				str = strconv.FormatUint(fieldValue.Uint(), 10)
			}
		case reflect.Float32, reflect.Float64:
			str = strconv.FormatFloat(fieldValue.Float(), 'f', -1, 64)
		case reflect.String:
			str = "'" + fieldValue.String() + "'"
		case reflect.Bool:
			if fieldValue.Bool() {
				str = "1"
			} else {
				str = "0"
			}
		case reflect.Slice:
			if fieldType.Type.Elem().Kind() == reflect.Uint8 {
				str = string(fieldValue.Interface().([]byte))
			}
		}

		first = false
		ret += str
	}

	ret += ")"
	return ret
}

func GenerateUpdateSql(info interface{}, tableId string) string {
	tableName := "tb_" + getTableName(info) + tableId
	where := ""
	ret := "UPDATE " + tableName + " set "

	dataStruct := reflect.Indirect(reflect.ValueOf(info))
	dataStructType := dataStruct.Type()
	first := true
	for i := 0; i < dataStructType.NumField(); i++ {
		fieldType := dataStructType.Field(i)
		fieldValue := dataStruct.Field(i)
		fieldTag := fieldType.Tag

		pk := fieldTag.Get("pk") == "true"

		if !pk && !first {
			ret += ", "
		}

		str := ""
		switch fieldType.Type.Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			if fieldTag.Get("type") == "time" {
				str = fieldType.Name + "=" + "FROM_UNIXTIME(" + strconv.FormatInt(fieldValue.Int(), 10) + ")"
			} else {
				str = fieldType.Name + "=" + strconv.FormatInt(fieldValue.Int(), 10)
			}
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if fieldTag.Get("type") == "time" {
				str = fieldType.Name + "=" + "FROM_UNIXTIME(" + strconv.FormatUint(fieldValue.Uint(), 10) + ")"
			} else {
				str = fieldType.Name + "=" + strconv.FormatUint(fieldValue.Uint(), 10)
			}
		case reflect.Float32, reflect.Float64:
			str = fieldType.Name + "=" + strconv.FormatFloat(fieldValue.Float(), 'f', -1, 64)
		case reflect.String:
			str = fieldType.Name + "='" + fieldValue.String() + "'"
		case reflect.Bool:
			if fieldValue.Bool() {
				str = fieldType.Name + "=1"
			} else {
				str = fieldType.Name + "=0"
			}
		case reflect.Slice:
			if fieldType.Type.Elem().Kind() == reflect.Uint8 {
				str = fieldType.Name + string(fieldValue.Interface().([]byte))
			}
		}

		if pk {
			if where == "" {
				where = str
			} else {
				where += ", " + str
			}
		} else {
			first = false
			ret += str
		}
	}

	if where != "" {
		ret += " WHERE " + where
	}
	return ret
}

func GenerateSelectSql(info interface{}, tableId string) string {
	tableName := "tb_" + getTableName(info) + tableId
	where := ""
	ret := "SELECT "

	dataStruct := reflect.Indirect(reflect.ValueOf(info))
	dataStructType := dataStruct.Type()
	first := true
	for i := 0; i < dataStructType.NumField(); i++ {
		fieldType := dataStructType.Field(i)
		fieldValue := dataStruct.Field(i)
		fieldTag := fieldType.Tag

		pk := fieldTag.Get("pk") == "true"

		if pk {
			str := ""
			switch fieldType.Type.Kind() {
			case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
				if fieldTag.Get("type") == "time" {
					str = fieldType.Name + "=" + "FROM_UNIXTIME(" + strconv.FormatInt(fieldValue.Int(), 10) + ")"
				} else {
					str = fieldType.Name + "=" + strconv.FormatInt(fieldValue.Int(), 10)
				}
			case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				if fieldTag.Get("type") == "time" {
					str = fieldType.Name + "=" + "FROM_UNIXTIME(" + strconv.FormatUint(fieldValue.Uint(), 10) + ")"
				} else {
					str = fieldType.Name + "=" + strconv.FormatUint(fieldValue.Uint(), 10)
				}
			case reflect.Float32, reflect.Float64:
				str = fieldType.Name + "=" + strconv.FormatFloat(fieldValue.Float(), 'f', -1, 64)
			case reflect.String:
				str = fieldType.Name + "='" + fieldValue.String() + "'"
			case reflect.Bool:
				if fieldValue.Bool() {
					str = fieldType.Name + "=1"
				} else {
					str = fieldType.Name + "=0"
				}
			case reflect.Slice:
				if fieldType.Type.Elem().Kind() == reflect.Uint8 {
					str = fieldType.Name + string(fieldValue.Interface().([]byte))
				}
			}

			if where == "" {
				where = str
			} else {
				where += ", " + str
			}
		} else {
			if !first {
				ret += ", "
			}

			str := ""
			switch fieldType.Type.Kind() {
			case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				if fieldTag.Get("type") == "time" {
					str = "UNIX_TIMESTAMP(" + fieldType.Name + ")"
				} else {
					str = fieldType.Name
				}
			case reflect.Float32, reflect.Float64:
				str = fieldType.Name
			case reflect.String:
				str = fieldType.Name
			case reflect.Bool:
				str = fieldType.Name
			case reflect.Slice:
				str = fieldType.Name
			}
			first = false
			ret += str
		}
	}

	if where != "" {
		ret += " FROM " + tableName + " WHERE " + where
	}
	return ret
}

func GenerateSelectAllSql(info interface{}, tableId string) string {
	tableName := "tb_" + getTableName(info) + tableId
	ret := "SELECT "

	dataStruct := reflect.Indirect(reflect.ValueOf(info))
	dataStructType := dataStruct.Type()
	first := true
	for i := 0; i < dataStructType.NumField(); i++ {
		fieldType := dataStructType.Field(i)
		fieldTag := fieldType.Tag

		if !first {
			ret += ", "
		}

		str := ""
		switch fieldType.Type.Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if fieldTag.Get("type") == "time" {
				str = "UNIX_TIMESTAMP(" + fieldType.Name + ")"
			} else {
				str = fieldType.Name
			}
		case reflect.Float32, reflect.Float64:
			str = fieldType.Name
		case reflect.String:
			str = fieldType.Name
		case reflect.Bool:
			str = fieldType.Name
		case reflect.Slice:
			str = fieldType.Name
		}
		first = false
		ret += str
	}
	ret += " FROM " + tableName
	return ret
}

func GenerateSelectOneFieldSql(info interface{}, field string, tableId string) string {
	tableName := "tb_" + getTableName(info) + tableId
	where := ""
	ret := "SELECT " + field

	dataStruct := reflect.Indirect(reflect.ValueOf(info))
	dataStructType := dataStruct.Type()
	for i := 0; i < dataStructType.NumField(); i++ {
		fieldType := dataStructType.Field(i)
		fieldValue := dataStruct.Field(i)
		fieldTag := fieldType.Tag

		if fieldTag.Get("pk") != "true" {
			continue
		}

		str := ""
		switch fieldType.Type.Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			if fieldTag.Get("type") == "time" {
				str = fieldType.Name + "=" + "FROM_UNIXTIME(" + strconv.FormatInt(fieldValue.Int(), 10) + ")"
			} else {
				str = fieldType.Name + "=" + strconv.FormatInt(fieldValue.Int(), 10)
			}
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if fieldTag.Get("type") == "time" {
				str = fieldType.Name + "=" + "FROM_UNIXTIME(" + strconv.FormatUint(fieldValue.Uint(), 10) + ")"
			} else {
				str = fieldType.Name + "=" + strconv.FormatUint(fieldValue.Uint(), 10)
			}
		case reflect.Float32, reflect.Float64:
			str = fieldType.Name + "=" + strconv.FormatFloat(fieldValue.Float(), 'f', -1, 64)
		case reflect.String:
			str = fieldType.Name + "='" + fieldValue.String() + "'"
		case reflect.Bool:
			if fieldValue.Bool() {
				str = fieldType.Name + "=1"
			} else {
				str = fieldType.Name + "=0"
			}
		case reflect.Slice:
			if fieldType.Type.Elem().Kind() == reflect.Uint8 {
				str = fieldType.Name + string(fieldValue.Interface().([]byte))
			}
		}

		if where == "" {
			where = str
		} else {
			where += ", " + str
		}

	}

	if where != "" {
		ret += " FROM " + tableName + " WHERE " + where
	}
	return ret
}

func GenerateDeleteSql(info interface{}, tableId string) string {
	tableName := "tb_" + getTableName(info) + tableId
	where := ""
	ret := "DELETE FROM " + tableName

	dataStruct := reflect.Indirect(reflect.ValueOf(info))
	dataStructType := dataStruct.Type()
	for i := 0; i < dataStructType.NumField(); i++ {
		fieldType := dataStructType.Field(i)
		fieldValue := dataStruct.Field(i)
		fieldTag := fieldType.Tag

		if fieldTag.Get("pk") != "true" {
			continue
		}

		str := ""
		switch fieldType.Type.Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			if fieldTag.Get("type") == "time" {
				str = fieldType.Name + "=" + "FROM_UNIXTIME(" + strconv.FormatInt(fieldValue.Int(), 10) + ")"
			} else {
				str = fieldType.Name + "=" + strconv.FormatInt(fieldValue.Int(), 10)
			}
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if fieldTag.Get("type") == "time" {
				str = fieldType.Name + "=" + "FROM_UNIXTIME(" + strconv.FormatUint(fieldValue.Uint(), 10) + ")"
			} else {
				str = fieldType.Name + "=" + strconv.FormatUint(fieldValue.Uint(), 10)
			}
		case reflect.Float32, reflect.Float64:
			str = fieldType.Name + "=" + strconv.FormatFloat(fieldValue.Float(), 'f', -1, 64)
		case reflect.String:
			str = fieldType.Name + "='" + fieldValue.String() + "'"
		case reflect.Bool:
			if fieldValue.Bool() {
				str = fieldType.Name + "=1"
			} else {
				str = fieldType.Name + "=0"
			}
		case reflect.Slice:
			if fieldType.Type.Elem().Kind() == reflect.Uint8 {
				str = fieldType.Name + string(fieldValue.Interface().([]byte))
			}
		}

		if where == "" {
			where = str
		} else {
			where += ", " + str
		}
	}

	if where != "" {
		ret += " WHERE " + where
	}
	return ret
}
