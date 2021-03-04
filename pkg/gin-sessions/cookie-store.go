package ginsessions

//import (
//	"gin-example/pkg/gin-sessions/sessions"
//)
//
//type cookieStore struct {
//	*sessions.CookieStore
//}
//
//func NewCookieStore(keyPairs ...[]byte) GinStoreInterface {
//	return &cookieStore{
//		CookieStore: sessions.NewCookieStore(keyPairs...),
//	}
//}
//
//func (c *cookieStore) CovertOptions(options Options) {
//	c.CookieStore.Options = options.ToOptions()
//}
