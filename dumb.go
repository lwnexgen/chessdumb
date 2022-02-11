package main

import (
	"fmt"
	"github.com/imroc/req/v3"
	"log"
	"time"
)

type UserArchives struct {
	Archives []string `json:"archives"`
}

type MonthlyGames struct {
	Games []map[string]string `json:"games"`
}

func usergames(client *req.Client, url string) []string {
	monthly := &MonthlyGames{}
	resp, err := client.R().SetHeader("Accept", "application/vnd.github.v3+json").
		SetResult(monthly).
		Get(url)
	if err != nil {
		log.Fatal(err)
	}
	if resp.IsSuccess() {
		for _, games := range monthly.Games {
			for index, game := range games {
				fmt.Println(
					fmt.Sprintf("%d: %+v", index, game))
			}
		}
	}

	return []string{}
}

func archives(client *req.Client, username string) {
	archives := &UserArchives{}
	resp, err := client.R().SetHeader("Accept", "application/vnd.github.v3+json").
		SetPathParam("username", username).
		SetResult(archives).
		Get("https://api.chess.com/pub/player/{username}/games/archives")
	if err != nil {
		log.Fatal(err)
	}
	// list of PGNs
	games := []string{}
	if resp.IsSuccess() {
		for _, url := range archives.Archives {
			games = append(games, usergames(client, url)...)
		}
	}
}

func main() {
	client := req.C().
		SetTimeout(5 * time.Second)

	fmt.Println("Chess is dumb")

	archives(client, "lwnexgen")
}
