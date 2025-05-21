package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	auth "forum/Auth"
	"forum/database"
	handlerfuncitons "forum/handlers"

	_ "github.com/mattn/go-sqlite3"
)

type RateLimiter struct {
	rate      float64
	capacity  float64 
	tokens    float64
	lastCheck time.Time
	mu        sync.Mutex
}

func NewRateLimiter(rate, capacity float64) *RateLimiter {
	return &RateLimiter{
		rate:      rate,
		capacity:  capacity,
		tokens:    capacity,
		lastCheck: time.Now(),
	}
}

func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastCheck).Seconds()
	rl.tokens += elapsed * rl.rate
	if rl.tokens > rl.capacity {
		rl.tokens = rl.capacity
	}
	rl.lastCheck = now

	if rl.tokens >= 1 {
		rl.tokens--
		return true
	}
	return false
}

type IPRateLimiter struct {
	ips map[string]*RateLimiter
	mu  sync.Mutex
	r   float64
	b   float64
}

func NewIPRateLimiter(r float64, b float64) *IPRateLimiter {
	return &IPRateLimiter{
		ips: make(map[string]*RateLimiter),
		r:   r,
		b:   b,
	}
}

func (i *IPRateLimiter) AddIP(ip string) *RateLimiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter := NewRateLimiter(i.r, i.b)
	i.ips[ip] = limiter
	return limiter
}

func (i *IPRateLimiter) GetLimiter(ip string) *RateLimiter {
	i.mu.Lock()
	limiter, exists := i.ips[ip]
	if !exists {
		i.mu.Unlock()
		return i.AddIP(ip)
	}
	i.mu.Unlock()
	return limiter
}

func rateLimitMiddleware(next http.Handler, limiter *IPRateLimiter) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip := r.RemoteAddr
        if !limiter.GetLimiter(ip).Allow() {
            w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%.0f", limiter.b))
            w.Header().Set("X-RateLimit-Remaining", "0")
            w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(time.Second).Unix()))
            http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
            log.Printf("Rate limit exceeded for IP: %s", ip)
            return 
        }
        next.ServeHTTP(w, r)
    })
}


func main() {
	database.Database()
	limiter := NewIPRateLimiter(1, 5) 

	mux := http.NewServeMux()

	// Static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.Handle("/javascript/", http.StripPrefix("/javascript/", http.FileServer(http.Dir("./javascript"))))

	mux.Handle("/", rateLimitMiddleware(http.HandlerFunc(handlerfuncitons.Homepage), limiter))
	mux.Handle("/register", rateLimitMiddleware(http.HandlerFunc(auth.Registerhandler), limiter))
	mux.Handle("/login", rateLimitMiddleware(http.HandlerFunc(auth.LoginHandler), limiter))
	mux.Handle("/logout", rateLimitMiddleware(http.HandlerFunc(auth.LogoutHandler), limiter))
	mux.Handle("/post", rateLimitMiddleware(http.HandlerFunc(handlerfuncitons.Post), limiter))
	mux.Handle("/category", rateLimitMiddleware(http.HandlerFunc(handlerfuncitons.Category), limiter))
	mux.Handle("/welcome", rateLimitMiddleware(http.HandlerFunc(handlerfuncitons.Welcomepage), limiter))
	mux.Handle("/createpost", rateLimitMiddleware(http.HandlerFunc(handlerfuncitons.Createpostshandler), limiter))
	mux.Handle("/like", rateLimitMiddleware(http.HandlerFunc(handlerfuncitons.LikesHandler), limiter))
	mux.Handle("/dislike", rateLimitMiddleware(http.HandlerFunc(handlerfuncitons.DisLikesHandler), limiter))
	mux.Handle("/comment", rateLimitMiddleware(http.HandlerFunc(handlerfuncitons.Comment), limiter))
	mux.Handle("/commentlike", rateLimitMiddleware(http.HandlerFunc(handlerfuncitons.CommentLikeHandler), limiter))
	mux.Handle("/commentdislike", rateLimitMiddleware(http.HandlerFunc(handlerfuncitons.CommentDisLikeHandler), limiter))
	mux.Handle("/filter", rateLimitMiddleware(http.HandlerFunc(handlerfuncitons.FilterHandler), limiter))
	mux.Handle("/createcomment", rateLimitMiddleware(http.HandlerFunc(handlerfuncitons.CreateComment), limiter))
	mux.Handle("/profile", rateLimitMiddleware(http.HandlerFunc(handlerfuncitons.Profilehandler), limiter))
	// mux.Handle("/auth/google/login", rateLimitMiddleware(http.HandlerFunc(auth.GoogleLoginHandler), limiter))
	// mux.Handle("/auth/google/callback", rateLimitMiddleware(http.HandlerFunc(auth.GoogleCallbackHandler), limiter))
	// mux.Handle("/auth/github/login", rateLimitMiddleware(http.HandlerFunc(auth.GitHubHandler), limiter))
	// mux.Handle("/auth/github/callback", rateLimitMiddleware(http.HandlerFunc(auth.GitHubCallbackHandler), limiter))

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		},
	}

	server := &http.Server{
		Addr:      "0.0.0.0:8443",
		TLSConfig: tlsConfig,
		Handler:   mux,
	}

	fmt.Println("Secure server started at https://0.0.0.0:8443")
	err := server.ListenAndServeTLS("server.crt", "server.key")
	if err != nil {
		fmt.Println("Failed to start HTTPS server:", err)
	}
}
