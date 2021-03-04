// Copyright 2012 Brian "bojo" Jones. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package sessions

import (
	"bytes"
	"context"
	"encoding/base32"
	"encoding/gob"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"

	"gin-example/pkg/cache"
	"gin-example/pkg/secure-cookie"
)

// Amount of time for cookies/redis keys to expire.
const (
	defaultSessionExpire    = 86400 * 1
	defaultSessionMaxLength = 4096
	defaultRedisKeyPrefix   = "session:"
)

// SessionSerializerInterface provides an interface hook for alternative serializers
type SessionSerializerInterface interface {
	Deserialize(d []byte, ss *Session) error
	Serialize(ss *Session) ([]byte, error)
}

// JSONSerializer encode the session map to JSON.
type JSONSerializer struct{}

// Serialize to JSON. Will err if there are unmarshalled key values
func (s JSONSerializer) Serialize(ss *Session) ([]byte, error) {
	m := make(map[string]interface{}, len(ss.Values))
	for k, v := range ss.Values {
		ks, ok := k.(string)
		if !ok {
			return nil, errors.Errorf("Non-string key value, cannot serialize session to JSON: %v", k)
		}
		m[ks] = v
	}
	return json.Marshal(m)
}

// Deserialize back to map[string]interface{}
func (s JSONSerializer) Deserialize(d []byte, ss *Session) error {
	m := make(map[string]interface{})
	err := json.Unmarshal(d, &m)
	if err != nil {
		return err
	}
	for k, v := range m {
		ss.Values[k] = v
	}
	return nil
}

// GobSerializer uses gob package to encode the session map
type GobSerializer struct{}

// Serialize using gob
func (s GobSerializer) Serialize(ss *Session) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(ss.Values)
	if err == nil {
		return buf.Bytes(), nil
	}
	return nil, err
}

// Deserialize back to map[interface{}]interface{}
func (s GobSerializer) Deserialize(d []byte, ss *Session) error {
	dec := gob.NewDecoder(bytes.NewBuffer(d))
	return dec.Decode(&ss.Values)
}

// RedisStore stores sessions in a redis backend.
type RedisStore struct {
	RedisClient cache.SessionCacheRedisClientInterface

	// need by secure cookie
	Codecs  []securecookie.Codec
	Options *Options

	defaultMaxAge    int // default redis TTL for session, 0 No limit
	defaultMaxLength int // default Max redis key size, 0 No limit

	redisKeyPrefix string

	serializer SessionSerializerInterface
}

// NewRedisStore returns a new RedisStore.
func NewRedisStore(redisClientInterface cache.SessionCacheRedisClientInterface, keyPairs ...[]byte) (*RedisStore, error) {
	rs := &RedisStore{
		RedisClient: redisClientInterface,
		Codecs:      securecookie.CodecsFromPairs(keyPairs...),
		Options: &Options{
			Path:   "/",
			MaxAge: defaultSessionExpire,
		},
		defaultMaxAge:    defaultSessionExpire, // 20 minutes seems like a reasonable default
		defaultMaxLength: defaultSessionMaxLength,
		redisKeyPrefix:   defaultRedisKeyPrefix,
		serializer:       GobSerializer{},
	}

	return rs, nil
}

// SetMaxLength sets RedisStore.maxLength if the `maxLength` argument is greater or equal 0
// maxLength restricts the maximum length of new sessions to l.
// If maxLength is 0 there is no limit to the size of a session, use with caution.
// The default for a new RedisStore is 4096. Redis allows for max.
// value sizes of up to 512MB (http://redis.io/topics/data-types)
// Default: 4096,
func (s *RedisStore) SetMaxLength(maxLength int) {
	if maxLength >= 0 {
		s.defaultMaxLength = maxLength
	}
}

// SetRedisKeyPrefix set the redis name prefix
func (s *RedisStore) SetRedisKeyPrefix(prefix string) {
	s.redisKeyPrefix = prefix
}

// SetSerializer sets the serializer
func (s *RedisStore) SetSerializer(i SessionSerializerInterface) {
	s.serializer = i
}

// SetMaxAge restricts the maximum age, in seconds, of the session record
// both in database and a browser. This is to change session storage configuration.
// If you want just to remove session use your session `s` object and change it's
// `Options.MaxAge` to -1, as specified in
//    http://godoc.org/github.com/gorilla/sessions#Options
//
// Default is the one provided by this package value - `sessionExpire`.
// Set it to 0 for no restriction.
// Because we use `MaxAge` also in SecureCookie crypting algorithm you should
// use this function to change `MaxAge` value.
func (s *RedisStore) SetMaxAge(maxAge int) {
	s.defaultMaxAge = maxAge
	//var c *securecookie.SecureCookie
	//var ok bool
	//s.Options.MaxAge = maxAge
	//for i := range s.Codecs {
	//	if c, ok = s.Codecs[i].(*securecookie.SecureCookie); ok {
	//		c.MaxAge(maxAge)
	//	} else {
	//		logging.Logger.Info("Can't change MaxAge", zap.Error(errors.Errorf("codec %v\n", s.Codecs[i])))
	//	}
	//}
}

// Get returns a session for the given name after adding it to the registry.
func (s *RedisStore) Get(r *http.Request, name string) (*Session, error) {
	return GetRegistry(r).Get(s, name)
}

// New returns a session for the given name without adding it to the registry.
func (s *RedisStore) New(r *http.Request, name string) (*Session, error) {
	var (
		err error
		ok  bool
	)
	session := NewSession(s, name)
	ops := *s.Options // make a copy
	ops.MaxAge = s.defaultMaxAge
	session.Options = &ops
	session.IsNew = true
	if c, errCookie := r.Cookie(name); errCookie == nil {
		err = securecookie.DecodeMulti(name, c.Value, &session.ID, s.Codecs...)
		if err == nil {
			ok, err = s.load(session)
			session.IsNew = !(err == nil && ok) // not new if no error and data available
		}
	}
	return session, err
}

// load reads the session from redis.
// returns true if there is a session data in DB
func (s *RedisStore) load(session *Session) (bool, error) {
	data, err := s.RedisClient.Get(context.Background(), s.redisKeyPrefix+session.ID).Bytes()
	if err != nil {
		return false, err
	}
	if data == nil {
		return false, nil // no data was associated with this key
	}

	return true, s.serializer.Deserialize(data, session)
}

// Save adds a single session to the response.
func (s *RedisStore) Save(r *http.Request, w http.ResponseWriter, session *Session) error {
	// Marked for deletion.
	if session.Options.MaxAge <= 0 {
		if err := s.delete(session); err != nil {
			return err
		}
		http.SetCookie(w, NewCookie(session.Name(), "", session.Options))
	} else {
		// Build an alphanumeric key for the redis store.
		if session.ID == "" {
			session.ID = strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32)), "=")
		}
		if err := s.save(session); err != nil {
			return err
		}
		encoded, err := securecookie.EncodeMulti(session.Name(), session.ID, s.Codecs...)
		if err != nil {
			return err
		}
		http.SetCookie(w, NewCookie(session.Name(), encoded, session.Options))
	}
	return nil
}

// save stores the session in redis.
func (s *RedisStore) save(session *Session) error {
	b, err := s.serializer.Serialize(session)
	if err != nil {
		return err
	}

	if s.defaultMaxLength != 0 && len(b) > s.defaultMaxLength {
		return errors.New("SessionStore: the value to store is too big")
	}

	age := session.Options.MaxAge
	if age == 0 {
		age = s.defaultMaxAge
	}
	err = s.RedisClient.Set(context.Background(), s.redisKeyPrefix+session.ID, b, time.Duration(age)*time.Second).Err()

	return nil
}

// Delete removes the session from redis, and sets the cookie to expire.
// WARNING: This method should be considered deprecated since it is not exposed via the gorilla/sessions interface.
// Set session.Options.MaxAge = -1 and call Save instead. - July 18th, 2013
func (s *RedisStore) Delete(r *http.Request, w http.ResponseWriter, session *Session) error {
	if err := s.RedisClient.Del(context.Background(), s.redisKeyPrefix+session.ID).Err(); err != nil {
		return err
	}
	// Set cookie to expire.
	options := *session.Options
	options.MaxAge = -1
	http.SetCookie(w, NewCookie(session.Name(), "", &options))
	// Clear session values.
	for k := range session.Values {
		delete(session.Values, k)
	}
	return nil
}

// delete removes keys from redis if MaxAge<0
func (s *RedisStore) delete(session *Session) error {
	if err := s.RedisClient.Del(context.Background(), s.redisKeyPrefix+session.ID).Err(); err != nil {
		return err
	}
	return nil
}
