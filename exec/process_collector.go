package exec

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type procMetric struct {
	Cpu  float64
	Mem  float64
	Args string
}

func convertLocalizedString(s string) string {
	if strings.ContainsAny(s, ",") {
		return strings.Replace(s, ",", ".", 1)
	} else {
		return s
	}
}

func ProcessCollect(pid int) (*procMetric, error) {
	output, err := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "%cpu=12345,%mem=12345,args").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute 'ps' command: %v", err)
	}

	linesOfProcStrings := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(linesOfProcStrings) < 2 {
		return nil, fmt.Errorf("failed to not find pid %d process ", pid)
	}

	line := linesOfProcStrings[1]

	cpu, err := strconv.ParseFloat(convertLocalizedString(strings.TrimSpace(line[0:5])), 64)
	if err != nil {
		log.Printf("failed to convert cpu to float: %v. split: %v", err, line)
	}

	mem, err := strconv.ParseFloat(convertLocalizedString(strings.TrimSpace(line[6:11])), 64)
	if err != nil {
		log.Printf("failed to convert mem to float: %v. split: %v", err, line)
	}

	return &procMetric{
		Args: line[12:],
		Cpu:  cpu,
		Mem:  mem,
	}, nil

}
