package srv

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/el10savio96/artistInsights/src/lib"
	"go.uber.org/zap"
)

type seriesRequest struct {
	ArtistID    string    `validate:"required"`
	Start       time.Time `validate:"required"`
	End         time.Time `validate:"required"`
	Granularity string    `validate:"required,oneof=hour day week month"`
}

func (s *Server) GetSeriesHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	var req seriesRequest
	req.ArtistID = q.Get("artist_id")
	req.Granularity = q.Get("granularity")

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

	points, err := s.artists.GetSeries(r.Context(), req.ArtistID, req.Start, req.End, req.Granularity)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err)
		return
	}

	if points == nil {
		points = []lib.SeriesPoint{}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string][]lib.SeriesPoint{"series": points}); err != nil {
		s.logger.Error("encode series", zap.Error(err))
	}
}
