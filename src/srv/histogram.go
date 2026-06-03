package srv

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/el10savio/artistInsights/src/lib"
	"go.uber.org/zap"
)

type histogramRequest struct {
	ArtistID string    `validate:"required"`
	Start    time.Time `validate:"required"`
	End      time.Time `validate:"required"`
	By       string    `validate:"required,oneof=listener track"`
}

func (s *Server) GetHistogramHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	var req histogramRequest
	req.ArtistID = q.Get("artist_id")
	req.By = q.Get("by")

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

	bins, err := s.artists.GetHistogram(r.Context(), req.ArtistID, req.Start, req.End, req.By)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err)
		return
	}

	if bins == nil {
		bins = []lib.HistogramBin{}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string][]lib.HistogramBin{"histogram": bins}); err != nil {
		s.logger.Error("encode histogram", zap.Error(err))
	}
}
