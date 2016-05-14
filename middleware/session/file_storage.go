package session

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/hickeroar/enliven"
	"github.com/jmcvetta/randutil"
)

// newFileSession Produces a file-based session instance
func newFileSession(sessID string, dir string) *fileSession {
	dir += (sessID + ".sess")

	fSess := &fileSession{
		sessionID: sessID,
		path:      dir,
	}

	return fSess
}

// fileSession implements the enliven.ISession interface
type fileSession struct {
	sessionID string
	path      string
}

func (fs *fileSession) getSessionData() map[string]string {
	if _, err := os.Stat(fs.path); os.IsNotExist(err) {
		return make(map[string]string)
	}

	var data map[string]string

	rawData, _ := ioutil.ReadFile(fs.path)
	json.Unmarshal(rawData, &data)

	return data
}

func (fs *fileSession) writeSessionData(data map[string]string) error {
	jsonData, _ := json.Marshal(data)

	return ioutil.WriteFile(fs.path, jsonData, 0755)
}

// Set sets a session variable
func (fs *fileSession) Set(key string, value string) error {
	sessionData := fs.getSessionData()
	sessionData[key] = value
	return fs.writeSessionData(sessionData)
}

// Get returns a session variable or empty string
func (fs *fileSession) Get(key string) string {
	sessionData := fs.getSessionData()
	if val, ok := sessionData[key]; ok {
		return val
	}
	return ""
}

// Delete removes a session variable
func (fs *fileSession) Delete(key string) error {
	sessionData := fs.getSessionData()
	if _, ok := sessionData[key]; ok {
		delete(sessionData, key)
	}
	return nil
}

// Destroy deletes this session from redis
func (fs *fileSession) Destroy() error {
	return os.Remove(fs.path)
}

// SessionID returns the current session id
func (fs *fileSession) SessionID() string {
	return fs.sessionID
}

// Path returns the path to the session storage file.
func (fs *fileSession) Path() string {
	return fs.path
}

// NewFileStorageMiddleware generates an instance of FileStorageMiddleware
func NewFileStorageMiddleware(suppliedConfig enliven.Config) *FileStorageMiddleware {
	var config = enliven.Config{
		"session.file.path":     "/tmp/",
		"session.file.ttl":      "86400",
		"session.file.purgettl": "1800",
	}

	config = enliven.MergeConfig(config, suppliedConfig)

	dir := config["session.file.path"]

	if dir[len(dir)-1:] != "/" {
		dir += "/"
	}

	purgeGap, _ := strconv.Atoi(config["session.file.purgettl"])
	sessionTTL, _ := strconv.Atoi(config["session.file.ttl"])

	return &FileStorageMiddleware{
		path:      dir,
		lastPurge: int32(time.Now().Unix()),
		purgeTTL:  int32(purgeGap),
		ttl:       int32(sessionTTL),
		purging:   false,
	}
}

// FileStorageMiddleware manages sessions, using the filesystem as the session storage mechanism
type FileStorageMiddleware struct {
	path      string
	lastPurge int32
	purgeTTL  int32
	ttl       int32
	purging   bool
}

func (fsm *FileStorageMiddleware) ServeHTTP(ctx *enliven.Context, next enliven.NextHandlerFunc) {
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

	session := newFileSession(sID, fsm.path)
	ctx.Session = session

	fsm.purgeSessions()

	next(ctx)
}

func (fsm *FileStorageMiddleware) purgeSessions() {
	// Returning if we're already in the process of purging
	if fsm.purging {
		return
	}

	current := int32(time.Now().Unix())

	if current > fsm.lastPurge+fsm.purgeTTL {
		fsm.purging = true

		// Holds all the file names we want to delete
		var toDelete []string

		// Finding all the files whose last modified time is more than our TTL ago
		files, _ := ioutil.ReadDir(fsm.path)
		for _, f := range files {
			fName := f.Name()

			// Only looking at session files
			if fName[len(fName)-5:] != ".sess" {
				continue
			}

			fmTime := int32(f.ModTime().Unix())
			if fmTime < current-fsm.ttl {
				toDelete = append(toDelete, fName)
			}
		}

		// Deleting each file which is beyond the ttl
		for _, d := range toDelete {
			os.Remove(fsm.path + d)
		}

		fsm.lastPurge = current
		fsm.purging = false
	}
}