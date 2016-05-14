package session

import (
	"encoding/base64"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/hickeroar/enliven"
	"github.com/jmcvetta/randutil"
	"gopkg.in/redis.v3"
)

// newRedisSession Produces a redis session instance
func newRedisSession(sessID string, rClient *redis.Client, existing bool) *redisSession {
	rSess := &redisSession{
		redisClient: rClient,
		sessionID:   sessID,
	}

	rSess.bump(existing)

	return rSess
}

// redisSession implements the enliven.Session interface
type redisSession struct {
	redisClient *redis.Client
	sessionID   string
}

// Resets the current session's expiration date to 24 hours in the future
// and sets the init time if this is a new session
func (rs *redisSession) bump(existing bool) {
	if !existing {
		rs.Set("init", time.Now().String())
	}
	// Setting the duration of the
	rs.redisClient.Expire(rs.sessionID, time.Duration(24)*time.Hour)
}

// Set sets a session variable
func (rs *redisSession) Set(key string, value string) error {
	_, err := rs.redisClient.HSet(rs.sessionID, key, value).Result()
	return err
}

// Get returns a session variable or empty string
func (rs *redisSession) Get(key string) string {
	value, err := rs.redisClient.HGet(rs.sessionID, key).Result()
	if err != nil {
		return ""
	}
	return value
}

// Delete removes a session variable
func (rs *redisSession) Delete(key string) error {
	_, err := rs.redisClient.HDel(rs.sessionID, key).Result()
	return err
}

// Destroy deletes this session from redis
func (rs *redisSession) Destroy() error {
	_, err := rs.redisClient.Del(rs.sessionID).Result()
	return err
}

// SessionID returns the current session id
func (rs *redisSession) SessionID() string {
	return rs.sessionID
}

// NewRedisStorageMiddleware generates an instance of RedisStorageMiddleware
func NewRedisStorageMiddleware(suppliedConfig enliven.Config) *RedisStorageMiddleware {
	var config = enliven.Config{
		"session.redis.address":  "127.0.0.1:6379",
		"session.redis.password": "",
		"session.redis.database": "0",
	}

	config = enliven.MergeConfig(config, suppliedConfig)

	database, _ := strconv.Atoi(config["session.redis.database"])

	return &RedisStorageMiddleware{
		redisClient: redis.NewClient(&redis.Options{
			Addr:     config["session.redis.address"],
			Password: config["session.redis.password"],
			DB:       int64(database),
		}),
	}
}

// RedisStorageMiddleware manages sessions, using redis as the session storage mechanism
type RedisStorageMiddleware struct {
	redisClient *redis.Client
}

// generateSessionID produces a unique session id that we can store on the user's end as a cookie.
func (rsm *RedisStorageMiddleware) generateSessionID() string {
	rb := make([]byte, 32)
	rand.Read(rb)
	return base64.URLEncoding.EncodeToString(rb)
}

func (rsm *RedisStorageMiddleware) ServeHTTP(ctx *enliven.Context, next enliven.NextHandlerFunc) {
	sessionID, err := ctx.Request.Cookie("enlivenSession")
	var existing bool
	var sID string
	// If there was no cookie, we create a session id, a session in redis, and a cookie to hold the ID.
	if err == nil {
		existing = true
		sID = sessionID.Value
	} else {
		existing = false
		sID, _ = randutil.AlphaString(32)
		cookie := http.Cookie{Name: "enlivenSession", Value: sID}
		http.SetCookie(ctx.Response, &cookie)
	}

	ctx.Session = newRedisSession(sID, rsm.redisClient, existing)

	next(ctx)
}
