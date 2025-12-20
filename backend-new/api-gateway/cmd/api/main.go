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
	// Define target URLs for services
	authURL, _ := url.Parse("http://localhost:8081")
	userURL, _ := url.Parse("http://localhost:8082")
	postURL, _ := url.Parse("http://localhost:8083")
	mediaURL, _ := url.Parse("http://localhost:8084")
	chatURL, _ := url.Parse("http://localhost:8085")
	groupURL, _ := url.Parse("http://localhost:8086")

	// Helper to create proxy with response modification
	newProxy := func(target *url.URL) *httputil.ReverseProxy {
		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.ModifyResponse = func(resp *http.Response) error {
			// Remove CORS headers from backend response to avoid duplicates
			resp.Header.Del("Access-Control-Allow-Origin")
			resp.Header.Del("Access-Control-Allow-Methods")
			resp.Header.Del("Access-Control-Allow-Headers")
			resp.Header.Del("Access-Control-Allow-Credentials")
			return nil
		}
		return proxy
	}

	// Create proxies
	authProxy := newProxy(authURL)
	userProxy := newProxy(userURL)
	postProxy := newProxy(postURL)
	mediaProxy := newProxy(mediaURL)
	chatProxy := newProxy(chatURL)
	groupProxy := newProxy(groupURL)

	// Main handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// CORS setup
		origin := r.Header.Get("Origin")
		// Allow any localhost origin for development
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

		// Routing logic
		switch {
		// Auth Service Routes
		case path == "/login" || path == "/register" || path == "/logout" || path == "/checkstate":
			authProxy.ServeHTTP(w, r)

		// User Service Routes
		case path == "/profile" ||
			strings.HasPrefix(path, "/profile/") ||
			path == "/toggle-follow" ||
			path == "/request_follow" ||
			path == "/Friends" ||
			path == "/users-sug" ||
			path == "/update-privacy":
			userProxy.ServeHTTP(w, r)

		// Post Service Routes
		case path == "/Createpost" ||
			path == "/getposts" ||
			strings.HasPrefix(path, "/GetCUser/") ||
			strings.HasPrefix(path, "/posts/") ||
			path == "/reactions":
			postProxy.ServeHTTP(w, r)

		// Media Service Routes
		case strings.HasPrefix(path, "/uploads/"):
			mediaProxy.ServeHTTP(w, r)

		// Chat Service Routes
		case path == "/ws" ||
			path == "/notifications" ||
			path == "/chat/ws" ||
			strings.HasPrefix(path, "/chat/"):
			chatProxy.ServeHTTP(w, r)

		// Group Service Routes
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
