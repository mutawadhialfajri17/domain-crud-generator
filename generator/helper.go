package generator

import (
	"domain-crud-generator/utils"
	"fmt"
	"os"
)

func createFile(path string) (f *os.File) {
	f, err := os.Create(path)
	if err != nil {
		utils.Print(err)
	}
	return
}

func generateStructValue(attrs []Attribute) string {
	result := ""
	for _, attr := range attrs {
		if attr.ColumnName == "" || attr.Type == "" {
			continue
		}
		// make attribute id to int64
		if attr.ColumnName == "id" {
			attr.Type = int64Str
		}
		jsonTag := "json:" + `"` + attr.ColumnName + `"`
		dbTag := "db:" + `"` + attr.ColumnName + `"`
		// replace id to uppercase (handle attribute which have suffix id)
		attributeName := utils.StringReplaceAll(utils.StringToCamelCase(attr.ColumnName), "Id", "ID")

		attributeDetail := fmt.Sprintf("%s %s `%s %s`", attributeName, convertType(attr.Type), jsonTag, dbTag)
		result += attributeDetail + "\n"
	}
	return result
}

func generateAppNameAcronym(appName string) string {
	// remove "engine"
	appName = utils.StringReplace(appName, "-engine", "", 1)

	splits := utils.StringSplit(appName, "-")
	appName = utils.StringJoin(splits, " ")
	return utils.StringToAcronym(appName)
}

func generateQueryFromAttributes(attrs []Attribute, tableName string) string {
	query := ""
	for i, attr := range attrs {
		// dont use coalesce at id
		if attr.ColumnName == "id" {
			query += fmt.Sprintf("%s.%s AS %s", tableName, attr.ColumnName, attr.ColumnName)
		} else {
			nilDataDB := mapNilDataDB[convertType(attr.Type)]
			query += fmt.Sprintf(`COALESCE(%s.%s, %s) AS "%s"`, tableName, attr.ColumnName, nilDataDB, attr.ColumnName)
		}

		// dont add "," at last index
		if i == len(attrs)-1 {
			break
		}
		query += ",\n"
	}
	return query
}

func generateParamNumberList(start, lenParam int) string {
	result := ""
	for i := start; i < start+lenParam; i++ {
		result = result + "$" + utils.IntToString(i) + ","
	}

	if result == "" {
		return ""
	}

	result = result[:len(result)-1]

	return result
}

func generateSetAttribute(idx int, columnName string) string {
	return fmt.Sprintf("%s = $%d", columnName, idx+1)
}

func replaceDashWithUnderscore(s string) string {
	return utils.StringReplaceAll(s, "-", "_")
}

func writePackageName(f *os.File, domainName string) {
	packageStr := utils.StringReplaceAll(packageCode, domainNameLowerCaseReplacer, utils.StringRemoveSpecialCharacter(domainName))
	f.Write(convertToByteln(packageStr))
}

func generateScanGetAttribute(attrs []Attribute) string {
	result := ""
	for _, attr := range attrs {
		// replace id to uppercase
		attributeName := utils.StringReplaceAll(utils.StringToCamelCase(attr.ColumnName), "Id", "ID")

		var attributeResult string
		if mapIsDataArray[convertType(attr.Type)] {
			attributeResult += fmt.Sprintf("pq.Array(&data.%s)", attributeName)
		} else {
			attributeResult += fmt.Sprintf("&data.%s", attributeName)
		}
		result += attributeResult + ",\n"
	}
	return result
}

func generateScanInsertAttribute(attrs []Attribute) string {
	result := ""
	for _, attr := range attrs {
		// ignore attribute id
		if attr.ColumnName == "id" {
			continue
		}

		// replace id to uppercase (handle attribute which have suffix id)
		attributeName := utils.StringReplaceAll(utils.StringToCamelCase(attr.ColumnName), "Id", "ID")

		var attributeResult string
		if mapIsDataArray[convertType(attr.Type)] {
			attributeResult += fmt.Sprintf("pq.Array(param.%s)", attributeName)
		} else {
			attributeResult += fmt.Sprintf("param.%s", attributeName)
		}
		result += attributeResult + ",\n"
	}
	return result
}

func generateScanUpdateAttribute(attrs []Attribute) string {
	result := ""
	for _, attr := range attrs {
		// replace id to uppercase
		attributeName := utils.StringReplaceAll(utils.StringToCamelCase(attr.ColumnName), "Id", "ID")

		var attributeResult string
		if mapIsDataArray[convertType(attr.Type)] {
			attributeResult += fmt.Sprintf("pq.Array(param.%s)", attributeName)
		} else {
			attributeResult += fmt.Sprintf("param.%s", attributeName)
		}
		result += attributeResult + ",\n"
	}
	return result
}

func generateCopyAttribute(attrs []Attribute) string {
	result := ""
	// copy slice first
	for _, attr := range attrs {
		attributeName := utils.StringReplaceAll(utils.StringToCamelCase(attr.ColumnName), "Id", "ID")
		if mapIsDataArray[convertType(attr.Type)] {
			result += fmt.Sprintf(`
			copy%s := make(%s, len([domainName].%s))
			copy(copy%s, [domainName].%s)
			`, attributeName, convertType(attr.Type), attributeName, attributeName, attributeName)
		}
	}

	result += "copyResult[i] = [DomainName]{\n"
	// copy all attributes
	for _, attr := range attrs {
		attributeName := utils.StringReplaceAll(utils.StringToCamelCase(attr.ColumnName), "Id", "ID")
		if mapIsDataArray[convertType(attr.Type)] {
			result += fmt.Sprintf("%s: copy%s,\n", attributeName, attributeName)
		} else {
			result += fmt.Sprintf("%s: [domainName].%s,\n", attributeName, attributeName)
		}
	}
	result += "}\n"

	return result
}
