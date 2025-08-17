package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"unicode"
	// bencode "github.com/jackpal/bencode-go" // Available if you need it!
)

// Ensures gofmt doesn't remove the "os" encoding/json import (feel free to remove this!)
var _ = json.Marshal

//this func decodes strings and int
func decodeStrInt(bencodedString string) (interface{}, int, error) {
	if unicode.IsDigit(rune(bencodedString[0])) {
		var firstColonIndex int

		for i := 0; i < len(bencodedString); i++ {
			if bencodedString[i] == ':' {
				firstColonIndex = i
				break
			}
		}

		lengthStr := bencodedString[:firstColonIndex]

		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			return "", 0, err
		}

		return bencodedString[firstColonIndex+1 : firstColonIndex+1+length], firstColonIndex+1+length, nil
	} else { // if bencodedString[0] == 'i'
		var end int

		for i := 1; i < len(bencodedString); i ++ {
			if bencodedString[i] == 'e' {
				end = i
				break
			}
		}

		res, err := strconv.Atoi(bencodedString[1:end])
		if err != nil {
			return "", 0, err
		}

		return res, end + 1, nil
	}
}

func main() {
	command := os.Args[1]

	if command == "decode" {
		bencodedValue := os.Args[2]

		if unicode.IsDigit(rune(bencodedValue[0])) || bencodedValue[0] == 'i' {
			decoded, _, err := decodeStrInt(bencodedValue)
			if err != nil {
				fmt.Println(err)
				return
			}
			jsonOutput, _ := json.Marshal(decoded)
			fmt.Println(string(jsonOutput))
		}
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
