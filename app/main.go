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

		for i := 1; i < len(bencodedString); i++ {
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

//this function decodes bencoded list
func decodeList(bencodedString string) (interface{}, int, error) {
	i := 1
	res := make([]interface{}, 0)
	for i < len(bencodedString) {
		if unicode.IsDigit(rune(bencodedString[i])) || bencodedString[i] == 'i' {
			curr, next, err := decodeStrInt(bencodedString[i:])
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return nil, 0, err
			}
			res = append(res, curr)
			i += next
		} else if bencodedString[i] == 'l' {
			//handle nested lists recursively
			curr, next, err := decodeList(bencodedString[i:])
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return nil, 0, err
			}
			res = append(res, curr)
			i += next
		} else {
			break
		}
	}
	return res, i + 1, nil
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
		} else if bencodedValue[0] == 'l' {
			decoded, _, err := decodeList(bencodedValue)
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
