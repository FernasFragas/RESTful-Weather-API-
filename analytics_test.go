package weatherservice

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const browserUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 14_5) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Safari/605.1.15"

// clearVisits empties the visits table. Tracking writes synchronously, so no earlier test
// can add rows afterwards.
func clearVisits(t *testing.T) {
	t.Helper()

	_, err := db.Exec(`DELETE FROM visits`)
	require.NoError(t, err)
}

// insertVisit stores a visit at a fixed time, the same way SaveVisit stores the current time.
func insertVisit(t *testing.T, at time.Time, visitor, path, city string) {
	t.Helper()

	_, err := db.Exec(`INSERT INTO visits (visited_at, visitor, path, city) VALUES (?, ?, ?, ?)`, at.UTC(), visitor, path, city)
	require.NoError(t, err)
}

type storedVisit struct {
	visitor string
	path    string
	city    string
}

func storedVisits(t *testing.T) []storedVisit {
	t.Helper()

	rows, err := db.Query(`SELECT visitor, path, city FROM visits ORDER BY id`)
	require.NoError(t, err)
	defer func() { _ = rows.Close() }()

	var visits []storedVisit
	for rows.Next() {
		var visit storedVisit
		require.NoError(t, rows.Scan(&visit.visitor, &visit.path, &visit.city))
		visits = append(visits, visit)
	}
	require.NoError(t, rows.Err())

	return visits
}

// Tests for GetVisitStats
func TestGetVisitStats(t *testing.T) {
	clearVisits(t)

	noon := time.Now().UTC().Truncate(24 * time.Hour).Add(12 * time.Hour)
	twoDaysAgo, yesterday := noon.AddDate(0, 0, -2), noon.AddDate(0, 0, -1)

	insertVisit(t, twoDaysAgo, "a", "/", "")
	insertVisit(t, twoDaysAgo, "a", "/process-form/", "porto")
	insertVisit(t, twoDaysAgo, "b", "/", "")
	insertVisit(t, yesterday, "a", "/", "")
	insertVisit(t, yesterday, "c", "/process-form/", "porto")
	insertVisit(t, yesterday, "c", "/process-form/", "faro")
	// Outside the 7-day window.
	insertVisit(t, noon.AddDate(0, 0, -8), "d", "/process-form/", "lisbon")

	stats, err := GetVisitStats(noon.AddDate(0, 0, -7))
	require.NoError(t, err)

	assert.Equal(t, 3, stats.UniqueVisitors)
	assert.Equal(t, 6, stats.PageViews)
	assert.Equal(t, 3, stats.Searches)
	assert.Equal(t, []DailyStats{
		{Day: twoDaysAgo.Format("2006-01-02"), UniqueVisitors: 2, PageViews: 3},
		{Day: yesterday.Format("2006-01-02"), UniqueVisitors: 2, PageViews: 3},
	}, stats.Daily)
	assert.Equal(t, []CityStats{
		{City: "porto", Searches: 2},
		{City: "faro", Searches: 1},
	}, stats.TopCities)
}

func TestGetVisitStats_NoVisits(t *testing.T) {
	clearVisits(t)

	stats, err := GetVisitStats(time.Now().AddDate(0, 0, -7))
	require.NoError(t, err)

	assert.Zero(t, stats.UniqueVisitors)
	assert.Zero(t, stats.PageViews)
	// Empty lists rather than nil, so the JSON has [] instead of null.
	assert.NotNil(t, stats.Daily)
	assert.NotNil(t, stats.TopCities)
}

func TestGetVisitStats_NonUTCSince(t *testing.T) {
	clearVisits(t)

	visitedAt := time.Now().UTC().Add(-2 * time.Hour)
	insertVisit(t, visitedAt, "a", "/", "")

	// The same instant must count the same visits whatever time zone it is expressed in.
	since := visitedAt.Add(-time.Hour)
	for _, loc := range []*time.Location{time.UTC, time.FixedZone("UTC+5", 5*60*60), time.FixedZone("UTC-8", -8*60*60)} {
		stats, err := GetVisitStats(since.In(loc))
		require.NoError(t, err)
		assert.Equal(t, 1, stats.PageViews, "since expressed in %s", loc)
	}
}

// Tests for showStats
func TestShowStats_Token(t *testing.T) {
	server := NewAppServer(nil, nil, nil, nil)

	tests := []struct {
		name       string
		token      string
		url        string
		wantStatus int
	}{
		{"no token configured", "", "/stats?token=", fiber.StatusNotFound},
		{"missing token", "secret", "/stats", fiber.StatusNotFound},
		{"wrong token", "secret", "/stats?token=guess", fiber.StatusNotFound},
		{"right token", "secret", "/stats?token=secret", fiber.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("STATS_TOKEN", tt.token)

			status, _ := doRequest(t, server, httptest.NewRequest("GET", tt.url, nil))

			assert.Equal(t, tt.wantStatus, status)
		})
	}
}

func TestShowStats_Days(t *testing.T) {
	t.Setenv("STATS_TOKEN", "secret")
	clearVisits(t)

	now := time.Now().UTC()
	insertVisit(t, now.Add(-time.Hour), "a", "/process-form/", "porto")
	insertVisit(t, now.AddDate(0, 0, -3), "b", "/", "")

	tests := []struct {
		name      string
		days      string
		wantSince time.Time
		wantViews int
	}{
		{"default window", "", now.AddDate(0, 0, -7), 2},
		{"custom window", "2", now.AddDate(0, 0, -2), 1},
		{"zero falls back to 7", "0", now.AddDate(0, 0, -7), 2},
		{"negative falls back to 7", "-3", now.AddDate(0, 0, -7), 2},
		{"non-numeric falls back to 7", "abc", now.AddDate(0, 0, -7), 2},
	}

	server := NewAppServer(nil, nil, nil, nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, body := doRequest(t, server, httptest.NewRequest("GET", "/stats?token=secret&days="+tt.days, nil))
			require.Equal(t, fiber.StatusOK, status)

			var stats VisitStats
			require.NoError(t, json.Unmarshal([]byte(body), &stats))
			assert.WithinDuration(t, tt.wantSince, stats.Since, time.Minute)
			assert.Equal(t, tt.wantViews, stats.PageViews)
		})
	}
}

// Tests for trackVisit
func TestTrackVisit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	server, mockWeather, _, _ := createTestServer(ctrl)

	// The page itself fails fast here. Tracking runs before it on the page routes.
	mockWeather.EXPECT().GenerateReport(gomock.Any(), gomock.Any()).Return(nil, errors.New("weather unavailable")).AnyTimes()

	visit := func(target, userAgent, clientIP string) {
		req := httptest.NewRequest("GET", target, nil)
		if userAgent != "" {
			req.Header.Set("User-Agent", userAgent)
		}
		if clientIP != "" {
			req.Header.Set("Fly-Client-IP", clientIP)
		}
		doRequest(t, server, req)
	}

	t.Run("home page visit is recorded", func(t *testing.T) {
		clearVisits(t)

		visit("/", browserUserAgent, "203.0.113.7")

		assert.Equal(t, []storedVisit{{visitorID("203.0.113.7", browserUserAgent), "/", ""}}, storedVisits(t))
	})

	t.Run("search is recorded with a normalized city", func(t *testing.T) {
		clearVisits(t)

		visit("/process-form/?city_name=%20Porto%20", browserUserAgent, "203.0.113.7")

		assert.Equal(t, []storedVisit{{visitorID("203.0.113.7", browserUserAgent), "/process-form/", "porto"}}, storedVisits(t))
	})

	t.Run("Fly-Client-IP identifies the visitor", func(t *testing.T) {
		clearVisits(t)

		visit("/", browserUserAgent, "203.0.113.7")
		visit("/", browserUserAgent, "203.0.113.7")
		visit("/", browserUserAgent, "198.51.100.23")

		stats, err := GetVisitStats(time.Now().Add(-time.Hour))
		require.NoError(t, err)
		assert.Equal(t, 2, stats.UniqueVisitors)
		assert.Equal(t, 3, stats.PageViews)
	})

	t.Run("bots and non-page requests are not recorded", func(t *testing.T) {
		clearVisits(t)

		visit("/", "", "203.0.113.7")
		visit("/", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)", "203.0.113.7")
		visit("/", "curl/8.7.1", "203.0.113.7")
		visit("/styles.css", browserUserAgent, "203.0.113.7")
		visit("/stats", browserUserAgent, "203.0.113.7")

		assert.Empty(t, storedVisits(t))
	})
}

func TestTrackVisit_NoDatabase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	server, mockWeather, _, _ := createTestServer(ctrl)

	// The page handler being called proves tracking passed the request on.
	mockWeather.EXPECT().GenerateReport(gomock.Any(), "Porto").Return(nil, errors.New("weather unavailable"))

	shared := db
	db = nil
	t.Cleanup(func() { db = shared })

	req := httptest.NewRequest("GET", "/?city_name=Porto", nil)
	req.Header.Set("User-Agent", browserUserAgent)
	doRequest(t, server, req)
}

func TestIsBot(t *testing.T) {
	tests := []struct {
		userAgent string
		want      bool
	}{
		{"", true},
		{"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)", true},
		{"Mozilla/5.0 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)", true},
		{"curl/8.7.1", true},
		{"python-requests/2.32.3", true},
		{"Go-http-client/1.1", true},
		{browserUserAgent, false},
		{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36", false},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, isBot(tt.userAgent), "user agent %q", tt.userAgent)
	}
}

func TestVisitorID(t *testing.T) {
	id := visitorID("203.0.113.7", browserUserAgent)

	assert.Regexp(t, `^[0-9a-f]{16}$`, id)
	assert.Equal(t, id, visitorID("203.0.113.7", browserUserAgent))
	assert.NotEqual(t, id, visitorID("198.51.100.23", browserUserAgent))
	assert.NotEqual(t, id, visitorID("203.0.113.7", "curl/8.7.1"))
}
