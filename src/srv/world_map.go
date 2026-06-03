package srv

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/el10savio/artistInsights/src/lib"
	"go.uber.org/zap"
)

type worldMapRequest struct {
	ArtistID string    `validate:"required"`
	Start    time.Time `validate:"required"`
	End      time.Time `validate:"required"`
}

func (s *Server) GetWorldMapHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	var req worldMapRequest
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

	data, err := s.artists.GetWorldMap(r.Context(), req.ArtistID, req.Start, req.End)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err)
		return
	}

	if data == nil {
		data = []lib.CountryListens{}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string][]lib.CountryListens{"world_map": data}); err != nil {
		s.logger.Error("encode world map", zap.Error(err))
	}
}
