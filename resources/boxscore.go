package resources

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/drewthor/wolves_reddit_bot/util"

	"github.com/drewthor/wolves_reddit_bot/services"
)

type BoxscoreResource struct {
	BoxscoreService *services.BoxscoreService
}

func (br BoxscoreResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", br.Get)

	return r
}

func (br BoxscoreResource) Get(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	gameID := chi.URLParam(r, "gameID")
	gameDate := r.FormValue("game_date")

	if gameID == "" {
		util.WriteJSON(http.StatusBadRequest, "invalid request: missing game_id", w)
		return
	}

	if gameDate == "" {
		util.WriteJSON(http.StatusBadRequest, "invalid request: missing game_date", w)
		return
	}

	util.WriteJSON(http.StatusOK, br.BoxscoreService.Get(gameID, gameDate), w)
}
