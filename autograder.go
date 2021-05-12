package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

var testCases = []TestCase{
	{40.0, "Basic simulator functionality", "Working", "asdf", "d"},
	{10.0, "Makefile", "", "", ""},
	{30.0, "Execution Statistics", "Prints out correct output", "asdf", "asdf"},
}

type TestCase struct {
	MaxScore    float64 `json:"max_score"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Expected    string
	Arguments   string
}

type Result struct {
	Score          float64       `json:"score"`
	Execution_time string        `json:"execution_time"`
	Output         string        `json:"output"`
	Visibility     string        `json:"visibility"`
	Leaderboard    []Leaderboard `json:"leaderboard"`
}

type Leaderboard struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Order string `json:"order"`
}

type Test struct {
	Score      float64 `json:"score"`
	MaxScore   float64 `json:"max_score"`
	Output     string  `json:"output"`
	Name       string  `json:"name"`
	Visibility string  `json:"visibility"`
}

func RunCmd(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	cmd.Dir, _ = os.Getwd()
	out, err := cmd.Output()
	return string(out), err
}

func TimeCmd(command string, args ...string) float64 {
	start := time.Now()
	RunCmd("./myISS", "sample.assembly")
	elapsed := time.Since(start)
	return float64(elapsed.Seconds())
}

func createTestJSON(score float64, maxScore float64, output string, name string, visibility string) Test {
	return Test{
		Score:      score,
		MaxScore:   maxScore,
		Output:     output,
		Name:       name,
		Visibility: visibility,
	}
}

func failTestCases(testCases []TestCase, failMsg string) []Test {
	results := []Test{}
	for _, t := range testCases {
		results = append(results, createTestJSON(0.0, t.MaxScore, t.Description+failMsg, t.Name, "after_published"))
	}
	return results
}

func Grade() string {
	score := 0.0
	out, err := RunCmd("make")
	if err != nil {
		result := failTestCases(testCases, "Did not compile:\n"+out)
		data, _ := json.MarshalIndent(result, "", "  ")
		return string(data)
	}
	cmd := "./myISS"
	args := "sample.assembly"
	out, err = RunCmd(cmd, args)
	if err != nil {
		result := failTestCases(testCases, "Did not execute:\n"+out)
		data, _ := json.MarshalIndent(result, "", "  ")
		return string(data)
	}

	// Get average times
	times := [10]float64{}
	sum := 0.0
	for i := 0; i < 10; i++ {
		times[i] = TimeCmd(cmd, args)
		sum += times[i]
	}
	avg_time := sum / float64(len(times))
	execution_time := fmt.Sprintf("%.6f", avg_time)

	execStruct := Leaderboard{
		Name:  "Time",
		Value: execution_time,
		Order: "asc",
	}
	arr := []Leaderboard{}
	arr = append(arr, execStruct)
	result := Result{
		Score:          score,
		Execution_time: execution_time,
		Output:         string(out),
		Visibility:     "visible",
		Leaderboard:    arr,
	}
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Println(err)
	}
	RunCmd("make", "clean")
	return string(data)
}

func main() {
	fmt.Println(Grade())
}
