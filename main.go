package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

func parseFile(path string) map[string]interface{} {
	file, _ := os.Open(path)
	defer file.Close()

	decoder := json.NewDecoder(file)
	data := make(map[string]interface{})
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Println("Error:", err)
	}
	return data
}

func validate() bool {
	dir, _ := os.Getwd()
	files, _ := os.ReadDir(dir)
	var jsons []string

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if len(file.Name()) > 5 && file.Name()[len(file.Name())-5:] == ".json" {
			jsons = append(jsons, file.Name())
		}
	}

	for _, file := range jsons {
		data := parseFile(file)
		if data["end"] == "1" && len(data["variants"].([]interface{})) > 0 {
			return false
		}
		for _, e := range data["variants"].([]interface{}) {
			if _, err := os.Stat(e.(map[string]interface{})["next_file"].(string)); os.IsNotExist(err) {
				return false
			}
		}
	}
	return true
}

func main() {
	if validate() {
		curr := "start.json"
		for parseFile(curr)["end"].(string) != "1" {
			fmt.Println(parseFile(curr)["text"])
			for i, variant := range parseFile(curr)["variants"].([]interface{}) {
				fmt.Println(i+1, ") ", variant.(map[string]interface{})["text"])
			}
			fmt.Print("введите ответ числом: ")
			var k string
			fmt.Scanln(&k)
			for {
				if _, err := strconv.Atoi(k); err == nil {
					break
				}
				fmt.Print("введите корректный ответ числом: ")
				fmt.Scanln(&k)
			}

			idx, _ := strconv.Atoi(k)
			curr = parseFile(curr)["variants"].([]interface{})[idx-1].(map[string]interface{})["next_file"].(string)
		}

		fmt.Println(parseFile(curr)["text"])
	}
}
