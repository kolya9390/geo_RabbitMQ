package reverproxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type ReverseProxy struct {
	host string
	port string
}

func NewReverseProxy(host, port string) *ReverseProxy {
	return &ReverseProxy{
		host: host,
		port: port,
	}
}

func (rp *ReverseProxy) ReverseProxy(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if strings.HasPrefix(r.URL.Path, "/api") {

			next.ServeHTTP(w, r)
			return
		}

		if strings.HasPrefix(r.URL.Path, "/address") {

			next.ServeHTTP(w, r)
			return
		}


		if strings.HasPrefix(r.URL.Path, "/swagger") {

			next.ServeHTTP(w, r)
			return
		}

		if strings.HasPrefix(r.URL.Path, "/public") {

			next.ServeHTTP(w, r)
			return
		}

		linc := fmt.Sprintf("http://%s:%s", rp.host, rp.port)

		uri, _ := url.Parse(linc)

		if uri.Host == r.Host {
			next.ServeHTTP(w, r)
			return
		}
		r.Header.Set("Reverse-Proxy", "true")

		proxy := httputil.ReverseProxy{Director: func(r *http.Request) {
			r.URL.Scheme = uri.Scheme
			r.URL.Host = uri.Host
			r.URL.Path = uri.Path + r.URL.Path
			r.Host = uri.Host
		}}

		proxy.ServeHTTP(w, r)
	})

}
