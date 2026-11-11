package realtime

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

// Two tabs signed in as the same admin both stay connected and both get
// events; before, each new connection closed the other.
func TestSameAdminCanHaveSeveralConnections(t *testing.T) {
	hub := NewHub()
	upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		hub.RegisterConnection(conn, 7)
	}))
	defer server.Close()

	url := "ws" + strings.TrimPrefix(server.URL, "http")
	dial := func() *websocket.Conn {
		t.Helper()
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { conn.Close() })
		return conn
	}
	first, second := dial(), dial()

	// Wait for both to be registered.
	deadline := time.Now().Add(2 * time.Second)
	for {
		hub.mu.RLock()
		n := len(hub.clients)
		hub.mu.RUnlock()
		if n == 2 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("expected 2 connections for one admin, have %d", n)
		}
		time.Sleep(10 * time.Millisecond)
	}

	hub.Broadcast("new_order", map[string]string{"orderNumber": "FG-1"})
	for name, conn := range map[string]*websocket.Conn{"first": first, "second": second} {
		_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		_, msg, err := conn.ReadMessage()
		if err != nil || !strings.Contains(string(msg), "FG-1") {
			t.Fatalf("%s tab should get the event, got %q %v", name, msg, err)
		}
	}
}
