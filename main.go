package main

import (
	"domain-crud-generator/generator"
	"domain-crud-generator/utils"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func main() {
	var requestData generator.RequestData

	// read file sample.json
	config, err := os.ReadFile(utils.FilenameConfig)
	if err != nil {
		utils.Print(err)
	}

	err = json.Unmarshal(config, &requestData)
	if err != nil {
		utils.Print(err)
	}

	os.Mkdir(utils.ResultFolderPath, 0770)

	// make domain name to lower case
	requestData.DomainName = strings.ToLower(requestData.DomainName)

	// remove folder domain if exist
	if _, err := os.Stat(utils.GenerateResultPath(requestData.DomainName)); !os.IsNotExist(err) {
		os.RemoveAll(utils.GenerateResultPath(requestData.DomainName))
	}

	// create new folder
	utils.Print(fmt.Sprintf("Creating Folder %s", requestData.DomainName))
	err = os.Mkdir(utils.GenerateResultPath(requestData.DomainName), 0770)
	if err != nil {
		utils.Print(err)
	}

	// TODO : make the process below to be async (exclude generate mock)

	// create file type
	utils.Print("Creating File type.go")
	pathType := utils.GenerateResultPath("type.go", requestData.DomainName)
	generator.GenerateType(pathType, requestData)

	// create file init
	utils.Print(fmt.Sprintf("Creating File %s.go", requestData.DomainName))
	pathInit := utils.GenerateResultPath(fmt.Sprintf("%s.go", requestData.DomainName), requestData.DomainName)
	generator.GenerateInit(pathInit, requestData)

	defer func() {
		// create file mock
		utils.Print("Creating File mock.go")
		pathMock := utils.GenerateResultPath("mock.go", requestData.DomainName)
		generator.GenerateMock(pathMock, requestData)
	}()

	// create file const
	utils.Print("Creating File const.go")
	pathConst := utils.GenerateResultPath("const.go", requestData.DomainName)
	generator.GenerateConst(pathConst, requestData)

	// create file database
	utils.Print("Creating File database.go")
	pathDatabase := utils.GenerateResultPath("database.go", requestData.DomainName)
	generator.GenerateDatabase(pathDatabase, requestData)

	// create file cache
	if requestData.Options.IsUseRedisAndMemcache {
		utils.Print("Creating File cache.go")
		pathCache := utils.GenerateResultPath("cache.go", requestData.DomainName)
		generator.GenerateCache(pathCache, requestData)
	}

	// create file singleflight
	if requestData.Options.IsUseSingleflight {
		utils.Print("Creating File singleflight.go")
		pathSingleflight := utils.GenerateResultPath("singleflight.go", requestData.DomainName)
		generator.GenerateSingleflight(pathSingleflight, requestData)
	}
}
