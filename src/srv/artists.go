package srv

import (
	"encoding/json"
	"net/http"

	"github.com/el10savio96/artistInsights/src/lib"
	"go.uber.org/zap"
)

func (s *Server) GetArtistsHandler(w http.ResponseWriter, r *http.Request) {
	artists, err := s.artists.GetArtists(r.Context())
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err)
		return
	}

	if artists == nil {
		artists = []lib.Artist{}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string][]lib.Artist{"artists": artists}); err != nil {
		s.logger.Error("encode artists", zap.Error(err))
	}
}
