package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func main() {
	authURL, _ := url.Parse("http://localhost:8081")
	userURL, _ := url.Parse("http://localhost:8082")
	postURL, _ := url.Parse("http://localhost:8083")
	mediaURL, _ := url.Parse("http://localhost:8084")
	chatURL, _ := url.Parse("http://localhost:8085")
	groupURL, _ := url.Parse("http://localhost:8086")

	newProxy := func(target *url.URL) *httputil.ReverseProxy {
		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.ModifyResponse = func(resp *http.Response) error {
			resp.Header.Del("Access-Control-Allow-Origin")
			resp.Header.Del("Access-Control-Allow-Methods")
			resp.Header.Del("Access-Control-Allow-Headers")
			resp.Header.Del("Access-Control-Allow-Credentials")
			return nil
		}
		return proxy
	}

	authProxy := newProxy(authURL)
	userProxy := newProxy(userURL)
	postProxy := newProxy(postURL)
	mediaProxy := newProxy(mediaURL)
	chatProxy := newProxy(chatURL)
	groupProxy := newProxy(groupURL)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		origin := r.Header.Get("Origin")
		if origin != "" && (strings.HasPrefix(origin, "http://localhost") || strings.HasPrefix(origin, "http://127.0.0.1")) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		switch {
		case path == "/login" || path == "/register" || path == "/logout" || path == "/checkstate":
			authProxy.ServeHTTP(w, r)

		case path == "/profile" ||
			strings.HasPrefix(path, "/profile/") ||
			path == "/toggle-follow" ||
			path == "/request_follow" ||
			path == "/Friends" ||
			path == "/users-sug" ||
			path == "/followers" ||
			path == "/following" ||
			path == "/update-privacy" ||
			path == "/update-profile":
			userProxy.ServeHTTP(w, r)

		case path == "/Createpost" ||
			path == "/getposts" ||
			strings.HasPrefix(path, "/GetCUser/") ||
			strings.HasPrefix(path, "/posts/") ||
			path == "/reactions":
			postProxy.ServeHTTP(w, r)

		case strings.HasPrefix(path, "/uploads/"):
			mediaProxy.ServeHTTP(w, r)

		case path == "/ws" ||
			path == "/notifications" ||
			path == "/chat/ws" ||
			strings.HasPrefix(path, "/chat/"):
			chatProxy.ServeHTTP(w, r)

		case path == "/Create_Group" ||
			path == "/Get_Groups" ||
			path == "/join" ||
			path == "/accept-reject-join" ||
			path == "/Get_Group_By_ID" ||
			path == "/CreateGroupPost" ||
			path == "/GetGroupPosts" ||
			strings.HasPrefix(path, "/groups/") ||
			strings.HasPrefix(path, "/group/") ||
			path == "/create_event" ||
			path == "/group-events" ||
			path == "/event-invite-response":
			groupProxy.ServeHTTP(w, r)

		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	})

	port := ":8080"
	fmt.Printf("API Gateway running on %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
