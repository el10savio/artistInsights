package srv

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/el10savio/artistInsights/src/lib"
	"go.uber.org/zap"
)

type topTracksRequest struct {
	ArtistID string    `validate:"required"`
	Start    time.Time `validate:"required"`
	End      time.Time `validate:"required"`
}

func (s *Server) GetTopTracksHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	var req topTracksRequest
	req.ArtistID = q.Get("artist_id")

	if startStr := q.Get("start"); startStr != "" {
		t, err := time.Parse(time.RFC3339, startStr)
		if err != nil {
			s.writeError(w, http.StatusBadRequest, errors.New("invalid start: must be RFC3339"))
			return
		}
		req.Start = t
	}

	if endStr := q.Get("end"); endStr != "" {
		t, err := time.Parse(time.RFC3339, endStr)
		if err != nil {
			s.writeError(w, http.StatusBadRequest, errors.New("invalid end: must be RFC3339"))
			return
		}
		req.End = t
	}

	if err := s.validate.Struct(req); err != nil {
		s.writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := validateTimeRange(req.Start, req.End); err != nil {
		s.writeError(w, http.StatusBadRequest, err)
		return
	}

	tracks, err := s.artists.GetTopTracks(r.Context(), req.ArtistID, req.Start, req.End)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err)
		return
	}

	if tracks == nil {
		tracks = []lib.TrackPlay{}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string][]lib.TrackPlay{"tracks": tracks}); err != nil {
		s.logger.Error("encode top tracks", zap.Error(err))
	}
}
