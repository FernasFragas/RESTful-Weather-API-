package weatherservice

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

var botUserAgents = []string{"bot", "crawl", "spider", "slurp", "curl", "wget", "python-requests", "go-http-client"}

func createVisitsTable() error {
	createTableSQL := `CREATE TABLE IF NOT EXISTS visits (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        visited_at TIMESTAMP NOT NULL,
        visitor TEXT NOT NULL,
        path TEXT NOT NULL,
        city TEXT
    );
    CREATE INDEX IF NOT EXISTS idx_visits_visited_at ON visits (visited_at);`

	if _, err := db.Exec(createTableSQL); err != nil {
		return fmt.Errorf("error creating visits table: %w", err)
	}

	return nil
}

// trackVisit records a page view before handing the request to the next handler.
// Failures are logged and never block the request.
func (s *Server) trackVisit(ctx *fiber.Ctx) error {
	userAgent := ctx.Get(fiber.HeaderUserAgent)
	if db != nil && !isBot(userAgent) {
		// Fly's proxy sets Fly-Client-IP to the real client address.
		ip := ctx.Get("Fly-Client-IP")
		if ip == "" {
			ip = ctx.IP()
		}

		if err := SaveVisit(visitorID(ip, userAgent), ctx.Path(), ctx.FormValue("city_name")); err != nil {
			log.Printf("Error saving visit: %v", err)
		}
	}

	return ctx.Next()
}

// SaveVisit stores a single page view.
func SaveVisit(visitor, path, city string) error {
	insertSQL := `INSERT INTO visits (visited_at, visitor, path, city) VALUES (?, ?, ?, ?);`

	_, err := db.Exec(insertSQL, time.Now().UTC(), visitor, path, strings.ToLower(strings.TrimSpace(city)))
	if err != nil {
		return fmt.Errorf("error saving visit: %w", err)
	}

	return nil
}

// visitorID returns an anonymous, stable identifier so raw IPs are never stored.
func visitorID(ip, userAgent string) string {
	sum := sha256.Sum256([]byte(ip + "|" + userAgent))
	return hex.EncodeToString(sum[:8])
}

func isBot(userAgent string) bool {
	if userAgent == "" {
		return true
	}

	userAgent = strings.ToLower(userAgent)
	for _, bot := range botUserAgents {
		if strings.Contains(userAgent, bot) {
			return true
		}
	}

	return false
}

type DailyStats struct {
	Day            string `json:"day"`
	UniqueVisitors int    `json:"unique_visitors"`
	PageViews      int    `json:"page_views"`
}

type CityStats struct {
	City     string `json:"city"`
	Searches int    `json:"searches"`
}

type VisitStats struct {
	Since          time.Time    `json:"since"`
	UniqueVisitors int          `json:"unique_visitors"`
	PageViews      int          `json:"page_views"`
	Searches       int          `json:"searches"`
	Daily          []DailyStats `json:"daily"`
	TopCities      []CityStats  `json:"top_cities"`
}

// GetVisitStats aggregates the visits recorded since the given time.
func GetVisitStats(since time.Time) (*VisitStats, error) {
	if db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	// Visits are stored as UTC text, so the comparison only works against UTC.
	since = since.UTC()

	stats := &VisitStats{Since: since, Daily: []DailyStats{}, TopCities: []CityStats{}}

	totalsSQL := `SELECT COUNT(DISTINCT visitor), COUNT(*), COUNT(NULLIF(city, ''))
        FROM visits WHERE visited_at >= ?;`
	if err := db.QueryRow(totalsSQL, since).Scan(&stats.UniqueVisitors, &stats.PageViews, &stats.Searches); err != nil {
		return nil, fmt.Errorf("error querying visit totals: %w", err)
	}

	dailySQL := `SELECT substr(visited_at, 1, 10) AS day, COUNT(DISTINCT visitor), COUNT(*)
        FROM visits WHERE visited_at >= ? GROUP BY day ORDER BY day;`
	rows, err := db.Query(dailySQL, since)
	if err != nil {
		return nil, fmt.Errorf("error querying daily visits: %w", err)
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var day DailyStats
		if err := rows.Scan(&day.Day, &day.UniqueVisitors, &day.PageViews); err != nil {
			return nil, fmt.Errorf("error reading daily visits: %w", err)
		}
		stats.Daily = append(stats.Daily, day)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading daily visits: %w", err)
	}

	citiesSQL := `SELECT city, COUNT(*) AS searches FROM visits
        WHERE visited_at >= ? AND city != '' GROUP BY city ORDER BY searches DESC LIMIT 10;`
	cityRows, err := db.Query(citiesSQL, since)
	if err != nil {
		return nil, fmt.Errorf("error querying top cities: %w", err)
	}
	defer func() { _ = cityRows.Close() }()

	for cityRows.Next() {
		var city CityStats
		if err := cityRows.Scan(&city.City, &city.Searches); err != nil {
			return nil, fmt.Errorf("error reading top cities: %w", err)
		}
		stats.TopCities = append(stats.TopCities, city)
	}
	if err := cityRows.Err(); err != nil {
		return nil, fmt.Errorf("error reading top cities: %w", err)
	}

	return stats, nil
}

// showStats returns visit stats as JSON. It is only reachable with the STATS_TOKEN
// set in the environment, e.g. /stats?token=...&days=7
func (s *Server) showStats(ctx *fiber.Ctx) error {
	token := os.Getenv("STATS_TOKEN")
	if token == "" || ctx.Query("token") != token {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	days := ctx.QueryInt("days", 7)
	if days < 1 {
		days = 7
	}

	stats, err := GetVisitStats(time.Now().UTC().AddDate(0, 0, -days))
	if err != nil {
		log.Printf("Error retrieving visit stats with error %s", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.JSON(stats)
}
