package utils

import "strings"

func GenerateResultPath(filename string, folderPath ...string) string {
	if len(folderPath) == 0 {
		return ResultFolderPath + "/" + filename
	}
	return ResultFolderPath + "/" + strings.Join(folderPath, "/") + "/" + filename
}
