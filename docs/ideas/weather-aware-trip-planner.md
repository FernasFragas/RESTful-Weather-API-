# Weather-aware trip planner

_Idea one-pager, written 2026-09-19. Replaces the old itinerary (Foursquare) and hotels (Google Places) features._

## TL;DR

- **What:** enter a city, a start date and a number of days. You get a **day-by-day plan that follows the forecast** (museums on the rainy day, viewpoints on the sunny one) and **the best area to stay for that plan**.
- **Cost:** $0, forever. Only open data, with no API keys that can bill.
- **Why it's worth building:** TravelTab is a weather app, and no other trip planner plans around the forecast. That's the portfolio story.
- **Next (v2):** shareable trip links, calendar and map export, and city intros written ahead of time by a local AI. See [Next Iteration (v2)](#next-iteration-v2).

## Problem Statement

How might we help someone who has already decided on a trip go from "I'm going to Lisbon" to "I know what to do each day and where to stay", using only data that is free forever?

## Recommended Direction

**Build the plan from famous places, split it into days by area, and let the weather choose which day gets what.**

1. **Places:** Wikidata gives the landmarks near the city. They're ranked by how many Wikipedia languages cover each place, which is a strong "this is famous" signal. That fixes the old problem of random results. Photos come from Wikimedia Commons.
2. **Days:** nearby places are grouped into N days, with 3–4 stops per day, and each day is ordered as a walking loop. This is plain Go with no API.
3. **Weather:** the Open-Meteo daily forecast decides the order. Rainy days get mostly indoor stops and sunny days get mostly outdoor ones.
4. **Stay:** find the center of the plan and show hotels, hostels and guest houses within 1 km from OpenStreetMap. Each gets a "Check prices" link to a Booking search with your dates filled in.

Hotels become part of the plan instead of a separate list. We don't show prices; we help you choose an area, then hand you off to book.

## Directions considered

| Direction | User value | Effort | What makes it different | Verdict |
|---|---|---|---|---|
| **A. Weather-aware planner** (this doc) | Solves a real planning problem | Medium: the hard part is grouping places into days | High: nobody plans around the forecast | ✅ **Build this** |
| B. Curated city guide (Wikivoyage + AI intros written in advance) | Inspiration, but these users have already decided to go | Medium: parsing wiki markup is fiddly | Low: it would be a worse Wikivoyage | ❌ As the main feature. The AI intros come back as a v2 add-on |
| C. Top 10 must-sees only | Quick to read, but not a plan | Low: about a weekend | Low: every travel site has this | ➖ Becomes **step 1** of A |

## Free data we use

| What | Source | Key? | Limits | We must |
|---|---|---|---|---|
| Places + ranking | Wikidata query service | No | Keep queries small; send a `User-Agent` | Cache per city |
| Photos | Wikimedia Commons | No | — | Credit each photo (licenses vary per image) |
| Forecast | Open-Meteo | No | 10k calls/day, **non-commercial only**, **16 days ahead max** | Credit "Open-Meteo" (CC BY 4.0) |
| Hotels | OpenStreetMap via the Overpass API | No | ~10k requests/day; **the public server says it isn't meant as an app backend** | Cache for 30 days; credit "© OpenStreetMap contributors" |
| Booking hand-off | A plain Booking.com search URL | No | — | Nothing. It's just a link |

## Key Assumptions to Validate

Check each one with a quick script **before** building the UI.

- [ ] **Ranking by Wikipedia coverage gives good places.** Test: fetch the top 15 for Lisbon, Porto, Tavira, Funchal and Kyoto. At least 12 of 15 should be places you'd actually visit.
- [ ] **We can tell indoor from outdoor.** Test: map the Wikidata place types (museum, church, park, beach and so on) to indoor or outdoor for the same 5 cities. At least 80% of places should get a label.
- [ ] **The forecast is useful even though it only covers 16 days.** People often plan further ahead. Test: ask 3–5 friends when they decide what to do each day of a trip. If most say "the week before", this holds up.
- [ ] **Light, cached use of the public Overpass server is acceptable.** Test: count new cities per day in `/stats`. If it's well under a few hundred, cache and move on. If it's more, switch hotels to Wikivoyage's "Sleep" listings.
- [ ] **OpenStreetMap has enough hotels near the plan center.** Test: in the 5 cities, count hotels within 1 km that have a website. Aim for at least 5.

## MVP Scope

**In:**

- [ ] One **"Plan my trip"** card under the weather, with a start date and **number of days (1–5)**. **No category checkboxes**, because they cluttered the old form.
- [ ] Wikidata places within 10 km: the top ~20 by fame, each with a photo and an indoor/outdoor label.
- [ ] Grouping into days: 3–4 stops per day, ordered as a walking loop, with the walking distance shown.
- [ ] Weather rule: if the rain chance is ≥ 60%, make it an indoor day; otherwise an outdoor day. If the trip starts more than 16 days away, build the plan without weather and show "Forecast appears 16 days before your trip".
- [ ] Stay card: up to 6 places to stay within 1 km of the plan center, showing stars, website and a "Check prices" link with the dates filled in.
- [ ] SQLite cache: places and hotels kept 30 days per city, the forecast kept 3 hours.
- [ ] A credits line in the footer for OpenStreetMap, Wikidata/Commons and Open-Meteo.
- [ ] Tests with stubbed APIs, like the existing ones, running in CI.
- [ ] Remove the Google Places and Foursquare code and keys.

**Done when:**

- [ ] "Lisbon, 3 days, next week" gives a sensible plan, and the rainy day really does get the museums.
- [ ] A cached city loads in under 2 seconds.
- [ ] Monthly API bill: **$0**.
- [ ] The README has a short "how the planner works" section with a diagram, ready for a portfolio.

## Not Doing (and Why)

- **Hotel prices, availability or reviews.** Only booking companies have that data, and it costs money. We hand off with a link.
- **AI in the MVP.** Simple rules can build the plan. Running AI on Fly needs a bigger paid machine. City intros written ahead of time on a laptop come in v2.
- **Category checkboxes.** They made the old form cluttered. Fame and the weather decide instead.
- **Restaurants.** OpenStreetMap has too many, with no ranking, which would bring back the "bad results" problem.
- **Weather for trips more than 16 days away.** Typical weather for those dates (averages from past years) can come later.
- **Real walking routes.** Straight-line distance is good enough for now. The free public routing server isn't meant for production use.
- **Parsing Wikivoyage listings.** Coverage is uneven and the markup is a project of its own. It's only the fallback if Overpass becomes a problem. (The v2 intros only need the page's plain text.)
- **Accounts or saved trips.** Too much for the MVP. Shareable links and export come in v2 and cover most of this.
- **Running our own Overpass server.** Only worth it if traffic grows.

## Open Questions

- [ ] **Is the old hotels card still costing money?** It still calls Google Places if `PLACES_API_NEW` is set on Fly. Check with `fly secrets list`, and remove the key if you don't need it.
- [ ] Maximum trip length for the MVP: 3 days or 5?
- [ ] Where does the planner live: under the weather on the same page, or on its own `/trip` page? (v2 gives every plan its own `/trip` URL either way.)
- [ ] Is "rain chance ≥ 60%" the right cut-off for a rainy day, or should it be an amount of rain in mm?
- [ ] Small towns with fewer than 8 good places: widen the search area, or suggest fewer days?

## Next Iteration (v2)

Three add-ons once the MVP works. **Build them in this order**, because each one uses the one before:

1. **Shareable links** give every plan a stable URL.
2. **Export** turns that same URL into a calendar file or a map file.
3. **AI city intros** go at the top of the trip page.

| Add-on | User value | Effort | Main risk |
|---|---|---|---|
| Shareable links | High: send the plan to whoever you're traveling with | Low | The plan has to come out the same every time |
| Export | High: the plan ends up in your calendar or map | Low | Time zones |
| AI city intros | Medium: nice to read, but the plan works without it | Medium | The model making things up |

### 1. Shareable trip links

`/trip/lisbon-pt?days=3&from=2026-10-02` always opens the same plan.

- [ ] Add a `GET /trip/:city` page. Submitting the planner puts this URL in the address bar (HTMX `hx-push-url`), so copying the address shares the plan.
- [ ] Make the plan deterministic: the same input always gives the same stops on the same days. No randomness when grouping; break ties with the Wikidata ID.
- [ ] Search basics: a page title and description ("3 days in Lisbon: …"), a canonical URL **without the date**, and a `sitemap.xml` that lists only cities already in the cache.

**Watch out:**

- **The weather changes, so the plan can't be fully frozen.** The stops and how they're grouped into days stay the same. Only the order of the days follows the latest forecast, with a note saying "Days reordered for the latest forecast". A friend opening your link next week gets fresher weather, which is a feature, not a bug.
- **Search engine crawlers** opening cities that aren't cached would hit Wikidata and Overpass. Add a simple limit on how many uncached cities can be fetched per minute.

### 2. Export the trip

Three buttons under the plan:

- [ ] **Add to calendar (`.ics`):** one event per stop at default times (for example 10:00, 12:00, 15:00 and 17:00, 1.5 hours each), with the place's location and a link back to the trip. Only shown when the trip has a start date.
- [ ] **Google My Maps (`.kml`):** one folder per day and one pin per stop. You import it at mymaps.google.com.
- [ ] **Bonus: "Open Day 1 in Google Maps":** a plain Google Maps link with walking directions through that day's stops. It needs no key, costs nothing, and is the most useful of the three on a phone. With 3–4 stops a day, it stays within Google's mobile limit of 3 stops between the start and the end.

The files are served at `/trip/lisbon-pt.ics?…` and `/trip/lisbon-pt.kml?…`: the same URL as the page, with a different extension.

**Watch out:**

- **Time zones.** Get the city's time zone from Open-Meteo (`timezone=auto`) and write event times in UTC. Import Go's `time/tzdata` so the slim Docker image doesn't need time zone files.
- **`.ics` files have a few formatting rules** (escape commas and semicolons, wrap long lines). Test the file in Google Calendar, Apple Calendar and Outlook.

### 3. City intros, written ahead of time by a local AI

A short intro at the top of the trip page: why go, what each area is like, and how to get around. It's written on your laptop and costs nothing to serve.

- [ ] Turn the empty `cmd/console` into `cmd/guides`. For each city in the list, it fetches the Wikivoyage page text and asks a small model in Ollama (7–8B is plenty) for a ~120-word intro **using only that text**.
- [ ] **Check for made-up places:** reject and retry any intro that names a place that isn't in the Wikivoyage text.
- [ ] Save the intros to `guides/guides.json` in the repo, embed the file in the app with `go:embed`, and load it at startup.
- [ ] Store the Wikivoyage revision for each city, so a rerun only rewrites cities whose page changed.
- [ ] Show "Based on Wikivoyage (CC BY-SA)" with a link under each intro. The license requires it.

**Why a JSON file in the repo instead of SQLite:** the database lives on the Fly volume, not on your laptop, so there's no simple way to copy rows into it. A file in git ships with every deploy, and **you can read every intro the AI wrote in the diff before it goes live.**

**Watch out:**

- **Write a general city intro, not "3 days in Lisbon".** The plan can be 1–5 days long, so a fixed "3 days" text would contradict it.
- **Which 100 cities?** `/stats` only returns the top 10, and it only started counting recently. Start with a hand-picked list in `guides/cities.txt`, then add your most-searched cities over time.

### v2 assumptions to validate

- [ ] **The plan comes out the same every time.** Test: build "Lisbon, 3 days" twice, a day apart and after a cache refresh, and compare the stops.
- [ ] **The calendar file works everywhere.** Test: import it into Google Calendar, Apple Calendar and Outlook, and check that the times are right.
- [ ] **A small local model sticks to the source text.** Test: generate 10 cities and read them all. Allow at most 1 wrong fact across the 10.
- [ ] **Search engines index the trip pages.** Test: submit the sitemap in Google Search Console (free) and check back after 2–4 weeks.

### v2 not doing

- **Writing intros on demand for any city.** That needs AI on the server, which means a paid machine.
- **AI text for each trip.** Intros are per city, not per plan.
- **Accounts or saved trips.** The link is the save button.
- **Short links (`/t/abc123`).** They'd need a database row for every trip, and readable URLs are enough.
- **Search engines indexing dated URLs.** Only the version without a date is canonical, which avoids thousands of near-duplicate pages.

### v2 open questions

- [ ] City in the URL: `lisbon` or `lisbon-pt`? There's more than one Lisbon.
- [ ] Default event times: are 10:00, 12:00, 15:00 and 17:00 sensible, or should days start later?
- [ ] Which cities go in the first list of 100?
- [ ] Will you read every AI intro before it ships, or spot-check some?
