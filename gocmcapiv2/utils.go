package gocmcapiv2

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func genUUID() string {
	b := make([]byte, 16)

	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	// Set version (4)
	b[6] = (b[6] & 0x0f) | 0x40

	// Set variant (RFC4122)
	b[8] = (b[8] & 0x3f) | 0x80

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4],
		b[4:6],
		b[6:8],
		b[8:10],
		b[10:16],
	)
}

func convert2JsonString(object interface{}) string {
	jsonString, err := json.Marshal(object)
	if err != nil {
		fmt.Printf("err %v", err)
		return ""
	}
	return string(jsonString)
}

// Logo log object
func Logo(pre string, object interface{}) {
	str, ok := os.LookupEnv("DEBUG_CMCCLOUD_TERRAFORM") // os.Getenv("DEBUG_CMCCLOUD_TERRAFORM")
	logfile, okLogFile := os.LookupEnv("DEBUG_CMCCLOUD_LOGFILE")

	if ok && str == "DEBUG" {
		if !okLogFile {
			logfile = "log.txt"
		}
		f, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		// defer f.Close()
		defer func() { _ = f.Close() }()

		logger := log.New(f, "", log.LstdFlags)
		// logger.Println(object)
		jsonString, err := json.Marshal(object)
		if err != nil {
			logger.Println("Error:", err)
			return
		}

		// Print JSON string
		logger.Println(pre + string(jsonString))
	}
}

func Logs(message string) {
	str, ok := os.LookupEnv("DEBUG_CMCCLOUD_TERRAFORM")
	logfile, okLogFile := os.LookupEnv("DEBUG_CMCCLOUD_LOGFILE")

	// Logall("str DEBUG_CMCCLOUD_TERRAFORM = " + str)
	if ok && str == "DEBUG" {
		if !okLogFile {
			logfile = "log.txt"
		}
		// Logall("str DEBUG_CMCCLOUD_LOGFILE = " + logfile)
		f, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		// defer f.Close()
		defer func() { _ = f.Close() }()

		logger := log.New(f, "", log.LstdFlags)
		logger.Println(message)
	}
}

func Logall(message string) {
	logfile := "log.txt"
	f, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer func() { _ = f.Close() }()

	logger := log.New(f, "", log.LstdFlags)
	logger.Println(message)
}
