package main

import(
	"io/ioutil"
	"gleipnir/error"
	"encoding/json"
	
)

type Configuration struct {

	Services []string `json: "services"`

}

func main() {

	data, err := ioutil.ReadFile("config.json")
	error.Check(err)
	
	var config Configuration
	
	err = json.Unmarshal([]byte(data), &config)
	error.Check(err)
}