package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Finding struct {
	Header string `json:"header"`
	Issue  string `json:"issue"`
	Advice string `json:"advice"`
	Score  int    `json:"score"`
}

type Report struct {
	TargetURL      string    `json:"target_url"`
	Findings       []Finding `json:"findings"`
	MissingHeaders []string  `json:"missing_headers"`
	Score          int       `json:"score"`
	Grade          string    `json:"grade"`
}

func main() {
	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve HTML pages
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.HandleFunc("/report.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "report.html")
	})

	// API endpoints
	http.HandleFunc("/api/scan", scanHandler)
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// Favicon
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/favicon.ico")
	})

	port := ":8080"
	log.Printf("Server running on %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func scanHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	target := r.URL.Query().Get("url")
	if target == "" {
		http.Error(w, "Missing url parameter", http.StatusBadRequest)
		return
	}

	report, err := scanSecurityHeaders(target)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(report)
}

func scanSecurityHeaders(target string) (*Report, error) {
	if !strings.HasPrefix(target, "http") {
		target = "https://" + target
	}

	_, err := url.ParseRequestURI(target)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %v", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("GET", target, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch target: %v", err)
	}
	defer resp.Body.Close()

	// Only essential headers
	checks := map[string]struct {
		Advice string
		Points int
	}{
		"Strict-Transport-Security": {"Forces HTTPS, protects from downgrade attacks.", 10},
		"Content-Security-Policy":   {"Prevents XSS by locking down allowed resources.", 10},
		"X-Frame-Options":           {"Stops clickjacking attacks.", 5},
		"X-Content-Type-Options":    {"Prevents browsers from misinterpreting files.", 5},
		"Referrer-Policy":           {"Controls what sensitive info leaks when users click links.", 5},
	}

	findings := []Finding{}
	missing := []string{}
	score := 100

	for header, cfg := range checks {
		if val := resp.Header.Get(header); val == "" {
			score -= cfg.Points
			missing = append(missing, header)
			findings = append(findings, Finding{
				Header: header,
				Issue:  "Missing",
				Advice: cfg.Advice,
				Score:  cfg.Points,
			})
		}
	}

	grade := "F"
	switch {
	case score >= 90:
		grade = "A"
	case score >= 80:
		grade = "B"
	case score >= 70:
		grade = "C"
	case score >= 60:
		grade = "D"
	}

	return &Report{
		TargetURL:      target,
		Findings:       findings,
		MissingHeaders: missing,
		Score:          score,
		Grade:          grade,
	}, nil
}