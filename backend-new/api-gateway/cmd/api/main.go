package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	authURL, _ := url.Parse(getEnv("AUTH_SERVICE_URL", "http://localhost:8081"))
	userURL, _ := url.Parse(getEnv("USER_SERVICE_URL", "http://localhost:8082"))
	postURL, _ := url.Parse(getEnv("POST_SERVICE_URL", "http://localhost:8083"))
	mediaURL, _ := url.Parse(getEnv("MEDIA_SERVICE_URL", "http://localhost:8084"))
	chatURL, _ := url.Parse(getEnv("CHAT_SERVICE_URL", "http://localhost:8085"))
	groupURL, _ := url.Parse(getEnv("GROUP_SERVICE_URL", "http://localhost:8086"))
	notificationURL, _ := url.Parse(getEnv("NOTIFICATION_SERVICE_URL", "http://localhost:8087"))

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
	notificationProxy := newProxy(notificationURL)

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
			path == "/GetPost" ||
			path == "/getposts" ||
			strings.HasPrefix(path, "/GetCUser/") ||
			strings.HasPrefix(path, "/posts/") ||
			path == "/reactions":
			postProxy.ServeHTTP(w, r)

		case strings.HasPrefix(path, "/uploads/"):
			mediaProxy.ServeHTTP(w, r)

		case path == "/chat/ws" ||
			strings.HasPrefix(path, "/chat/"):
			chatProxy.ServeHTTP(w, r)

		case path == "/ws" ||
			path == "/notifications" ||
			strings.HasPrefix(path, "/notifications/") ||
			path == "/notifications/read" ||
			path == "/notifications/send":
			notificationProxy.ServeHTTP(w, r)

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
	log.Printf("API Gateway running on %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
