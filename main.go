package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"slices"
)

var options = make(map[int]string)
var endpoints = []string{"https://google.com"}
var optionsHandlers = make(map[int]func())

const ADD_ENDPOINT = 1
const LIST_ENDPOINTS = 2
const TEST_ALL = 3
const REMOVE_ENDPOINT = 4

type DelayConfiguration struct {
	Unit string  `json:"unit"`
	Time float32 `json:"time"`
}

type RetryConfiguration struct {
	Tries int                `json:"retries"`
	Delay DelayConfiguration `json:"delay"`
}

type Endpoint struct {
	Url   string             `json:"url"`
	Retry RetryConfiguration `json:"retry"`
}

type ConfigurationFile struct {
	Endpoints []Endpoint `json:"endpoints"`
}

func generateConfigurationFile(path string) {
	saveConfigurationFile(path, generateBaseConfigurationFile())
}

func generateBaseConfigurationFile() *ConfigurationFile {

	return &ConfigurationFile{
		Endpoints: []Endpoint{
			{
				Url: "https://github.com/natanmatosz",
				Retry: RetryConfiguration{
					Tries: 3,
					Delay: DelayConfiguration{
						Unit: "s",
						Time: 5,
					},
				},
			},
		},
	}
}

func saveConfigurationFile(path string, file *ConfigurationFile) {
	configurationFileContent, err := json.MarshalIndent(file, "", "\t")

	if err != nil {
		fmt.Println("Error while generating the configuration file.")
		return
	}

	err = os.WriteFile("./test.json", configurationFileContent, 0644)

	if err != nil {
		fmt.Println("Error while saving configuration file.")
	}
}

func main() {
	// generateConfigurationFile("./test.json")
	setupApplication()
	showOptions()

	selectedOption := getOptionFromCli()

	handleSelectedOption(selectedOption)
}

func setupApplication() {
	registerOptions()
	registerOptionsHandlers()
}

func registerOptions() {
	options[TEST_ALL] = "Add enpoint"
	options[LIST_ENDPOINTS] = "List endpoints"
	options[TEST_ALL] = "Test All"
	options[REMOVE_ENDPOINT] = "Remove Endpoint"
}

func registerOptionsHandlers() {
	optionsHandlers[TEST_ALL] = handleTestAll
	optionsHandlers[LIST_ENDPOINTS] = handleListAll
}

func showOptions() {
	for idx, option := range options {
		fmt.Printf("[%d] - %s\n", idx, option)
	}
}

func getOptionFromCli() int {
	var option int

	validOptionSelected := false

	for !validOptionSelected {
		fmt.Print("\033[0;32mSelect an option:\033[0m ")
		_, err := fmt.Scan(&option)

		validOptionSelected = err == nil

		if !validOptionSelected {
			fmt.Println("\033[0;31mSelect a valid option!\033[0m")
		}
	}

	return option
}

func handleSelectedOption(selectedOption int) {
	optionsHandlers[selectedOption]()
}

func handleTestAll() {
	if len(endpoints) < 1 {
		fmt.Println("\033[0;31mNo endpoints to test!\033[0m")
		return
	}

	for _, endpoint := range endpoints {
		fmt.Printf("\033[0;33mTesting endpoint \"%s\"\033[0m\n", endpoint)

		if !isEndpointWorking(endpoint) {
			fmt.Printf("\033[0;31mEndpoint \"%s\" is not working!\033[0m\n", endpoint)
			return
		}

		fmt.Printf("\033[0;32mEndpoint \"%s\" is working fine!\033[0m\n", endpoint)
	}
}

func isEndpointWorking(endpoint string) bool {
	response, err := http.Get(endpoint)

	const MIN_SUCCESS_STATUS_CODE = 200
	const MAX_SUCCESS_STATUS_CODE = 300

	if err != nil {
		return false
	}

	return response.StatusCode >= MIN_SUCCESS_STATUS_CODE && response.StatusCode < MAX_SUCCESS_STATUS_CODE
}

func handleListAll() {
	for _, endpoint := range endpoints {
		fmt.Println(endpoint)
	}
}

func handleRemoveEndpoint() {
	selectedIdx := getOptionFromCli()

	slices.Delete(endpoints, selectedIdx, 0)
}
