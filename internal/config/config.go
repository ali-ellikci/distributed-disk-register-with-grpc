package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const toleranceKey = "TOLERANCE"

func LoadConfig() (int, error) {
	file, err := os.Open("tolerance.conf")
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		data := strings.Join(strings.Fields(scanner.Text()), "")
		parts := strings.SplitN(data, "=", 2)
		if len(parts) == 2 && parts[0] == toleranceKey {
			var tolerance int
			_, err := fmt.Sscanf(parts[1], "%d", &tolerance)
			if err != nil {
				return 0, err
			}
			return tolerance, nil
		} else {
			return 0, fmt.Errorf("invalid config format")
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return 0, fmt.Errorf("tolerance.conf is empty or missing TOLERANCE")

}
