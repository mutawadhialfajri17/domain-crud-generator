package generator

import (
	"domain-crud-generator/utils"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gobeam/stringy"
	"github.com/iancoleman/strcase"
)

func GenerateInit(path string, requestData RequestData) {
	f := createFile(path)
	defer f.Close()

	// package
	writePackageName(f, requestData.DomainName)

	// import
	importCode := importInitCode
	strs := strings.Split(importInitCode, "\n\n")
	prefix := strs[0]
	suffix := strs[1]
	importCode = prefix + "\n" + resourceTDKLib
	if requestData.Options.IsUseRedisAndMemcache {
		importCode += "\n" + redisLib
		importCode += "\n" + memcacheLib
		importCode += "\n" + cacheLib
	}
	importCode += "\n" + suffix

	importCode = strings.ReplaceAll(importCode, appNameReplacer, requestData.AppName)
	f.Write(convertToBytelnln(importCode))

	// struct and interface
	f.Write(convertToByteln("type ("))

	f.Write(convertToByteln(strings.ReplaceAll(resourceStruct, domainNameCamelCaseReplacer, stringy.New(requestData.DomainName).CamelCase()) + "{"))
	f.Write(convertToByteln(dbResourceStruct))
	if requestData.Options.IsUseRedisAndMemcache {
		f.Write(convertToByteln(redisResourceStruct))
		f.Write(convertToByteln(memcacheResourceStruct))
	}
	f.Write(convertToBytelnln("}"))

	f.Write(convertToByteln(strings.ReplaceAll(domainStruct, domainNameCamelCaseReplacer, stringy.New(requestData.DomainName).CamelCase()) + "{"))
	f.Write(convertToByteln(strings.ReplaceAll(resourceDomainStruct, domainNameCamelCaseReplacer, stringy.New(requestData.DomainName).CamelCase())))
	f.Write(convertToBytelnln("}"))

	f.Write(convertToByteln(strings.ReplaceAll(domainInterface, domainNameCamelCaseReplacer, stringy.New(requestData.DomainName).CamelCase()) + "{"))
	domainItfData := domainInterfaceData
	domainItfData = strings.ReplaceAll(domainItfData, domainNameCamelCaseReplacer, stringy.New(requestData.DomainName).CamelCase())
	domainItfData = strings.ReplaceAll(domainItfData, domainNameLowerCamelCaseReplacer, strcase.ToLowerCamel(requestData.DomainName))
	if requestData.Options.IsUseRedisAndMemcache {
		domainItfData = strings.ReplaceAll(domainItfData, useCacheParamReplacer, ", "+cacheParam)
	} else {
		domainItfData = strings.ReplaceAll(domainItfData, useCacheParamReplacer, "")
	}
	if requestData.Options.IsUseSingleflight {
		singleflightDomainItfData := singleflightDomainInterfaceData
		singleflightDomainItfData = strings.ReplaceAll(singleflightDomainItfData, domainNameCamelCaseReplacer, stringy.New(requestData.DomainName).CamelCase())
		singleflightDomainItfData = strings.ReplaceAll(singleflightDomainItfData, domainNameLowerCamelCaseReplacer, strcase.ToLowerCamel(requestData.DomainName))
		if requestData.Options.IsUseRedisAndMemcache {
			singleflightDomainItfData = strings.ReplaceAll(singleflightDomainItfData, useCacheParamReplacer, ", "+cacheParam)
		} else {
			singleflightDomainItfData = strings.ReplaceAll(singleflightDomainItfData, useCacheParamReplacer, "")
		}
		domainItfData += singleflightDomainItfData
	}
	f.Write(convertToByteln(domainItfData))
	f.Write(convertToBytelnln("}"))

	f.Write(convertToByteln(strings.ReplaceAll(resourceInterface, domainNameCamelCaseReplacer, stringy.New(requestData.DomainName).CamelCase()) + "{"))
	resourceItfData := resourceInterfaceDataDB
	resourceItfData = strings.ReplaceAll(resourceItfData, domainNameCamelCaseReplacer, stringy.New(requestData.DomainName).CamelCase())
	resourceItfData = strings.ReplaceAll(resourceItfData, domainNameLowerCamelCaseReplacer, strcase.ToLowerCamel(requestData.DomainName))
	if requestData.Options.IsUseRedisAndMemcache {
		additionalResourceItfData := resourceInterfaceDataCache
		additionalResourceItfData = strings.ReplaceAll(additionalResourceItfData, domainNameCamelCaseReplacer, stringy.New(requestData.DomainName).CamelCase())
		additionalResourceItfData = strings.ReplaceAll(additionalResourceItfData, domainNameLowerCamelCaseReplacer, strcase.ToLowerCamel(requestData.DomainName))
		resourceItfData += additionalResourceItfData

	}
	f.Write(convertToByteln(resourceItfData))
	f.Write(convertToByteln("}"))

	f.Write(convertToBytelnln(")"))

	// func
	initCode := funcInitDomainCode
	initCode = strings.ReplaceAll(initCode, domainNameCamelCaseReplacer, stringy.New(requestData.DomainName).CamelCase())
	f.Write(convertToByteln(initCode))

	for _, component := range crudComponent {
		switch component {
		case "CREATE":
			funcCode := funcInsertDomainCode
			funcCode = strings.ReplaceAll(funcCode, domainNameCamelCaseReplacer, stringy.New(requestData.DomainName).CamelCase())
			funcCode = strings.ReplaceAll(funcCode, domainNameLowerCaseReplacer, stringy.New(requestData.DomainName).RemoveSpecialCharacter())
			funcCode = strings.ReplaceAll(funcCode, domainNameLowerCamelCaseReplacer, strcase.ToLowerCamel(requestData.DomainName))
			f.Write(convertToByteln(funcCode))
		case "READ":
			var funcCode string
			if requestData.Options.IsUseRedisAndMemcache {
				funcCode = funcGetDomainWithCacheCode
			} else {
				funcCode = funcGetDomainWithoutCacheCode
			}
			funcCode = strings.ReplaceAll(funcCode, domainNameCamelCaseReplacer, stringy.New(requestData.DomainName).CamelCase())
			funcCode = strings.ReplaceAll(funcCode, domainNameLowerCaseReplacer, stringy.New(requestData.DomainName).RemoveSpecialCharacter())
			funcCode = strings.ReplaceAll(funcCode, domainNameLowerCamelCaseReplacer, strcase.ToLowerCamel(requestData.DomainName))
			f.Write(convertToByteln(funcCode))
		case "UPDATE":
			funcCode := funcUpdateDomainCode
			funcCode = strings.ReplaceAll(funcCode, domainNameCamelCaseReplacer, stringy.New(requestData.DomainName).CamelCase())
			funcCode = strings.ReplaceAll(funcCode, domainNameLowerCaseReplacer, stringy.New(requestData.DomainName).RemoveSpecialCharacter())
			funcCode = strings.ReplaceAll(funcCode, domainNameLowerCamelCaseReplacer, strcase.ToLowerCamel(requestData.DomainName))
			f.Write(convertToByteln(funcCode))
		case "DELETE":
			funcCode := funcDeleteDomainCode
			funcCode = strings.ReplaceAll(funcCode, domainNameCamelCaseReplacer, stringy.New(requestData.DomainName).CamelCase())
			funcCode = strings.ReplaceAll(funcCode, domainNameLowerCaseReplacer, stringy.New(requestData.DomainName).RemoveSpecialCharacter())
			funcCode = strings.ReplaceAll(funcCode, domainNameLowerCamelCaseReplacer, strcase.ToLowerCamel(requestData.DomainName))
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
			structCode = strings.ReplaceAll(structCode, domainNameCamelCaseReplacer, stringy.New(requestData.DomainName).CamelCase())

			// remove id from attributes
			fixAttributes := []Attribute{}
			for _, attr := range requestData.Attributes {
				if attr.ColumnName == "id" {
					continue
				}
				fixAttributes = append(fixAttributes, attr)
			}
			structCode = strings.ReplaceAll(structCode, attributeReplacer, generateStructValue(fixAttributes))
			f.Write(convertToByteln(structCode))
		case "READ":
			structCode := structRead
			structCode = strings.ReplaceAll(structCode, domainNameCamelCaseReplacer, stringy.New(requestData.DomainName).CamelCase())
			structCode = strings.ReplaceAll(structCode, attributeReplacer, generateStructValue(requestData.Attributes))
			f.Write(convertToByteln(structCode))
		case "UPDATE":
			structCode := structUpdate
			structCode = strings.ReplaceAll(structCode, domainNameCamelCaseReplacer, stringy.New(requestData.DomainName).CamelCase())
			structCode = strings.ReplaceAll(structCode, attributeReplacer, generateStructValue(requestData.Attributes))
			f.Write(convertToByteln(structCode))
		case "DELETE":
			structCode := structDelete
			structCode = strings.ReplaceAll(structCode, domainNameCamelCaseReplacer, stringy.New(requestData.DomainName).CamelCase())

			// only retrieve id
			fixAttributes := []Attribute{}
			for _, attr := range requestData.Attributes {
				if attr.ColumnName == "id" {
					fixAttributes = append(fixAttributes, attr)
					break
				}
			}
			structCode = strings.ReplaceAll(structCode, attributeReplacer, generateStructValue(fixAttributes))
			f.Write(convertToByteln(structCode))
		}
	}

	if requestData.Options.IsUseSingleflight {
		structCode := structSingleflight
		structCode = strings.ReplaceAll(structCode, domainNameCamelCaseReplacer, stringy.New(requestData.DomainName).CamelCase())
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
		cacheCode = strings.ReplaceAll(cacheCode, domainNameCamelCaseReplacer, stringy.New(requestData.DomainName).CamelCase())
		cacheCode = strings.ReplaceAll(cacheCode, domainNameLowerCaseReplacer, stringy.New(requestData.DomainName).RemoveSpecialCharacter())
		cacheCode = strings.ReplaceAll(cacheCode, acronymReplacer, generateAppNameAcronym(requestData.AppName))
		cacheCode = strings.ReplaceAll(cacheCode, cacheTTLReplacer, strconv.Itoa(requestData.RedisTTL))
		cacheCode = strings.ReplaceAll(cacheCode, memcacheTTLReplacer, strconv.Itoa(requestData.MemcacheTTL))
		f.Write(convertToByteln(cacheCode))
	}

	f.Write(convertToByteln("const ("))
	for _, component := range crudComponent {
		switch component {
		case "CREATE":
			createQuery := constCreateQuery
			createQuery = strings.ReplaceAll(createQuery, domainNameCamelCaseReplacer, stringy.New(requestData.DomainName).CamelCase())
			createQuery = strings.ReplaceAll(createQuery, tableNameReplacer, requestData.TableName)
			attrList := []string{}
			for _, attr := range requestData.Attributes {
				attrList = append(attrList, attr.ColumnName)
			}
			createQuery = strings.ReplaceAll(createQuery, attributesListReplacer, strings.Join(attrList, ","))
			createQuery = strings.ReplaceAll(createQuery, paramNumberListReplacer, generateParamNumberList(1, len(requestData.Attributes)))
			f.Write(convertToByteln(createQuery))
		case "READ":
			readQuery := constReadQuery
			readQuery = strings.ReplaceAll(readQuery, domainNameCamelCaseReplacer, stringy.New(requestData.DomainName).CamelCase())
			readQuery = strings.ReplaceAll(readQuery, tableNameReplacer, requestData.TableName)
			readQuery = strings.ReplaceAll(readQuery, attributeQueryReplacer, generateQueryFromAttributes(requestData.Attributes, requestData.TableName))
			f.Write(convertToByteln(readQuery))
		case "UPDATE":
			updateQuery := constUpdateQuery
			updateQuery = strings.ReplaceAll(updateQuery, domainNameCamelCaseReplacer, stringy.New(requestData.DomainName).CamelCase())
			updateQuery = strings.ReplaceAll(updateQuery, tableNameReplacer, requestData.TableName)
			attributeSetList := []string{}
			for i, attr := range requestData.Attributes {
				// set id in query, and dont append to list
				if attr.ColumnName == "id" {
					updateQuery = strings.ReplaceAll(updateQuery, idSetReplacer, generateSetAttribute(i, attr.ColumnName))
				} else {
					attributeSetList = append(attributeSetList, generateSetAttribute(i, attr.ColumnName))
				}
			}
			updateQuery = strings.ReplaceAll(updateQuery, attributeSetListReplacer, strings.Join(attributeSetList, ",\n"))
			f.Write(convertToByteln(updateQuery))
		case "DELETE":
			deleteQuery := constDeleteQuery
			deleteQuery = strings.ReplaceAll(deleteQuery, domainNameCamelCaseReplacer, stringy.New(requestData.DomainName).CamelCase())
			deleteQuery = strings.ReplaceAll(deleteQuery, tableNameReplacer, requestData.TableName)
			for i, attr := range requestData.Attributes {
				if attr.ColumnName == "id" {
					deleteQuery = strings.ReplaceAll(deleteQuery, idSetReplacer, generateSetAttribute(i, attr.ColumnName))
					break
				}
			}
			f.Write(convertToByteln(deleteQuery))
		}
	}
	f.Write(convertToByteln(")"))

	if requestData.Options.IsUseSingleflight {
		singleflightCode := constSingleflight
		singleflightCode = strings.ReplaceAll(singleflightCode, domainNameCamelCaseReplacer, stringy.New(requestData.DomainName).CamelCase())
		singleflightCode = strings.ReplaceAll(singleflightCode, domainNameLowerCaseReplacer, stringy.New(requestData.DomainName).RemoveSpecialCharacter())
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
		oldPackage := fmt.Sprintf("package mock_%s", stringy.New(requestData.DomainName).RemoveSpecialCharacter())
		newPackage := fmt.Sprintf("package %s", stringy.New(requestData.DomainName).RemoveSpecialCharacter())
		fileContent = []byte(strings.Replace(string(fileContent), oldPackage, newPackage, 1))

		// remove self import
		selfImportPattern := fmt.Sprintf(`%s "%s/result/%s"`, replaceDashWithUnderscore(requestData.DomainName), utils.GeneratorName, requestData.DomainName)
		fileContent = []byte(strings.Replace(string(fileContent), selfImportPattern, "", 1))

		// remove all self import reference
		selfImportReferencePattern := fmt.Sprintf("%s.", replaceDashWithUnderscore(requestData.DomainName))
		fileContent = []byte(strings.ReplaceAll(string(fileContent), selfImportReferencePattern, ""))

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
	importCode = strings.ReplaceAll(importCode, appNameReplacer, requestData.AppName)
	f.Write(convertToBytelnln(importCode))

	// all func
	funcCode := funcDatabaseCode
	funcCode = strings.ReplaceAll(funcCode, domainNameCamelCaseReplacer, stringy.New(requestData.DomainName).CamelCase())
	funcCode = strings.ReplaceAll(funcCode, domainNameLowerCamelCaseReplacer, strcase.ToLowerCamel(requestData.DomainName))
	funcCode = strings.ReplaceAll(funcCode, domainNameLowerCaseReplacer, stringy.New(requestData.DomainName).RemoveSpecialCharacter())
	funcCode = strings.ReplaceAll(funcCode, scanGetAttributeReplacer, generateScanGetAttribute(requestData.Attributes))
	funcCode = strings.ReplaceAll(funcCode, scanInsertAttributeReplacer, generateScanInsertAttribute(requestData.Attributes))
	funcCode = strings.ReplaceAll(funcCode, scanUpdateAttributeReplacer, generateScanUpdateAttribute(requestData.Attributes))
	f.Write(convertToBytelnln(funcCode))
}

func GenerateCache(path string, requestData RequestData) {
	f := createFile(path)
	defer f.Close()

	// package
	writePackageName(f, requestData.DomainName)

	// import
	importCode := importCacheCode
	importCode = strings.ReplaceAll(importCode, appNameReplacer, requestData.AppName)
	f.Write(convertToBytelnln(importCode))

	// all func
	funcCode := funcCacheCode
	funcCode = strings.ReplaceAll(funcCode, domainNameCamelCaseReplacer, stringy.New(requestData.DomainName).CamelCase())
	funcCode = strings.ReplaceAll(funcCode, domainNameLowerCamelCaseReplacer, strcase.ToLowerCamel(requestData.DomainName))
	funcCode = strings.ReplaceAll(funcCode, domainNameLowerCaseReplacer, stringy.New(requestData.DomainName).RemoveSpecialCharacter())
	funcCode = strings.ReplaceAll(funcCode, scanGetAttributeReplacer, generateScanGetAttribute(requestData.Attributes))
	funcCode = strings.ReplaceAll(funcCode, scanInsertAttributeReplacer, generateScanInsertAttribute(requestData.Attributes))
	funcCode = strings.ReplaceAll(funcCode, scanUpdateAttributeReplacer, generateScanUpdateAttribute(requestData.Attributes))
	f.Write(convertToBytelnln(funcCode))
}

func GenerateSingleflight(path string, requestData RequestData) {
	f := createFile(path)
	defer f.Close()

	// package
	writePackageName(f, requestData.DomainName)

	// import
	importCode := importSingleflightCode
	importCode = strings.ReplaceAll(importCode, appNameReplacer, requestData.AppName)
	f.Write(convertToBytelnln(importCode))

	// all func
	funcCode := funcSingleflightCode
	funcCode = strings.ReplaceAll(funcCode, copyAttributeReplacer, generateCopyAttribute(requestData.Attributes))
	funcCode = strings.ReplaceAll(funcCode, domainNameCamelCaseReplacer, stringy.New(requestData.DomainName).CamelCase())
	funcCode = strings.ReplaceAll(funcCode, acronymReplacer, generateAppNameAcronym(requestData.AppName))
	funcCode = strings.ReplaceAll(funcCode, domainNameLowerCamelCaseReplacer, strcase.ToLowerCamel(requestData.DomainName))
	funcCode = strings.ReplaceAll(funcCode, domainNameLowerCaseReplacer, stringy.New(requestData.DomainName).RemoveSpecialCharacter())
	f.Write(convertToBytelnln(funcCode))
}
