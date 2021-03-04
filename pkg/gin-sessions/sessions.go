package ginsessions

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gin-example/pkg/logging"
	"gin-example/pkg/sessions"
)

const (
	defaultKey = "ginSessions"
)

// Wraps thinly gorilla-session methods.
// Session stores the values and optional configuration for a session.
type GinSessionInterface interface {
	// Get returns the session value associated to the given key.
	Get(key interface{}) interface{}
	// Set sets the session value associated to the given key.
	Set(key interface{}, val interface{})
	// Delete removes the session value associated to the given key.
	Delete(key interface{})
	// Clear deletes all values in the session.
	Clear()
	// AddFlash adds a flash message to the session.
	// A single variadic argument is accepted, and it is optional: it defines the flash key.
	// If not defined "_flash" is used by default.
	AddFlash(value interface{}, vars ...string)
	// Flashes returns a slice of flash messages from the session.
	// A single variadic argument is accepted, and it is optional: it defines the flash key.
	// If not defined "_flash" is used by default.
	Flashes(vars ...string) []interface{}
	// Options sets configuration for a session.
	Options(Options)
	// Save saves all sessions used during the current request.
	Save() error
}

type ginSession struct {
	name    string
	request *http.Request

	store   GinStoreInterface
	session *sessions.Session

	needWritten bool
	writer      http.ResponseWriter
}

func Sessions(name string, store GinStoreInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		s := &ginSession{name, c.Request, store, nil, false, c.Writer}
		c.Set(defaultKey, s)
		defer sessions.Clear(c.Request)
		c.Next()
	}
}

func SessionsMany(names []string, store GinStoreInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		s := make(map[string]GinSessionInterface, len(names))
		for _, name := range names {
			s[name] = &ginSession{name, c.Request, store, nil, false, c.Writer}
		}
		c.Set(defaultKey, s)
		defer sessions.Clear(c.Request)
		c.Next()
	}
}

func (s *ginSession) Session() *sessions.Session {
	if s.session == nil {
		var err error
		s.session, err = s.store.Get(s.request, s.name)
		if err != nil {
			logging.Logger.Warn("get sessions.Session failed")
		}
	}
	return s.session
}

func (s *ginSession) Get(key interface{}) interface{} {
	return s.Session().Values[key]
}

func (s *ginSession) Set(key interface{}, val interface{}) {
	s.Session().Values[key] = val
	s.needWritten = true
}

func (s *ginSession) Delete(key interface{}) {
	delete(s.Session().Values, key)
	s.needWritten = true
}

func (s *ginSession) Clear() {
	for key := range s.Session().Values {
		s.Delete(key)
	}
}

func (s *ginSession) AddFlash(value interface{}, vars ...string) {
	s.Session().AddFlash(value, vars...)
	s.needWritten = true
}

func (s *ginSession) Flashes(vars ...string) []interface{} {
	s.needWritten = true
	return s.Session().Flashes(vars...)
}

func (s *ginSession) Options(options Options) {
	s.Session().Options = options.ToOptions()
}

func (s *ginSession) NeedWritten() bool {
	return s.needWritten
}

func (s *ginSession) Save() error {
	if s.NeedWritten() {
		if err := s.Session().Save(s.request, s.writer); err != nil {
			return err
		}
		s.needWritten = false
	}
	return nil
}

// shortcut to get session
func GetSession(c *gin.Context) GinSessionInterface {
	return c.MustGet(defaultKey).(GinSessionInterface)
}

// shortcut to get session with given name
func GetSessionMany(c *gin.Context, name string) GinSessionInterface {
	return c.MustGet(defaultKey).(map[string]GinSessionInterface)[name]
}
