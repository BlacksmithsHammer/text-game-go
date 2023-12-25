package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
)

func parseFile(path string) map[string]interface{} {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	var data map[string]interface{}
	err = json.Unmarshal(file, &data)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return data
}

func validate() bool {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		fmt.Println(err)
		return false
	}
	jsons := []string{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if file.Name()[len(file.Name())-5:] == ".json" {
			jsons = append(jsons, file.Name())
		}
	}
	for _, file := range jsons {
		data := parseFile(file)
		if data["end"] == "1" && len(data["variants"].([]interface{})) > 0 {
			return false
		}
		for _, e := range data["variants"].([]interface{}) {
			if !contains(jsons, e.(map[string]interface{})["next_file"].(string)) {
				return false
			}
		}
	}
	return true
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func main() {
	if validate() {
		curr := "start.json"
		for parseFile(curr)["end"].(string) != "1" {
			fmt.Println(parseFile(curr)["text"].(string))
			for i, variant := range parseFile(curr)["variants"].([]interface{}) {
				fmt.Printf("%d) %s\n", i+1, variant.(map[string]interface{})["text"].(string))
			}
			fmt.Print("введите ответ числом: ")
			var k string
			fmt.Scanln(&k)
			for !isNumeric(k) || toInt(k) < 1 || toInt(k) > len(parseFile(curr)["variants"].([]interface{})) {
				fmt.Print("введите корректный ответ числом: ")
				fmt.Scanln(&k)
			}
			curr = parseFile(curr)["variants"].([]interface{})[toInt(k)-1].(map[string]interface{})["next_file"].(string)
		}
		fmt.Println(parseFile(curr)["text"].(string))
	}
}

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func toInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

