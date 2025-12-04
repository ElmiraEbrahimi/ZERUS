package api

import (
    "encoding/json"
    "io"
    "net/http"
    "strings"

    "github.com/example/bridge-service/internal/bridge"
)

// Server ties HTTP handlers to the bridge service. It defines endpoints
// for sending messages and querying their status.
type Server struct {
    bridge *bridge.Service
}

// NewServer constructs a Server with the provided bridge. Call Start
// to begin listening for requests.
func NewServer(bridge *bridge.Service) *Server {
    return &Server{bridge: bridge}
}

// Start begins listening on the given address (e.g. ":8080"). It
// blocks until the underlying http.ListenAndServe returns.
func (s *Server) Start(addr string) error {
    return http.ListenAndServe(addr, s.routes())
}

// routes registers handler functions on a new ServeMux.
func (s *Server) routes() http.Handler {
    mux := http.NewServeMux()
    mux.HandleFunc("/bridge/l1-to-l2", s.handleL1ToL2)
    mux.HandleFunc("/bridge/l2-to-l1", s.handleL2ToL1)
    mux.HandleFunc("/bridge/status/", s.handleStatus)
    return mux
}

// handleL1ToL2 accepts POST requests with a JSON body defining the
// L1ToL2Request. It returns the generated message ID in the response.
func (s *Server) handleL1ToL2(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }
    defer r.Body.Close()
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    var req bridge.L1ToL2Request
    if err := json.Unmarshal(body, &req); err != nil {
        http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
        return
    }
    id, err := s.bridge.SendL1ToL2Message(r.Context(), req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"id": id})
}

// handleL2ToL1 accepts POST requests for bridging from L2 to L1.
func (s *Server) handleL2ToL1(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }
    defer r.Body.Close()
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    var req bridge.L2ToL1Request
    if err := json.Unmarshal(body, &req); err != nil {
        http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
        return
    }
    id, err := s.bridge.SendL2ToL1Message(r.Context(), req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"id": id})
}

// handleStatus returns the status of a message given its ID, which is
// extracted from the path ("/bridge/status/<id>").
func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
    // Expected path: /bridge/status/{id}
    id := strings.TrimPrefix(r.URL.Path, "/bridge/status/")
    if id == "" {
        http.Error(w, "missing id", http.StatusBadRequest)
        return
    }
    status, err := s.bridge.TrackStatus(r.Context(), id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"id": id, "status": string(status)})
}