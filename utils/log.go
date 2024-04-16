package utils

import (
	"fmt"
	"log"
)

func Print(in interface{}) {
	if e, ok := in.(error); ok {
		log.Fatalln("ERROR : ", e)
	} else {
		log.Println("INFO : ", in)
	}
}

func init() {
	fmt.Println("======= Domain CRUD Generator =======")
}
