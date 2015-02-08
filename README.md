##go-turbolinks


###Martini middleware example

https://gist.github.com/liuzhe0223/5dca0531f5a07bf5cb32

```go
package turbolinks

import (
  "net/http"

  "github.com/go-martini/martini"
  "github.com/liuzhe0223/go-turbolinks/turbolinks"
  "github.com/martini-contrib/sessions"
)

type Session struct {
  session sessions.Session
}

func (s *Session) Get(key string) string {
  val := s.session.Get(key)
  if val == nil {
    return ""
  }

  return val.(string)
}

func (s *Session) Set(key, val string) {
  s.session.Set(key, val)
}

func (s *Session) Del(key string) {
  s.session.Delete(key)
}

func Trubolinks() martini.Handler {
  return func(s sessions.Session, c martini.Context, r *http.Request, w http.ResponseWriter) {
    session := &Session{s}
    rw := w.(martini.ResponseWriter)

    turbolinks.Before(r)

    rw.Before(func(martini.ResponseWriter) {
      turbolinks.After(r, w, session)
    })

    c.Next()
  }
}
```
