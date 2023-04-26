package classic

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const SETTINGS_FILENAME = "settings.txt"

type Settings struct {
	IP          string
	Port        uint16
	Name        string
	MOTD        string
	Online      bool
	Public      bool
	Password    string
	WorldX      int16
	WorldY      int16
	WorldZ      int16
	PlayerCount uint8
}

func CreateDefaultSettings() *Settings {
	return &Settings{
		IP:          "",
		Port:        25565,
		Name:        "Classic Server",
		MOTD:        "An implementation of Classic Minecraft written in Go",
		Online:      false,
		Public:      false,
		Password:    "",
		WorldX:      256,
		WorldY:      256,
		WorldZ:      256,
		PlayerCount: 128,
	}
}

func LoadSettings() (*Settings, error) {
	settings := CreateDefaultSettings()
	if _, err := os.Stat(SETTINGS_FILENAME); err != nil {
		if os.IsNotExist(err) {
			log.Println("Creating new settings.txt...")
			err := SaveSettings(settings)
			if err != nil {
				return nil, err
			}
			return settings, nil
		} else {
			return nil, err
		}
	}

	contents, err := os.ReadFile(SETTINGS_FILENAME)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(contents), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		segments := strings.SplitN(line, "=", 2)
		if len(segments) != 2 {
			log.Println("Unable to interpret setting \"", line, "\"")
			continue
		}

		key := segments[0]
		value := segments[1]

		switch key {

		case "ip":
			settings.IP = value

		case "port":
			if parsed, err := strconv.ParseUint(value, 10, 16); err != nil {
				log.Println(err)
				log.Printf("Unable to interpret setting \"port\" value \"%s\"\n", value)
			} else {
				settings.Port = uint16(parsed)
			}

		case "name":
			settings.Name = value

		case "motd":
			settings.MOTD = value

		case "online":
			if parsed, err := strconv.ParseBool(value); err != nil {
				log.Println(err)
				log.Printf("Unable to interpret setting \"online\" value \"%s\"\n", value)
			} else {
				settings.Online = parsed
			}

		case "public":
			if parsed, err := strconv.ParseBool(value); err != nil {
				log.Println(err)
				log.Printf("Unable to interpret setting \"public\" value \"%s\"\n", value)
			} else {
				settings.Public = parsed
			}

		case "password":
			settings.Password = value

		case "worldX":
			if parsed, err := strconv.ParseInt(value, 10, 16); err != nil {
				log.Println(err)
				log.Printf("Unable to interpret setting \"worldX\" value \"%s\"\n", value)
			} else {
				settings.WorldX = int16(parsed)
			}

		case "worldY":
			if parsed, err := strconv.ParseInt(value, 10, 16); err != nil {
				log.Println(err)
				log.Printf("Unable to interpret setting \"worldY\" value \"%s\"\n", value)
			} else {
				settings.WorldY = int16(parsed)
			}

		case "worldZ":
			if parsed, err := strconv.ParseInt(value, 10, 16); err != nil {
				log.Println(err)
				log.Printf("Unable to interpret setting \"worldZ\" value \"%s\"\n", value)
			} else {
				settings.WorldZ = int16(parsed)
			}

		case "playerCount":
			if parsed, err := strconv.ParseUint(value, 10, 8); err != nil {
				log.Println(err)
				log.Printf("Unable to interpret setting \"playerCount\" value \"%s\"\n", value)
			} else {
				settings.PlayerCount = uint8(parsed)
			}

		}
	}

	return settings, nil
}

func SaveSettings(settings *Settings) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("ip=%s\n", settings.IP))
	sb.WriteString(fmt.Sprintf("port=%d\n", settings.Port))
	sb.WriteString(fmt.Sprintf("name=%s\n", settings.Name))
	sb.WriteString(fmt.Sprintf("motd=%s\n", settings.MOTD))
	sb.WriteString(fmt.Sprintf("online=%t\n", settings.Online))
	sb.WriteString(fmt.Sprintf("public=%t\n", settings.Public))
	sb.WriteString(fmt.Sprintf("password=%s\n", settings.Password))
	sb.WriteString(fmt.Sprintf("worldX=%d\n", settings.WorldX))
	sb.WriteString(fmt.Sprintf("worldY=%d\n", settings.WorldY))
	sb.WriteString(fmt.Sprintf("worldZ=%d\n", settings.WorldZ))
	sb.WriteString(fmt.Sprintf("playerCount=%d\n", settings.PlayerCount))

	err := os.WriteFile(SETTINGS_FILENAME, []byte(sb.String()), 0644)
	return err
}
