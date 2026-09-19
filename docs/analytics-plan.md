# Plan: see how many people used TravelTab last week

_Written 2026-09-19. Prices were checked on that date._

## First: Plausible has no free tier

- **Plausible Cloud** gives you a 30-day free trial with no card needed. After that it's **$9/month** for up to 10k pageviews.
- **Self-hosted Plausible** (Community Edition) is free software, but it needs its own server with **2 GB+ RAM** to run its databases. On Fly that costs money, and you'd have to maintain it.

## You may not need Plausible

TravelTab already has visit tracking. `/stats` returns unique visitors, page views, searches, a day-by-day breakdown and the top cities, and it skips bots. It stays hidden (404) until you set a token.

Either way, counting only starts once the code is live. **Nothing counts visits from before that**, so "last week" means a week after you deploy.

## Option A: use your own `/stats` (recommended: free, about 2 minutes)

- [ ] **1. Make a token and save it somewhere.** Fly won't show it to you again.
  ```bash
  openssl rand -hex 16
  ```
- [ ] **2. Give it to Fly.** Setting a secret restarts the app.
  ```bash
  fly secrets set STATS_TOKEN=<your-token>
  ```
- [ ] **3. Deploy**, if the analytics commit isn't live yet.
  ```bash
  fly deploy
  ```
- [ ] **4. Check it works:** visit the site once, then open:
  ```
  https://traveltabbypatronfragas.com/stats?token=<your-token>&days=7
  ```
  You should see `page_views: 1` or more.
- [ ] **5. Wait 7 days, open the same URL, and read `unique_visitors`.** ✅

## Option B: Plausible Cloud trial (about 15 minutes, free for 30 days)

- [ ] **1.** Sign up at plausible.io and add the site `traveltabbypatronfragas.com`.
- [ ] **2.** Copy the script snippet Plausible shows you.
- [ ] **3.** Paste it into the `<head>` of `views/index.go.tpl`. That file is the only full page. Search results are swapped in by HTMX, so they won't be counted twice.
- [ ] **4.** Run `fly deploy`.
- [ ] **5. Check it works:** open the site in a private window with no ad blocker. You should show up in Plausible's live view.
- [ ] **6.** Wait 7 days, set the dashboard to **"Last 7 days"** and read **Unique visitors**. ✅
- [ ] **7. Set a reminder for day 28:** either pay, or remove the script.

### Optional: count searches

Searches don't change the URL, so Plausible only sees the first page load. To count them, add this script to `views/index.go.tpl`, then create a goal called `Search` in Plausible under **Settings → Goals**:

```html
<script>
  document.addEventListener("htmx:afterRequest", function (e) {
    if (e.detail.successful && e.detail.pathInfo.requestPath.startsWith("/process-form/") && window.plausible) {
      plausible("Search", { props: { city: new URLSearchParams(e.detail.pathInfo.requestPath.split("?")[1]).get("city_name") } });
    }
  });
</script>
```

`/stats` already counts searches, so Option A doesn't need this.

## Recommendation

**Option A.** It's already built, it's free, it counts searches and cities, and your visit data stays on your own server. Choose Plausible only if you want a dashboard to look at.

## Sources

- [Plausible pricing](https://plausible.io/#pricing)
- [Plausible custom events](https://plausible.io/docs/custom-event-goals)
- [Plausible self-hosting](https://plausible.io/docs/self-hosting)
- [Plausible Community Edition requirements](https://github.com/plausible/community-edition)
