package generator

import "strings"

func convertType(t string) string {
	switch strings.ToLower(t) {
	case "int", "integer":
		return intStr
	case "int32":
		return int32Str
	case "int64", "bigint", "int8", "unix":
		return int64Str
	case "string", "varchar", "text":
		return stringStr
	case "bool", "boolean":
		return boolStr
	case "float32":
		return float32Str
	case "float64", "double", "float":
		return float64Str
	case "[]int", "int[]":
		return arrayOfIntStr
	case "[]int32":
		return arrayOfInt32Str
	case "[]int64", "bigint[]", "int8[]":
		return arrayOfInt64Str
	case "[]string", "varchar[]":
		return arrayOfStringStr
	case "[]float32":
		return arrayOfFloat32Str
	case "[]float64", "double[]":
		return arrayOfFloat64Str
	case "time", "time.Time", "timestamp", "date":
		return timeStr
	default:
		return ""
	}
}

func convertToByte(s string) []byte {
	return []byte(s)
}

func convertToByteln(s string) []byte {
	return []byte(s + "\n")
}

func convertToBytelnln(s string) []byte {
	return []byte(s + "\n\n")
}
