package main

import(
	"io/ioutil"
	"encoding/json"
	"gleipnir/error"
	"gleipnir/service"
)

type Configuration struct {
	Services service.Services `json: "services"`
}


func main() {
	loadConfig()
}

func loadConfig() {
	data, err := ioutil.ReadFile("config.json")
	error.Check(err)
	
	var config Configuration
	
	err = json.Unmarshal([]byte(data), &config)
	error.Check(err)
}