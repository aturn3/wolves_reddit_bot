package main

import (
	"fmt"
	"github.com/drewthor/wolves_reddit_bot/apis/nba"
	"log"
	"time"
)

func main() {
	currentTimeUTC := time.Now().UTC()
	fmt.Println(currentTimeUTC)
	currentTimeEastern := time.Now().UTC()
	currentDateEastern := currentTimeEastern.Format(nba.TimeDayFormat)
	dailyAPIPaths := nba.GetDailyAPIPaths()
	teams := nba.GetTeams(dailyAPIPaths.Teams)
	wolvesID := teams["MIN"].ID
	scheduledGames := nba.GetScheduledGames(dailyAPIPaths.TeamSchedule, wolvesID)
	todaysGame, gameToday := scheduledGames[currentDateEastern]
	if gameToday {
		todaysGameScoreboard := nba.GetGameScoreboard(dailyAPIPaths.Scoreboard, currentDateEastern, todaysGame.GameID)
		if todaysGameScoreboard.EndTimeUTC != "" {
			gameEndTime, err := time.Parse(nba.UTCFormat, todaysGameScoreboard.EndTimeUTC)
			if err != nil {
				log.Fatal(err)
			}
			timeSinceGameEnded := currentTimeUTC.Sub(gameEndTime)
			if timeSinceGameEnded.Minutes() < 1 {
				// make post game thread
			}
		}
		fmt.Println(todaysGameScoreboard)
	}
}
