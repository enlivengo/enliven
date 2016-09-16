package session

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/enlivengo/enliven"
	"github.com/enlivengo/enliven/config"
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
func NewFileStorageMiddleware() *FileStorageMiddleware {
	return &FileStorageMiddleware{}
}

// FileStorageMiddleware manages sessions, using the filesystem as the session storage mechanism
type FileStorageMiddleware struct {
	path      string
	lastPurge int32
	purgeTTL  int32
	ttl       int32
	purging   bool
}

// Initialize sets up the session middleware
func (fsm *FileStorageMiddleware) Initialize(ev *enliven.Enliven) {
	conf := config.Config{
		"session_file_path":     "/tmp/",
		"session_file_ttl":      "86400",
		"session_file_purgettl": "1800",
	}

	conf = config.UpdateConfig(config.MergeConfig(conf, config.GetConfig()))

	dir := conf["session_file_path"]

	if dir[len(dir)-1:] != "/" {
		dir += "/"
	}

	purgeGap, _ := strconv.Atoi(conf["session_file_purgettl"])
	sessionTTL, _ := strconv.Atoi(conf["session_file_ttl"])

	fsm.path = dir
	fsm.lastPurge = int32(time.Now().Unix())
	fsm.purgeTTL = int32(purgeGap)
	fsm.ttl = int32(sessionTTL)
	fsm.purging = false
}

// GetName returns the middleware's name
func (fsm *FileStorageMiddleware) GetName() string {
	return "session"
}

func (fsm *FileStorageMiddleware) ServeHTTP(ctx *enliven.Context, next enliven.NextHandlerFunc) {
	sessionID, err := ctx.Request.Cookie("enlivenSession")
	var sID string
	// If there was no cookie, we create a session id
	if err == nil {
		sID = sessionID.Value
	} else {
		sID, _ = randutil.AlphaString(32)
		cookie := http.Cookie{Name: "enlivenSession", Value: sID, Path: "/"}
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
