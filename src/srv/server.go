package srv

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Server struct {
	artists  ArtistService
	logger   *zap.Logger
	mux      *http.ServeMux
	validate *validator.Validate
}

func New(artists ArtistService, logger *zap.Logger) *Server {
	s := &Server{
		artists:  artists,
		logger:   logger,
		mux:      http.NewServeMux(),
		validate: validator.New(),
	}
	s.mux.HandleFunc("GET /ping", s.PingHandler)
	s.mux.HandleFunc("GET /artists", s.GetArtistsHandler)
	s.mux.HandleFunc("GET /summary", s.GetSummaryHandler)
	s.mux.HandleFunc("GET /series", s.GetSeriesHandler)
	s.mux.HandleFunc("GET /histogram", s.GetHistogramHandler)
	s.mux.HandleFunc("GET /top-tracks", s.GetTopTracksHandler)
	s.mux.HandleFunc("GET /heatmap", s.GetHeatmapHandler)
	s.mux.HandleFunc("GET /world-map", s.GetWorldMapHandler)
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) Start(addr string) error {
	s.logger.Info("server starting", zap.String("addr", addr))
	return http.ListenAndServe(addr, s.mux)
}

func validateTimeRange(start, end time.Time) error {
	if end.Before(start) {
		return errors.New("end must not be before start")
	}
	return nil
}

func (s *Server) writeError(w http.ResponseWriter, status int, err error) {
	s.logger.Error("request error", zap.Int("status", status), zap.Error(err))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": http.StatusText(status)})
}
