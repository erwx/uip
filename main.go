package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

// ActionStep represents one row from the CSV file
type ActionStep struct {
	District string
	School   string
	Step     string
	Strategy string
}

// NewStep creates a new ActionStep with the given values
func NewStep(district, school, step, strategy string) *ActionStep {
	return &ActionStep{
		District: district,
		School:   school,
		Step:     step,
		Strategy: strategy,
	}
}

// parseCSV reads a CSV file and returns a slice of ActionSteps
func parseCSV(filename string) ([]*ActionStep, error) {
	// Open file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer file.Close()

	// Create CSV reader
	reader := csv.NewReader(file)

	// Read headers
	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read headers: %w", err)
	}

	// Create header map
	headerMap := make(map[string]int)
	for idx, header := range headers {
		headerMap[header] = idx
	}

	var steps []*ActionStep

	// Read data rows
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read record: %w", err)
		}

		// Create ActionStep from record
		step := NewStep(
			record[headerMap["UIP: District Name"]],
			record[headerMap["UIP: School Name"]],
			record[headerMap["Improvement Action Step"]],
			record[headerMap["Major Improvement Strategy"]],
		)
		steps = append(steps, step)
	}

	return steps, nil
}

func main() {
	steps, err := parseCSV("ActionSteps.csv")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Loaded %d steps:\n", len(steps))
	for i := 0; i < 4 && i < len(steps); i++ {
		fmt.Printf("Step %d: %+v\n", i+1, *steps[i])
	}
}
