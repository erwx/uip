package main

import (
	"encoding/csv"
  "encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

// ActionStep represents one row from the CSV file
type ActionStep struct {
	District    string
	Step        string
  Description string
  Start       string
  Target      string
	Strategy    string
  Resources   string
}

type DistrictBatch struct {
  DistrictName string       `json:"district_name"` 
  ActionSteps  []ActionInfo `json:"action_steps"`
}

type ActionInfo struct {
  Step        string `json:"step"`
  Description string `json:"description"`
  Start       string `json:"start"`
  Target      string `json:"target"`
  Strategy    string `json:"strategy"`
  Resources   string `json:"resources"`
}

func createBatch(districtName string, steps []*ActionStep) *DistrictBatch {
  batch := DistrictBatch{
    DistrictName: districtName,
    ActionSteps: []ActionInfo{},
  }

  for _, step := range steps {
    info := ActionInfo{
      Step:        step.Step,
      Description: step.Description,
      Start:       step.Start,
      Target:      step.Target,
      Strategy:    step.Strategy,
      Resources:   step.Resources,
    }
    batch.ActionSteps = append(batch.ActionSteps, info)
  }
  return &batch
}

func parseRecord(record []string, headerMap map[string]int) *ActionStep {
  return &ActionStep {
    District:    record[headerMap["UIP: District Name"]],
    Step:        record[headerMap["Improvement Action Step"]],
    Description: record[headerMap["Description of Action Step"]],
    Start:       record[headerMap["Start Date"]],
    Target:      record[headerMap["Target Date"]],
    Strategy:    record[headerMap["Major Improvement Strategy"]],
    Resources:   record[headerMap["Resources"]],
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

    step := parseRecord(record, headerMap)
		steps = append(steps, step)
	}

	return steps, nil
}

func main() {
	steps, err := parseCSV("ActionSteps.csv")

	if err != nil {
		log.Fatal(err)
	}

  // Map district names to ActionSteps
  districtMap := make(map[string][]*ActionStep)
  for _, step := range steps {
    district := step.District
    districtMap[district] = append(districtMap[district], step)
  }

  var batches []*DistrictBatch

  for districtName, actionSteps := range districtMap {
    b := createBatch(districtName, actionSteps)
    batches = append(batches, b)
  }

  for i := range batches {
    j, err := json.Marshal(batches[i])
    if err != nil {
      fmt.Printf("JSON error: %v", err)
    }
    fmt.Println(string(j))
  }
}
