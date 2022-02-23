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

type Game struct {
	PGN string `json:"pgn"`
}

type MonthlyGames struct {
	Games []Game `json:"games"`
}

func usergames(client *req.Client, url string) []Game {
	monthly := &MonthlyGames{}
	resp, err := client.R().SetHeader("Accept", "application/vnd.github.v3+json").
		SetResult(monthly).
		Get(url)
	if err != nil {
		log.Fatal(err)
	}
	games := []Game{}
	if resp.IsSuccess() {
		for _, game := range monthly.Games {
			games = append(games, game)
		}
	}
	return games
}

func archives(client *req.Client, username string) []Game {
	archives := &UserArchives{}
	resp, err := client.R().SetHeader("Accept", "application/vnd.github.v3+json").
		SetPathParam("username", username).
		SetResult(archives).
		Get("https://api.chess.com/pub/player/{username}/games/archives")
	if err != nil {
		log.Fatal(err)
	}
	// list of PGNs
	games := []Game{}
	if resp.IsSuccess() {
		for _, url := range archives.Archives {
			games = append(games, usergames(client, url)...)
		}
	}
	return games
}

func main() {
	client := req.C().
		SetTimeout(5 * time.Second)

	fmt.Println("Chess is dumb")

	for _, game := range archives(client, "lwnexgen") {
		Blunders(game.PGN)
		break
	}
}
