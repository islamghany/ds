package main

import (
	"encoding/json"
	"io/ioutil"
	"islamghany/ds/rate-limiting/rateLimitingStrategies"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
)

// unit can be second, minute, hour, day
type Unit string

const (
	Second Unit = "second"
	Minute Unit = "minute"
	Hour   Unit = "hour"
	Day    Unit = "day"
)

type RateLimit struct {
	Unit           Unit `json:"unit"`
	RequestPerUnit int  `json:"request_per_unit"`
}
type Rule struct {
	Path      string    `json:"path"`
	RateLimit RateLimit `json:"rate_limit"`
}
type LimiterConfig struct {
	Origin string `json:"origin"`
	Rules  []Rule `json:"rules"`
}

var limiterConfig LimiterConfig

func getMatchingRules(path string) []Rule {
	log.Print("path: ", path)
	var matchingRules []Rule
	for _, rule := range limiterConfig.Rules {
		matched, err := regexp.MatchString(rule.Path, path)
		if err != nil {
			log.Fatal(err)
		}
		if matched {
			matchingRules = append(matchingRules, rule)
		}
	}

	return matchingRules
}

type IPPaths map[string]rateLimitingStrategies.Limiter

func createIPaths(rule Rule) IPPaths {
	ipPaths := make(IPPaths)
	unitInSeconds := 0
	if rule.RateLimit.Unit == Second {
		unitInSeconds = rule.RateLimit.RequestPerUnit
	} else if rule.RateLimit.Unit == Minute {
		unitInSeconds = rule.RateLimit.RequestPerUnit * 60
	} else if rule.RateLimit.Unit == Hour {
		unitInSeconds = rule.RateLimit.RequestPerUnit * 60 * 60
	} else if rule.RateLimit.Unit == Day {
		unitInSeconds = rule.RateLimit.RequestPerUnit * 60 * 60 * 24
	}

	ipPaths[rule.Path] = *rateLimitingStrategies.NewLimiter(rateLimitingStrategies.NewtokenBucket(rule.RateLimit.RequestPerUnit, unitInSeconds))
	return ipPaths
}

type IPMap map[string]IPPaths

func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}
func main() {
	// Todo: run a background goroutine that periodically cleans up the IPMap
	// phase 1: read the limiter.config.json file
	limiterConfigFile, err := os.Open("limiter.config.json")
	if err != nil {
		log.Fatal(err)
	}
	limiterConfigBytes, err := ioutil.ReadAll(limiterConfigFile)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(limiterConfigBytes, &limiterConfig)
	if err != nil {
		log.Fatal(err)
	}
	limiterConfigFile.Close()
	// phase 2: create a map of origin -> path -> limiter
	ipMap := make(IPMap)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		ip := GetIP(r)
		log.Panicln("IP: ", ip, path)
		matchedRules := getMatchingRules(path)
		// log.Printf("Matching rule: %v", matchedRules)
		if len(matchedRules) == 0 {
			log.Printf("No matching rules for path %s", path)
			w.WriteHeader(http.StatusAccepted)
			return
		}
		for _, rule := range matchedRules {
			ipMatcher, ok := ipMap[ip]
			if !ok {
				ipMatcher = createIPaths(rule)
				ipMap[ip] = ipMatcher
			}
			limiter, ok := ipMatcher[path]
			if !ok {
				limiter = *rateLimitingStrategies.NewLimiter(rateLimitingStrategies.NewtokenBucket(rule.RateLimit.RequestPerUnit, rule.RateLimit.RequestPerUnit))
				ipMatcher[path] = limiter
			}
			if !limiter.AddRequest() {
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
			http.Redirect(w, r, "https://www.google.com", http.StatusFound)
		}
		// log.Printf("======================================================")
		// log.Printf("Query-string: %s", r.URL.RawQuery)
		// log.Printf("Path: %s", r.URL.Path)
		// log.Printf("Method: %s", r.Method)
		// log.Printf("Path: %s", r.Host)
		// for k, v := range r.Header {
		// 	log.Printf("Header %s=%s", k, v)
		// }
		// if r.Body != nil {
		// 	body, _ := ioutil.ReadAll(r.Body)
		// 	log.Printf("Body: %s", string(body))
		// }
		// log.Printf("======================================================")
		// w.WriteHeader(http.StatusAccepted)
	})
	srv := &http.Server{
		Addr:           ":8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(srv.ListenAndServe())
}
