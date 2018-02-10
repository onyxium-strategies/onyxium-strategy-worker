package main

import (
//"encoding/json"
"github.com/tidwall/gjson"
//"reflect"
"fmt"
"io/ioutil"
"os"
)

func parseJson(jsonFile string) []interface{} {
	file, e := ioutil.ReadFile(jsonFile)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	myJson := string(file)

	data, ok := gjson.Parse(myJson).Value().([]interface{})
	if !ok {
		fmt.Println("Error")
	}
	return data
}

func main(){
	data := parseJson("./tree-example.json")
	fmt.Println(data)
}
