package classic

import (
	"log"
	"os"
	"strings"
)

const OPS_FILENAME = "ops.txt"

func LoadOPs() ([]string, error) {
	blankOps := []string{}
	if _, err := os.Stat(OPS_FILENAME); err != nil {
		if os.IsNotExist(err) {
			log.Println("Creating new ops.txt file")
			f, err := os.Create(OPS_FILENAME)
			if err != nil {
				return blankOps, err
			}
			f.Close()
			return blankOps, nil
		} else {
			return blankOps, err
		}
	}

	contents, err := os.ReadFile(OPS_FILENAME)
	if err != nil {
		return blankOps, err
	}

	lines := strings.Split(string(contents), "\n")
	for i, val := range lines {
		lines[i] = strings.TrimSpace(val)
	}

	return lines, nil
}

func SaveOPs(ops []string) error {
	var sb strings.Builder
	for _, op := range ops {
		sb.WriteString(op)
		sb.WriteByte('\n')
	}

	err := os.WriteFile(OPS_FILENAME, []byte(sb.String()), 0644)
	return err
}
