package ginsessions

import (
	"net/http"

	"gin-example/pkg/sessions"
)

// Options stores configuration for a session or session store.
// Fields are a subset of http.Cookie fields.
type Options struct {
	Path   string
	Domain string
	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'.
	// MaxAge>0 means Max-Age attribute present and given in seconds.
	MaxAge   int
	Secure   bool
	HttpOnly bool
	// rfc-draft to preventing CSRF: https://tools.ietf.org/html/draft-west-first-party-cookies-07
	//   refer: https://godoc.org/net/http
	//          https://www.sjoerdlangkemper.nl/2016/04/14/preventing-csrf-with-samesite-cookie-attribute/
	SameSite http.SameSite
}

func (o Options) ToOptions() *sessions.Options {
	return &sessions.Options{
		Path:     o.Path,
		Domain:   o.Domain,
		MaxAge:   o.MaxAge,
		Secure:   o.Secure,
		HttpOnly: o.HttpOnly,
		SameSite: o.SameSite,
	}
}
