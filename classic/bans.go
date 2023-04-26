package classic

import (
	"log"
	"os"
	"strings"
)

const BANS_FILENAME = "bans.txt"

func LoadBans() ([]string, error) {
	blankBans := []string{}
	if _, err := os.Stat(BANS_FILENAME); err != nil {
		if os.IsNotExist(err) {
			log.Println("Creating new bans.txt...")
			f, err := os.Create(BANS_FILENAME)
			if err != nil {
				return blankBans, err
			}
			f.Close()
			return blankBans, nil
		} else {
			return blankBans, err
		}
	}

	contents, err := os.ReadFile(BANS_FILENAME)
	if err != nil {
		return blankBans, err
	}

	lines := strings.Split(string(contents), "\n")
	for i, val := range lines {
		lines[i] = strings.TrimSpace(val)
	}

	return lines, nil
}

func SaveBans(bans []string) error {
	var sb strings.Builder
	for _, ban := range bans {
		sb.WriteString(ban)
		sb.WriteByte('\n')
	}

	err := os.WriteFile(BANS_FILENAME, []byte(sb.String()), 0644)
	return err
}
