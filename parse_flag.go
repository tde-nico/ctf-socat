package main

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var commandMap map[string]func(string) (string, error)
var re *regexp.Regexp
var flag_format string

func genRandomBytes(amount string) (string, error) {
	size, err := strconv.Atoi(amount)
	if err != nil {
		return "", err
	}

	bytes := make([]byte, size)
	_, err = rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func parseRegex(match string) string {
	command := strings.Fields(match[1 : len(match)-1])
	fn, found := commandMap[command[0]]
	if found {
		result, err := fn(command[1])
		if err != nil {
			log.Fatalf("%v\n", err)
		}
		return result
	}
	return ""
}

func getUniqueFlag() string {
	unique_flag := re.ReplaceAllStringFunc(flag_format, parseRegex)
	debugLog("Flag: %v\n", unique_flag)
	return unique_flag
}

func init_regex(f string) {
	flag_format = f
	re = regexp.MustCompile(`\[(.*?)\]`)
	commandMap = map[string]func(string) (string, error){
		"random": genRandomBytes,
	}
}
