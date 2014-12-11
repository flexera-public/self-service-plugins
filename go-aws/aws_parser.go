package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {

	var aws map[string]Service

	if err := json.NewDecoder(os.Stdin).Decode(&aws); err != nil {
		log.Fatalf("Error: %s", err.Error())
	}

	fmt.Println("package aws")
	fmt.Print("var Aws = ")
	def := fmt.Sprintf("%#v\n", aws)
	def = strings.Replace(def, "struct {", "struct {\n", -1)
	def = strings.Replace(def, "}{", "}{\n", -1)
	fmt.Println(def)
	log.Println("Done")
}
