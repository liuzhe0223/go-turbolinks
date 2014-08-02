package turbolinks

import (
	"net/http"
	"net/url"
)

type SessionWriter interface {
	Set(string, string)
	Get(string) string
	Del(string)
}

func Before(r *http.Request) {
	referrer := r.Header.Get("X-Xhr-Referer")
	if referrer != "" {
		r.Header.Set("Referer", referrer)
	}
}

func After(r *http.Request, w http.ResponseWriter, session SessionWriter) {
	referrer := r.Header.Get("X-Xhr-Referer")
	if referrer == "" {
		return
	}

	requestMethodCookie, _ := r.Cookie("request_method")
	if requestMethodCookie == nil || requestMethodCookie.Value != r.Method {
		cookie := http.Cookie{
			Name:  "request_method",
			Value: r.Method,
			Path:  "/",
		}

		http.SetCookie(w, &cookie)
	}

	if location := w.Header().Get("Location"); location != "" {
		session.Set("_turbolinks_redirect_to", location)

		if isDiffHost(location, referrer) {
			w.WriteHeader(http.StatusForbidden)
		}
	} else if location := session.Get("_turbolinks_redirect_to"); location != "" {
		session.Del("_turbolinks_redirect_to")
		w.Header().Add("X-Xhr-Redirected-To", location)
	}
}

func isDiffHost(location, referrer string) bool {

	locationURL, err1 := url.Parse(location)
	referrerURL, err2 := url.Parse(referrer)
	if err1 != nil || err2 != nil {
		return false
	}

	if locationURL.Host == "" {
		return false
	}

	return locationURL.Host != referrerURL.Host
}
