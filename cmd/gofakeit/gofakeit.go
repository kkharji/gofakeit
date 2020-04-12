package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/brianvoe/gofakeit/v4"
)

var noFuncRunMsg = "Could not find function to run\nCall list for available functions"

func main() {
	gofakeit.Seed(0)

	args := os.Args[1:]
	argsLen := len(args)

	// Make sure they passed first argument for function call
	if argsLen < 1 {
		fmt.Println(noFuncRunMsg)
		return
	}

	// Get function name
	function := args[0]

	// If function is list output list
	if function == "list" {
		selectedCat := ""
		if argsLen >= 2 {
			selectedCat = args[1]
		}
		listOutput(selectedCat)
		return
	}

	// Lookup fake data method
	info := gofakeit.GetLookupData(function)
	if info == nil {
		fmt.Println(noFuncRunMsg)
		return
	}

	// Set function and params
	var params map[string][]string
	paramsLen := len(info.Params)
	if paramsLen > 0 {
		for i := 0; i < argsLen; i++ {
			if i == 0 {
				continue
			}

			// if params not set make it
			if params == nil {
				params = make(map[string][]string)
			}

			// Map argument to param field
			if paramsLen >= i {
				p := info.Params[i-1]
				if strings.Contains(p.Type, "array") {
					params[p.Field] = strings.Split(args[i], ",")
				} else {
					params[p.Field] = []string{args[i]}
				}
			}
		}
	}

	value, err := info.Call(&params, info)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(fmt.Sprintf("%v", value))
}

func listOutput(selectedCategory string) {
	stringInSlice := func(a string, list []string) bool {
		for _, b := range list {
			if b == a {
				return true
			}
		}
		return false
	}

	// Get list of categories
	categories := []string{}
	for _, l := range gofakeit.MapLookups.Map {
		// If selected category is set only grab of that category
		if selectedCategory != "" && selectedCategory != l.Category {
			continue
		}

		if !stringInSlice(l.Category, categories) {
			categories = append(categories, l.Category)
		}
	}

	// Sort categories
	sort.Strings(categories)

	for i := 0; i < len(categories); i++ {
		fmt.Println(categories[i])

		funcNames := []string{}

		// Get all in category
		for fName, l := range gofakeit.MapLookups.Map {
			if categories[i] == l.Category && !stringInSlice(fName, funcNames) {
				funcNames = append(funcNames, fName)
			}
		}

		// Sort categories
		sort.Strings(funcNames)

		// Output func info
		for _, fName := range funcNames {
			info := gofakeit.GetLookupData(fName)
			fmt.Println("    " + fName + " - " + info.Description)
			for _, p := range info.Params {
				fmt.Println("        Field Name: " + p.Field + " Type: " + p.Type + " Default: " + p.Default + " - " + p.Description)
			}
		}

		fmt.Println()
	}
}
