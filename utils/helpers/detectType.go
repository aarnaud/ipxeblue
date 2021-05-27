package helpers

import "strconv"

const TYPE_INT = "int"
const TYPE_FlOAT = "float"
const TYPE_BOOL = "bool"
const TYPE_STRING = "string"

func StringToType(value string) string {
	_, err := strconv.ParseInt(value, 10, 64)
	if err == nil {
		return TYPE_INT
	}
	_, err = strconv.ParseFloat(value, 64)
	if err == nil {
		return TYPE_FlOAT
	}
	_, err = strconv.ParseBool(value)
	if err == nil {
		return TYPE_BOOL
	}
	return TYPE_STRING
}
