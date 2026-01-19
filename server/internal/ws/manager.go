package ws

import (
	"sync"
	"time"

	"github.com/gofiber/websocket/v2"
)

type Manager struct {
	mu              sync.RWMutex
	rooms           map[string]*room
	cacheSize       int
	roomTTL         time.Duration
	cleanupInterval time.Duration
	done            chan struct{}
}

type Config struct {
	CacheSize       int
	RoomTTL         time.Duration
	CleanupInterval time.Duration
}

type room struct {
	mu           sync.Mutex
	clients      map[*Client]struct{}
	cache        []cachedMessage
	lastActivity time.Time
}

type cachedMessage struct {
	payload []byte
	at      time.Time
}

type Client struct {
	conn    *websocket.Conn
	writeMu sync.Mutex
}

func NewManager(cfg Config) *Manager {
	cacheSize := cfg.CacheSize
	if cacheSize <= 0 {
		cacheSize = 3
	}
	roomTTL := cfg.RoomTTL
	if roomTTL <= 0 {
		roomTTL = 30 * time.Second
	}
	cleanupInterval := cfg.CleanupInterval
	if cleanupInterval <= 0 {
		cleanupInterval = 5 * time.Second
	}
	manager := &Manager{
		rooms:           make(map[string]*room),
		cacheSize:       cacheSize,
		roomTTL:         roomTTL,
		cleanupInterval: cleanupInterval,
		done:            make(chan struct{}),
	}
	go manager.cleanupLoop()
	return manager
}

func (m *Manager) Join(roomKey string, conn *websocket.Conn) (*Client, [][]byte) {
	if roomKey == "" || conn == nil {
		return nil, nil
	}

	m.mu.Lock()
	rm := m.rooms[roomKey]
	if rm == nil {
		rm = &room{
			clients:      make(map[*Client]struct{}),
			cache:        []cachedMessage{},
			lastActivity: time.Now(),
		}
		m.rooms[roomKey] = rm
	}
	cl := &Client{conn: conn}
	rm.mu.Lock()
	rm.clients[cl] = struct{}{}
	rm.lastActivity = time.Now()
	cached := make([][]byte, len(rm.cache))
	for i, msg := range rm.cache {
		cached[i] = append([]byte(nil), msg.payload...)
	}
	rm.mu.Unlock()
	m.mu.Unlock()

	return cl, cached
}

func (m *Manager) Leave(roomKey string, cl *Client) {
	if roomKey == "" || cl == nil {
		return
	}
	m.mu.RLock()
	rm := m.rooms[roomKey]
	m.mu.RUnlock()
	if rm == nil {
		return
	}
	rm.mu.Lock()
	delete(rm.clients, cl)
	rm.lastActivity = time.Now()
	rm.mu.Unlock()
}

func (m *Manager) Broadcast(roomKey string, payload []byte) {
	if roomKey == "" || len(payload) == 0 {
		return
	}
	m.mu.RLock()
	rm := m.rooms[roomKey]
	m.mu.RUnlock()
	if rm == nil {
		return
	}

	rm.mu.Lock()
	rm.cache = append(rm.cache, cachedMessage{payload: append([]byte(nil), payload...), at: time.Now()})
	if len(rm.cache) > m.cacheSize {
		rm.cache = rm.cache[len(rm.cache)-m.cacheSize:]
	}
	rm.lastActivity = time.Now()
	clients := make([]*Client, 0, len(rm.clients))
	for cl := range rm.clients {
		clients = append(clients, cl)
	}
	rm.mu.Unlock()

	for _, cl := range clients {
		if err := cl.Write(payload); err != nil {
			m.Leave(roomKey, cl)
		}
	}
}

func (m *Manager) Close() {
	close(m.done)
}

func (c *Client) Write(payload []byte) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	return c.conn.WriteMessage(websocket.TextMessage, payload)
}

func (m *Manager) cleanupLoop() {
	ticker := time.NewTicker(m.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.cleanupRooms()
		case <-m.done:
			return
		}
	}
}

func (m *Manager) cleanupRooms() {
	now := time.Now()
	m.mu.Lock()
	for key, rm := range m.rooms {
		rm.mu.Lock()
		expired := len(rm.clients) == 0 && now.Sub(rm.lastActivity) > m.roomTTL
		rm.mu.Unlock()
		if expired {
			delete(m.rooms, key)
		}
	}
	m.mu.Unlock()
}
