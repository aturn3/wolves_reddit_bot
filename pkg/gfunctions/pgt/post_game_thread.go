package pgt

import (
	"log"
	"time"

	"github.com/drewthor/wolves_reddit_bot/apis/nba"
	"github.com/drewthor/wolves_reddit_bot/apis/reddit"
	"github.com/drewthor/wolves_reddit_bot/pkg/gcloud"
)

func CreatePostGameThread(teamTriCode nba.TriCode) {
	currentTimeUTC := time.Now().UTC()
	// Issues occur when using eastern time for "today's games" as games on the west coast can still be going on
	// when the eastern time rolls over into the next day
	westCoastLocation, locationError := time.LoadLocation("America/Los_Angeles")
	if locationError != nil {
		log.Fatal(locationError)
	}
	currentTimeWestern := currentTimeUTC.In(westCoastLocation)
	currentDateWestern := currentTimeWestern.Format(nba.TimeDayFormat)
	dailyAPIPaths := nba.GetDailyAPIPaths()
	teams := nba.GetTeams(dailyAPIPaths.Teams)
	teamID := teams[teamTriCode].ID
	scheduledGames := nba.GetScheduledGames(dailyAPIPaths.TeamSchedule, teamID)
	todaysGame, gameToday := scheduledGames[currentDateWestern]
	if gameToday {
		log.Println("game today")
		todaysGameScoreboard := nba.GetGameScoreboard(dailyAPIPaths.Scoreboard, currentDateWestern, todaysGame.GameID)
		boxscore := nba.GetBoxscore(dailyAPIPaths.Boxscore, currentDateWestern, todaysGame.GameID)
		if boxscore.GameEnded() {
			log.Println("game ended")
			if todaysGameScoreboard.EndTimeUTC != "" {
				gameEndTimeUTC, err := time.Parse(nba.UTCFormat, todaysGameScoreboard.EndTimeUTC)
				if err != nil {
					log.Fatal(err)
				}
				log.Println(gameEndTimeUTC)
			}
			datastore := new(gcloud.Datastore)
			gameEvent, exists := datastore.GetTeamGameEvent(todaysGame.GameID, teamID)
			if exists && gameEvent.PostGameThread {
				log.Println("already found post")
				return
			}
			log.Println("making post")
			redditClient := reddit.Client{}
			redditClient.Authorize()
			log.Println("authorized")
			subreddit := "SeattleSockeye"
			title := boxscore.GetRedditPostGameThreadTitle(teamTriCode, teams)
			content := boxscore.GetRedditPostGameThreadBodyString(nba.GetPlayers(dailyAPIPaths.Players))
			redditClient.SubmitNewPost(subreddit, title, content)

			gameEvent.CreatedTime = time.Now()
			gameEvent.GameID = todaysGame.GameID
			gameEvent.TeamID = teamID
			gameEvent.PostGameThread = true
			datastore.SaveTeamGameEvent(gameEvent)
		}
	}
}