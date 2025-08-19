package main

import (
	"encoding/json"
	"html/template"
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
	TargetURL string    `json:"target_url"`
	Findings  []Finding `json:"findings"`
	Score     int       `json:"score"`
	Grade     string    `json:"grade"`
}

func main() {
	tpl := template.Must(template.ParseGlob("templates/*.html"))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Explicitly serve favicon (if stored in ./static/favicon.ico)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/favicon.ico")
	})

	// Home page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			tpl.ExecuteTemplate(w, "index.html", nil)
			return
		}
	})

	// Ping endpoint for uptime monitoring
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Ping received from %s at %s", r.RemoteAddr, time.Now().Format(time.RFC3339))
		w.Write([]byte("pong"))
	})

	// API endpoint for scanning
	http.HandleFunc("/api/scan", func(w http.ResponseWriter, r *http.Request) {
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
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(report)
	})

	// HTML report endpoint
	http.HandleFunc("/scan", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		target := r.FormValue("url")
		report, err := scanSecurityHeaders(target)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tpl.ExecuteTemplate(w, "report.html", report)
	})

	port := ":8080"
	log.Printf("Server running on %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func scanSecurityHeaders(target string) (*Report, error) {
	if !strings.HasPrefix(target, "http") {
		target = "https://" + target
	}
	_, err := url.ParseRequestURI(target)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(target)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	findings := []Finding{}
	score := 100

	checks := map[string]struct {
		expected string
		advice   string
		points   int
	}{
		"Strict-Transport-Security": {
			expected: "present",
			advice:   "Add HSTS to force HTTPS.",
			points:   10,
		},
		"Content-Security-Policy": {
			expected: "present",
			advice:   "Define a strong CSP to prevent XSS.",
			points:   10,
		},
		"X-Frame-Options": {
			expected: "DENY or SAMEORIGIN",
			advice:   "Prevent clickjacking with X-Frame-Options.",
			points:   5,
		},
		"X-Content-Type-Options": {
			expected: "nosniff",
			advice:   "Prevent MIME sniffing.",
			points:   5,
		},
		"Referrer-Policy": {
			expected: "strict or no-referrer",
			advice:   "Control referrer data leakage.",
			points:   5,
		},
	}

	for header, cfg := range checks {
		val := resp.Header.Get(header)
		if val == "" {
			score -= cfg.points
			findings = append(findings, Finding{
				Header: header,
				Issue:  "Missing",
				Advice: cfg.advice,
				Score:  cfg.points,
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
		TargetURL: target,
		Findings:  findings,
		Score:     score,
		Grade:     grade,
	}, nil
}