package generator

import (
	"domain-crud-generator/utils"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

func GenerateInit(path string, requestData RequestData) {
	f := createFile(path)
	defer f.Close()

	// package
	writePackageName(f, requestData.DomainName)

	// import
	importCode := importInitCode
	strs := utils.StringSplit(importInitCode, "\n\n")
	prefix := strs[0]
	suffix := strs[1]
	importCode = prefix + "\n" + resourceTDKLib
	if requestData.Options.IsUseRedisAndMemcache {
		importCode += "\n" + redisLib
		importCode += "\n" + memcacheLib
		importCode += "\n" + cacheLib
	}
	importCode += "\n" + suffix

	importCode = utils.StringReplaceAll(importCode, appNameReplacer, requestData.AppName)
	f.Write(convertToBytelnln(importCode))

	// struct and interface
	f.Write(convertToByteln("type ("))

	f.Write(convertToByteln(utils.StringReplaceAll(resourceStruct, domainNameCamelCaseReplacer, utils.StringToCamelCase(requestData.DomainName)) + "{"))
	f.Write(convertToByteln(dbResourceStruct))
	if requestData.Options.IsUseRedisAndMemcache {
		f.Write(convertToByteln(redisResourceStruct))
		f.Write(convertToByteln(memcacheResourceStruct))
	}
	f.Write(convertToBytelnln("}"))

	f.Write(convertToByteln(utils.StringReplaceAll(domainStruct, domainNameCamelCaseReplacer, utils.StringToCamelCase(requestData.DomainName)) + "{"))
	f.Write(convertToByteln(utils.StringReplaceAll(resourceDomainStruct, domainNameCamelCaseReplacer, utils.StringToCamelCase(requestData.DomainName))))
	f.Write(convertToBytelnln("}"))

	f.Write(convertToByteln(utils.StringReplaceAll(domainInterface, domainNameCamelCaseReplacer, utils.StringToCamelCase(requestData.DomainName)) + "{"))
	domainItfData := domainInterfaceData
	domainItfData = utils.StringReplaceAll(domainItfData, domainNameCamelCaseReplacer, utils.StringToCamelCase(requestData.DomainName))
	domainItfData = utils.StringReplaceAll(domainItfData, domainNameLowerCamelCaseReplacer, utils.StringToLowerCamel(requestData.DomainName))
	if requestData.Options.IsUseRedisAndMemcache {
		domainItfData = utils.StringReplaceAll(domainItfData, useCacheParamReplacer, ", "+cacheParam)
	} else {
		domainItfData = utils.StringReplaceAll(domainItfData, useCacheParamReplacer, "")
	}
	if requestData.Options.IsUseSingleflight {
		singleflightDomainItfData := singleflightDomainInterfaceData
		singleflightDomainItfData = utils.StringReplaceAll(singleflightDomainItfData, domainNameCamelCaseReplacer, utils.StringToCamelCase(requestData.DomainName))
		singleflightDomainItfData = utils.StringReplaceAll(singleflightDomainItfData, domainNameLowerCamelCaseReplacer, utils.StringToLowerCamel(requestData.DomainName))
		if requestData.Options.IsUseRedisAndMemcache {
			singleflightDomainItfData = utils.StringReplaceAll(singleflightDomainItfData, useCacheParamReplacer, ", "+cacheParam)
		} else {
			singleflightDomainItfData = utils.StringReplaceAll(singleflightDomainItfData, useCacheParamReplacer, "")
		}
		domainItfData += singleflightDomainItfData
	}
	f.Write(convertToByteln(domainItfData))
	f.Write(convertToBytelnln("}"))

	f.Write(convertToByteln(utils.StringReplaceAll(resourceInterface, domainNameCamelCaseReplacer, utils.StringToCamelCase(requestData.DomainName)) + "{"))
	resourceItfData := resourceInterfaceDataDB
	resourceItfData = utils.StringReplaceAll(resourceItfData, domainNameCamelCaseReplacer, utils.StringToCamelCase(requestData.DomainName))
	resourceItfData = utils.StringReplaceAll(resourceItfData, domainNameLowerCamelCaseReplacer, utils.StringToLowerCamel(requestData.DomainName))
	if requestData.Options.IsUseRedisAndMemcache {
		additionalResourceItfData := resourceInterfaceDataCache
		additionalResourceItfData = utils.StringReplaceAll(additionalResourceItfData, domainNameCamelCaseReplacer, utils.StringToCamelCase(requestData.DomainName))
		additionalResourceItfData = utils.StringReplaceAll(additionalResourceItfData, domainNameLowerCamelCaseReplacer, utils.StringToLowerCamel(requestData.DomainName))
		resourceItfData += additionalResourceItfData

	}
	f.Write(convertToByteln(resourceItfData))
	f.Write(convertToByteln("}"))

	f.Write(convertToBytelnln(")"))

	// func
	initCode := funcInitDomainCode
	initCode = utils.StringReplaceAll(initCode, domainNameCamelCaseReplacer, utils.StringToCamelCase(requestData.DomainName))
	f.Write(convertToByteln(initCode))

	for _, component := range crudComponent {
		switch component {
		case "CREATE":
			funcCode := funcInsertDomainCode
			funcCode = utils.StringReplaceAll(funcCode, domainNameCamelCaseReplacer, utils.StringToCamelCase(requestData.DomainName))
			funcCode = utils.StringReplaceAll(funcCode, domainNameLowerCaseReplacer, utils.StringRemoveSpecialCharacter(requestData.DomainName))
			funcCode = utils.StringReplaceAll(funcCode, domainNameLowerCamelCaseReplacer, utils.StringToLowerCamel(requestData.DomainName))
			f.Write(convertToByteln(funcCode))
		case "READ":
			var funcCode string
			if requestData.Options.IsUseRedisAndMemcache {
				funcCode = funcGetDomainWithCacheCode
			} else {
				funcCode = funcGetDomainWithoutCacheCode
			}
			funcCode = utils.StringReplaceAll(funcCode, domainNameCamelCaseReplacer, utils.StringToCamelCase(requestData.DomainName))
			funcCode = utils.StringReplaceAll(funcCode, domainNameLowerCaseReplacer, utils.StringRemoveSpecialCharacter(requestData.DomainName))
			funcCode = utils.StringReplaceAll(funcCode, domainNameLowerCamelCaseReplacer, utils.StringToLowerCamel(requestData.DomainName))
			f.Write(convertToByteln(funcCode))
		case "UPDATE":
			funcCode := funcUpdateDomainCode
			funcCode = utils.StringReplaceAll(funcCode, domainNameCamelCaseReplacer, utils.StringToCamelCase(requestData.DomainName))
			funcCode = utils.StringReplaceAll(funcCode, domainNameLowerCaseReplacer, utils.StringRemoveSpecialCharacter(requestData.DomainName))
			funcCode = utils.StringReplaceAll(funcCode, domainNameLowerCamelCaseReplacer, utils.StringToLowerCamel(requestData.DomainName))
			f.Write(convertToByteln(funcCode))
		case "DELETE":
			funcCode := funcDeleteDomainCode
			funcCode = utils.StringReplaceAll(funcCode, domainNameCamelCaseReplacer, utils.StringToCamelCase(requestData.DomainName))
			funcCode = utils.StringReplaceAll(funcCode, domainNameLowerCaseReplacer, utils.StringRemoveSpecialCharacter(requestData.DomainName))
			funcCode = utils.StringReplaceAll(funcCode, domainNameLowerCamelCaseReplacer, utils.StringToLowerCamel(requestData.DomainName))
			f.Write(convertToByteln(funcCode))
		}
	}
}

func GenerateType(path string, requestData RequestData) {
	f := createFile(path)
	defer f.Close()

	// package
	writePackageName(f, requestData.DomainName)

	// import
	isNeedImportTime := false
	for _, reqAttribute := range requestData.Attributes {
		if convertType(reqAttribute.Type) == timeStr {
			isNeedImportTime = true
			break
		}
	}
	if isNeedImportTime {
		f.Write(convertToByteln(`import "time"`))
	}

	// struct
	for _, component := range crudComponent {
		switch component {
		case "CREATE":
			structCode := structCreate
			structCode = utils.StringReplaceAll(structCode, domainNameCamelCaseReplacer, utils.StringToCamelCase(requestData.DomainName))

			// remove id from attributes
			fixAttributes := []Attribute{}
			for _, attr := range requestData.Attributes {
				if attr.ColumnName == "id" {
					continue
				}
				fixAttributes = append(fixAttributes, attr)
			}
			structCode = utils.StringReplaceAll(structCode, attributeReplacer, generateStructValue(fixAttributes))
			f.Write(convertToByteln(structCode))
		case "READ":
			structCode := structRead
			structCode = utils.StringReplaceAll(structCode, domainNameCamelCaseReplacer, utils.StringToCamelCase(requestData.DomainName))
			structCode = utils.StringReplaceAll(structCode, attributeReplacer, generateStructValue(requestData.Attributes))
			f.Write(convertToByteln(structCode))
		case "UPDATE":
			structCode := structUpdate
			structCode = utils.StringReplaceAll(structCode, domainNameCamelCaseReplacer, utils.StringToCamelCase(requestData.DomainName))
			structCode = utils.StringReplaceAll(structCode, attributeReplacer, generateStructValue(requestData.Attributes))
			f.Write(convertToByteln(structCode))
		case "DELETE":
			structCode := structDelete
			structCode = utils.StringReplaceAll(structCode, domainNameCamelCaseReplacer, utils.StringToCamelCase(requestData.DomainName))

			// only retrieve id
			fixAttributes := []Attribute{}
			for _, attr := range requestData.Attributes {
				if attr.ColumnName == "id" {
					fixAttributes = append(fixAttributes, attr)
					break
				}
			}
			structCode = utils.StringReplaceAll(structCode, attributeReplacer, generateStructValue(fixAttributes))
			f.Write(convertToByteln(structCode))
		}
	}

	if requestData.Options.IsUseSingleflight {
		structCode := structSingleflight
		structCode = utils.StringReplaceAll(structCode, domainNameCamelCaseReplacer, utils.StringToCamelCase(requestData.DomainName))
		f.Write(convertToByteln(structCode))
	}
}

func GenerateConst(path string, requestData RequestData) {
	f := createFile(path)
	defer f.Close()

	// package
	writePackageName(f, requestData.DomainName)

	// const
	if requestData.Options.IsUseRedisAndMemcache {
		cacheCode := constCache
		cacheCode = utils.StringReplaceAll(cacheCode, domainNameCamelCaseReplacer, utils.StringToCamelCase(requestData.DomainName))
		cacheCode = utils.StringReplaceAll(cacheCode, domainNameLowerCaseReplacer, utils.StringRemoveSpecialCharacter(requestData.DomainName))
		cacheCode = utils.StringReplaceAll(cacheCode, acronymReplacer, generateAppNameAcronym(requestData.AppName))
		cacheCode = utils.StringReplaceAll(cacheCode, cacheTTLReplacer, utils.IntToString(requestData.RedisTTL))
		cacheCode = utils.StringReplaceAll(cacheCode, memcacheTTLReplacer, utils.IntToString(requestData.MemcacheTTL))
		f.Write(convertToByteln(cacheCode))
	}

	f.Write(convertToByteln("const ("))
	for _, component := range crudComponent {
		switch component {
		case "CREATE":
			createQuery := constCreateQuery
			createQuery = utils.StringReplaceAll(createQuery, domainNameCamelCaseReplacer, utils.StringToCamelCase(requestData.DomainName))
			createQuery = utils.StringReplaceAll(createQuery, tableNameReplacer, requestData.TableName)
			attrList := []string{}
			for _, attr := range requestData.Attributes {
				attrList = append(attrList, attr.ColumnName)
			}
			createQuery = utils.StringReplaceAll(createQuery, attributesListReplacer, utils.StringJoin(attrList, ","))
			createQuery = utils.StringReplaceAll(createQuery, paramNumberListReplacer, generateParamNumberList(1, len(requestData.Attributes)))
			f.Write(convertToByteln(createQuery))
		case "READ":
			readQuery := constReadQuery
			readQuery = utils.StringReplaceAll(readQuery, domainNameCamelCaseReplacer, utils.StringToCamelCase(requestData.DomainName))
			readQuery = utils.StringReplaceAll(readQuery, tableNameReplacer, requestData.TableName)
			readQuery = utils.StringReplaceAll(readQuery, attributeQueryReplacer, generateQueryFromAttributes(requestData.Attributes, requestData.TableName))
			f.Write(convertToByteln(readQuery))
		case "UPDATE":
			updateQuery := constUpdateQuery
			updateQuery = utils.StringReplaceAll(updateQuery, domainNameCamelCaseReplacer, utils.StringToCamelCase(requestData.DomainName))
			updateQuery = utils.StringReplaceAll(updateQuery, tableNameReplacer, requestData.TableName)
			attributeSetList := []string{}
			for i, attr := range requestData.Attributes {
				// set id in query, and dont append to list
				if attr.ColumnName == "id" {
					updateQuery = utils.StringReplaceAll(updateQuery, idSetReplacer, generateSetAttribute(i, attr.ColumnName))
				} else {
					attributeSetList = append(attributeSetList, generateSetAttribute(i, attr.ColumnName))
				}
			}
			updateQuery = utils.StringReplaceAll(updateQuery, attributeSetListReplacer, utils.StringJoin(attributeSetList, ",\n"))
			f.Write(convertToByteln(updateQuery))
		case "DELETE":
			deleteQuery := constDeleteQuery
			deleteQuery = utils.StringReplaceAll(deleteQuery, domainNameCamelCaseReplacer, utils.StringToCamelCase(requestData.DomainName))
			deleteQuery = utils.StringReplaceAll(deleteQuery, tableNameReplacer, requestData.TableName)
			for i, attr := range requestData.Attributes {
				if attr.ColumnName == "id" {
					deleteQuery = utils.StringReplaceAll(deleteQuery, idSetReplacer, generateSetAttribute(i, attr.ColumnName))
					break
				}
			}
			f.Write(convertToByteln(deleteQuery))
		}
	}
	f.Write(convertToByteln(")"))

	if requestData.Options.IsUseSingleflight {
		singleflightCode := constSingleflight
		singleflightCode = utils.StringReplaceAll(singleflightCode, domainNameCamelCaseReplacer, utils.StringToCamelCase(requestData.DomainName))
		singleflightCode = utils.StringReplaceAll(singleflightCode, domainNameLowerCaseReplacer, utils.StringRemoveSpecialCharacter(requestData.DomainName))
		f.Write(convertToByteln(singleflightCode))
	}
}

func GenerateMock(path string, requestData RequestData) {
	source := fmt.Sprintf("-source=%s", utils.GenerateResultPath(fmt.Sprintf("%s.go", requestData.DomainName), requestData.DomainName))
	dest := fmt.Sprintf("-destination=%s", path)
	cmd := exec.Command("mockgen", source, dest)
	_, err := cmd.Output()
	if err != nil {
		utils.Print(err)
	}

	err = filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		fileContent, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// replace package
		oldPackage := fmt.Sprintf("package mock_%s", utils.StringRemoveSpecialCharacter(requestData.DomainName))
		newPackage := fmt.Sprintf("package %s", utils.StringRemoveSpecialCharacter(requestData.DomainName))
		fileContent = []byte(utils.StringReplace(string(fileContent), oldPackage, newPackage, 1))

		// remove self import
		selfImportPattern := fmt.Sprintf(`%s "%s/result/%s"`, replaceDashWithUnderscore(requestData.DomainName), utils.GeneratorName, requestData.DomainName)
		fileContent = []byte(utils.StringReplace(string(fileContent), selfImportPattern, "", 1))

		// remove all self import reference
		selfImportReferencePattern := fmt.Sprintf("%s.", replaceDashWithUnderscore(requestData.DomainName))
		fileContent = []byte(utils.StringReplaceAll(string(fileContent), selfImportReferencePattern, ""))

		err = os.WriteFile(path, fileContent, 0)
		if err != nil {
			panic(err)
		}

		return nil
	})
	if err != nil {
		utils.Print(err)
	}

}

func GenerateDatabase(path string, requestData RequestData) {
	f := createFile(path)
	defer f.Close()

	// package
	writePackageName(f, requestData.DomainName)

	// import
	importCode := importDatabaseCode
	importCode = utils.StringReplaceAll(importCode, appNameReplacer, requestData.AppName)
	f.Write(convertToBytelnln(importCode))

	// all func
	funcCode := funcDatabaseCode
	funcCode = utils.StringReplaceAll(funcCode, domainNameCamelCaseReplacer, utils.StringToCamelCase(requestData.DomainName))
	funcCode = utils.StringReplaceAll(funcCode, domainNameLowerCamelCaseReplacer, utils.StringToLowerCamel(requestData.DomainName))
	funcCode = utils.StringReplaceAll(funcCode, domainNameLowerCaseReplacer, utils.StringRemoveSpecialCharacter(requestData.DomainName))
	funcCode = utils.StringReplaceAll(funcCode, scanGetAttributeReplacer, generateScanGetAttribute(requestData.Attributes))
	funcCode = utils.StringReplaceAll(funcCode, scanInsertAttributeReplacer, generateScanInsertAttribute(requestData.Attributes))
	funcCode = utils.StringReplaceAll(funcCode, scanUpdateAttributeReplacer, generateScanUpdateAttribute(requestData.Attributes))
	f.Write(convertToBytelnln(funcCode))
}

func GenerateCache(path string, requestData RequestData) {
	f := createFile(path)
	defer f.Close()

	// package
	writePackageName(f, requestData.DomainName)

	// import
	importCode := importCacheCode
	importCode = utils.StringReplaceAll(importCode, appNameReplacer, requestData.AppName)
	f.Write(convertToBytelnln(importCode))

	// all func
	funcCode := funcCacheCode
	funcCode = utils.StringReplaceAll(funcCode, domainNameCamelCaseReplacer, utils.StringToCamelCase(requestData.DomainName))
	funcCode = utils.StringReplaceAll(funcCode, domainNameLowerCamelCaseReplacer, utils.StringToLowerCamel(requestData.DomainName))
	funcCode = utils.StringReplaceAll(funcCode, domainNameLowerCaseReplacer, utils.StringRemoveSpecialCharacter(requestData.DomainName))
	funcCode = utils.StringReplaceAll(funcCode, scanGetAttributeReplacer, generateScanGetAttribute(requestData.Attributes))
	funcCode = utils.StringReplaceAll(funcCode, scanInsertAttributeReplacer, generateScanInsertAttribute(requestData.Attributes))
	funcCode = utils.StringReplaceAll(funcCode, scanUpdateAttributeReplacer, generateScanUpdateAttribute(requestData.Attributes))
	f.Write(convertToBytelnln(funcCode))
}

func GenerateSingleflight(path string, requestData RequestData) {
	f := createFile(path)
	defer f.Close()

	// package
	writePackageName(f, requestData.DomainName)

	// import
	importCode := importSingleflightCode
	importCode = utils.StringReplaceAll(importCode, appNameReplacer, requestData.AppName)
	f.Write(convertToBytelnln(importCode))

	// all func
	funcCode := funcSingleflightCode
	funcCode = utils.StringReplaceAll(funcCode, copyAttributeReplacer, generateCopyAttribute(requestData.Attributes))
	funcCode = utils.StringReplaceAll(funcCode, domainNameCamelCaseReplacer, utils.StringToCamelCase(requestData.DomainName))
	funcCode = utils.StringReplaceAll(funcCode, acronymReplacer, generateAppNameAcronym(requestData.AppName))
	funcCode = utils.StringReplaceAll(funcCode, domainNameLowerCamelCaseReplacer, utils.StringToLowerCamel(requestData.DomainName))
	funcCode = utils.StringReplaceAll(funcCode, domainNameLowerCaseReplacer, utils.StringRemoveSpecialCharacter(requestData.DomainName))
	f.Write(convertToBytelnln(funcCode))
}
