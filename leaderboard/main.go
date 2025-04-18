package main

import (
	"fmt"
	"log"
)

func main() {
	InitDB()
	InitRedis()

	go PeriodicSyncWorker()

	_ = UpdatePlayerScore("player1", 10)
	_ = UpdatePlayerScore("player2", 5)
	_ = UpdatePlayerScore("player1", 20)

	topPlayers, err := GetTopPlayers(10)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("🏆 Leaderboard:")
	for idx, p := range topPlayers {
		fmt.Printf("%d. %s - %d points\n", idx+1, p.Member.(string), int(p.Score))
	}

	select {}
}
