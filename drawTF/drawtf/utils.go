package drawtf

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

// readOutput reads the content of a file specified by the path and returns it as a string.
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

// saveToFile writes the given content to a file specified by the path.
func saveToFile(content []byte, path string) error {
	filepath := path
	err := os.WriteFile(filepath, content, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return err
	}
	return nil
}

// timeDelta logs the time taken by a script and appends it to a "time.log" file.
func timeDelta(start time.Time, name string) {
	log.Println("\n*********************")
	log.Printf("%s script called ", name)
	log.Printf("%s takes: %s ", name, time.Since(start))
	time := fmt.Sprintf("%s takes: %s \n", name, time.Since(start)) + "\n\n"
	file, err := os.OpenFile("time.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Can't write time benchmarks to the file. ", err)
	}
	defer file.Close()
	file.Write([]byte(time))
	log.Println("*********************\n")
}

// launchGradio launches the Gradio application using the specified Python interpreter and Gradio launcher.
func launchGradio() error {
	pythonCmd := exec.Command(PythonInterpreter, GradioLauncher)
	pythonCmd.Stdout = os.Stdout
	pythonCmd.Stderr = os.Stderr
	err := pythonCmd.Start()
	fmt.Println("Error is:", err)
	if err != nil {
		fmt.Println("Error running Gradio Launcher:", err)
		return err
	}
	return nil
}
