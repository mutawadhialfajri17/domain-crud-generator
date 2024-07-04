package utils

import (
	"strconv"
	"strings"

	"github.com/gobeam/stringy"
	"github.com/iancoleman/strcase"
)

var (
	StringToLower    = strings.ToLower
	StringSplit      = strings.Split
	StringReplace    = strings.Replace
	StringReplaceAll = strings.ReplaceAll
	StringJoin       = strings.Join

	// hello world => Hello World
	StringToCamelCase = func(s string) string {
		return stringy.New(s).CamelCase()
	}

	// hello *#$! world => hello world
	StringRemoveSpecialCharacter = func(s string) string {
		return stringy.New(s).RemoveSpecialCharacter()
	}

	// Hello World => hw
	StringToAcronym = func(s string) string {
		return stringy.New(s).Acronym().ToLower()
	}

	// HelloWorld => helloWorld
	StringToLowerCamel = strcase.ToLowerCamel

	IntToString = strconv.Itoa
)
