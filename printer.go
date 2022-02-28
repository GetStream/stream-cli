package cli

import (
	"encoding/json"
	"fmt"
)

func PrintMessage(message string) {
	fmt.Println(" > " + message)
}

func PrintMessageFormat(message string, args ...interface{}) {
	PrintMessage(fmt.Sprintf(message, args...))
}

func PrintHappyMessage(message string) {
	PrintMessage(message + " ✅")
}

func PrintHappyMessageFormatted(message string, args ...interface{}) {
	PrintHappyMessage(fmt.Sprintf(message, args...))
}

func PrintSadMessage(message string) {
	PrintMessage(message + " ❌")
}

func PrintRawJson(data interface{}) {
	res, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(res))
}
