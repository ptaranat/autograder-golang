package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

type Result struct {
	Score          float64     `json:"score"`
	Execution_time float64     `json:"execution_time"`
	Output         string      `json:"output"`
	Visibility     string      `json:"visibility"`
	Leaderboard    Leaderboard `json:"leaderboard"`
}

type Leaderboard struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
	Order string  `json:"order"`
}

func RunExe(executable string, arguments ...string) string {
	cmd := exec.Command(executable, arguments...)
	cmd.Dir, _ = os.Getwd()
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	return string(out)
}

func Grade() string {
	RunExe("make")
	start := time.Now()
	out := RunExe("./myISS", "sample.assembly")
	elapsed := time.Since(start)
	execution_time := float64(elapsed.Seconds())
	result := Result{
		Score:          0.0,
		Execution_time: execution_time,
		Output:         string(out),
		Visibility:     "visible",
		Leaderboard: Leaderboard{
			Name:  "Time",
			Value: execution_time,
			Order: "asc",
		},
	}
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Println(err)
	}
	RunExe("make", "clean")
	return string(data)
}

func main() {
	fmt.Println(Grade())
}
