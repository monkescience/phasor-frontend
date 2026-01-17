package integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"
)

// mockBackendServer creates a mock backend server for testing.
// It returns the instance info endpoint that the frontend expects.
type mockBackendServer struct {
	server    *httptest.Server
	version   string
	hostname  string
	startTime time.Time
}

func newMockBackend(version string) *mockBackendServer {
	m := &mockBackendServer{
		version:   version,
		hostname:  "test-host",
		startTime: time.Now(),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/instance/info", m.instanceInfoHandler)
	mux.HandleFunc("/health/ready", m.healthReadyHandler)
	mux.HandleFunc("/health/live", m.healthLiveHandler)

	m.server = httptest.NewServer(mux)

	return m
}

func (m *mockBackendServer) URL() string {
	return m.server.URL
}

func (m *mockBackendServer) Close() {
	m.server.Close()
}

//nolint:errchkjson // Test helper, error handling not critical.
func (m *mockBackendServer) instanceInfoHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp := struct {
		Version   string `json:"version"`
		Hostname  string `json:"hostname"`
		Uptime    string `json:"uptime"`
		GoVersion string `json:"go_version"`
		Timestamp string `json:"timestamp"`
	}{
		Version:   m.version,
		Hostname:  m.hostname,
		Uptime:    time.Since(m.startTime).String(),
		GoVersion: "go1.25.5",
		Timestamp: time.Now().Format(time.RFC3339Nano),
	}

	_ = json.NewEncoder(w).Encode(resp)
}

//nolint:errchkjson // Test helper, error handling not critical.
func (m *mockBackendServer) healthReadyHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp := map[string]string{
		"status": "ok",
	}

	_ = json.NewEncoder(w).Encode(resp)
}

//nolint:errchkjson // Test helper, error handling not critical.
func (m *mockBackendServer) healthLiveHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp := map[string]string{
		"status": "ok",
	}

	_ = json.NewEncoder(w).Encode(resp)
}
