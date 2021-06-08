package helpers

import (
	"github.com/google/uuid"
	"strconv"
)

const TYPE_INT = "int"
const TYPE_FlOAT = "float"
const TYPE_BOOL = "bool"
const TYPE_STRING = "string"
const TYPE_UUID = "uuid"

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
	_, err = uuid.Parse(value)
	if err == nil {
		return TYPE_UUID
	}
	return TYPE_STRING
}
