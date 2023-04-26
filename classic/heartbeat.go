package classic

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

func NewHeartbeatTicker() *time.Ticker {
	return time.NewTicker(time.Second * 45)
}

func Heartbeat(port uint16, maxPlayers uint8, name string, public bool, salt string, players int) {
	qs := fmt.Sprintf(
		"port=%d&max=%d&name=%s&public=%t&version=7&salt=%s&users=%d",
		port,
		maxPlayers,
		url.QueryEscape(name),
		public,
		salt,
		players,
	)
	url := fmt.Sprintf("https://www.classicube.net/server/heartbeat?%s", qs)
	_, err := http.Get(url)
	if err != nil {
		log.Println(err)
		log.Println("Error requesting heartbeat")
	}
}
