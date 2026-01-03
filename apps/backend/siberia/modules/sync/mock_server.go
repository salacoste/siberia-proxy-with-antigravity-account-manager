package sync

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"

	"github.com/go-chi/chi/v5"
)

// MockServer mimics the behavior of the Cloud Sync API
type MockServer struct {
	Store map[string]SyncPayload // profileID -> payload
	Mu    sync.RWMutex
}

func NewMockServer() *MockServer {
	return &MockServer{
		Store: make(map[string]SyncPayload),
	}
}

func (s *MockServer) Router() http.Handler {
	r := chi.NewRouter()

	r.Post("/sync/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var newPayload SyncPayload
		if err := json.NewDecoder(r.Body).Decode(&newPayload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		s.Mu.Lock()
		defer s.Mu.Unlock()

		existing, exists := s.Store[id]
		if exists && existing.Timestamp > newPayload.Timestamp {
			// Conflict: Server has newer data
			w.WriteHeader(http.StatusConflict)
			return
		}

		// Accept update (Last Write Wins or First Write)
		s.Store[id] = newPayload
		w.WriteHeader(http.StatusOK)
	})

	r.Get("/sync/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		s.Mu.RLock()
		defer s.Mu.RUnlock()

		payload, exists := s.Store[id]
		if !exists {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(payload)
	})

	return r
}

func (s *MockServer) Start() *httptest.Server {
	return httptest.NewServer(s.Router())
}
