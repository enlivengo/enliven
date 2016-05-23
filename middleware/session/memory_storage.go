package session

import (
	"net/http"
	"strconv"
	"time"

	"github.com/enlivengo/enliven"
	"github.com/enlivengo/enliven/config"
	"github.com/jmcvetta/randutil"
)

var sessions map[string]*StoredSession

// StoredSession represents a user's session in memory
type StoredSession struct {
	mTime int32
	data  map[string]string
}

// newMemorySession Produces a memory-based session instance
func newMemorySession(sessID string) *memorySession {
	fSess := &memorySession{
		sessionID: sessID,
	}

	if _, ok := sessions[sessID]; !ok {
		sessions[sessID] = &StoredSession{
			mTime: int32(time.Now().Unix()),
			data:  make(map[string]string),
		}
	}

	return fSess
}

// memorySession implements the enliven.ISession interface
type memorySession struct {
	sessionID string
}

// Set sets a session variable
func (ms *memorySession) Set(key string, value string) error {
	storedSession := sessions[ms.sessionID]
	storedSession.data[key] = value
	return nil
}

// Get returns a session variable or empty string
func (ms *memorySession) Get(key string) string {
	storedSession := sessions[ms.sessionID]
	if val, ok := storedSession.data[key]; ok {
		return val
	}
	return ""
}

// Delete removes a session variable
func (ms *memorySession) Delete(key string) error {
	storedSession := sessions[ms.sessionID]
	if _, ok := storedSession.data[key]; ok {
		delete(storedSession.data, key)
	}
	return nil
}

// Destroy deletes this session from redis
func (ms *memorySession) Destroy() error {
	delete(sessions, ms.sessionID)
	return nil
}

// SessionID returns the current session id
func (ms *memorySession) SessionID() string {
	return ms.sessionID
}

// NewMemoryStorageMiddleware generates an instance of MemoryStorageMiddleware
func NewMemoryStorageMiddleware() *MemoryStorageMiddleware {
	return &MemoryStorageMiddleware{}
}

// MemoryStorageMiddleware manages sessions, using memory as the session storage mechanism
type MemoryStorageMiddleware struct {
	lastPurge int32
	purgeTTL  int32
	ttl       int32
	purging   bool
}

// Initialize sets up the session middleware
func (msm *MemoryStorageMiddleware) Initialize(ev *enliven.Enliven) {
	sessions = make(map[string]*StoredSession)

	conf := config.Config{
		"session_memory_ttl":      "86400",
		"session_memory_purgettl": "1800",
	}

	conf = config.UpdateConfig(config.MergeConfig(conf, config.GetConfig()))

	purgeGap, _ := strconv.Atoi(conf["session_memory_purgettl"])
	sessionTTL, _ := strconv.Atoi(conf["session_memory_ttl"])

	msm.lastPurge = int32(time.Now().Unix())
	msm.purgeTTL = int32(purgeGap)
	msm.ttl = int32(sessionTTL)
	msm.purging = false
}

// GetName returns the middleware's name
func (msm *MemoryStorageMiddleware) GetName() string {
	return "session"
}

func (msm *MemoryStorageMiddleware) ServeHTTP(ctx *enliven.Context, next enliven.NextHandlerFunc) {
	sessionID, err := ctx.Request.Cookie("enlivenSession")
	var sID string
	// If there was no cookie, we create a session id
	if err == nil {
		sID = sessionID.Value
	} else {
		sID, _ = randutil.AlphaString(32)
		cookie := http.Cookie{Name: "enlivenSession", Value: sID}
		http.SetCookie(ctx.Response, &cookie)
	}

	ctx.Session = newMemorySession(sID)

	msm.purgeSessions()

	next(ctx)
}

func (msm *MemoryStorageMiddleware) purgeSessions() {
	// Returning if we're already in the process of purging
	if msm.purging {
		return
	}

	current := int32(time.Now().Unix())

	if current > msm.lastPurge+msm.purgeTTL {
		msm.purging = true

		// Finding all the sessions whose last modified time is more than our TTL ago
		for key, session := range sessions {
			if session.mTime < current-msm.ttl {
				delete(sessions, key)
			}
		}

		msm.lastPurge = current
		msm.purging = false
	}
}
