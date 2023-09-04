package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func readOutput(path string) (string, error) {
	// Open the file for reading
	filePath := path
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return "", err
	}
	defer file.Close()

	// Read the file content line by line
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return "", err
	}
	return string(content), nil
}

func saveToFile(content []byte, path string) error {
	filepath := path
	err := os.WriteFile(filepath, content, 0644)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return err
	}
	return nil
}

func timeDelta(start time.Time, name string) {
	log.Println("\n*********************")
	log.Printf("%s script called ", name)
	log.Printf("%s takes: %s ", name, time.Since(start))
	time := fmt.Sprintf("%s takes: %s \n", name, time.Since(start)) + "\n\n"
	file, err := os.OpenFile("time.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Cant write time benchmarks to the file. ", err)
	}
	defer file.Close()
	file.Write([]byte(time))
	log.Println("*********************\n")

}
