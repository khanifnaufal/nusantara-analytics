# Nusantara Analytics API Backend

This is the backend service for Nusantara Analytics, built with Go and the Fiber framework. It aggregates data from various external APIs (exchange rates, weather, commodities, market analytics, and earthquakes) with internal caching and route-specific rate limiting.

## Features

- **Route-specific Rate Limiting**:
  - `StrictLimiter` (20 req/min): applied to `/rates`, `/commodities`, and `/market`.
  - `NormalLimiter` (30 req/min): applied to `/weather` and `/quakes`.
  - Exempt: `/health` endpoint is not rate limited.
- **In-Memory Caching**: Cache entries stored globally using `go-cache` to protect external API rate quotas.
- **Rate Limit Headers**: `X-Ratelimit-Limit` and `X-Ratelimit-Remaining` are included on all responses (including HTTP 429).
- **Graceful Shutdown**: Intercepts `SIGINT`/`SIGTERM` to safely drain connection pools within a 5-second window.
- **Production-ready Fallbacks**: Safe defaults for all configuration variables when `.env` is absent.

---

## Local Setup

### 1. Prerequisites
- Go 1.26 or higher installed.

### 2. Configure Environment Variables
Copy `.env.example` to `.env`:
```bash
cp .env.example .env
```
Default fallbacks are defined in code. You may adjust them in `.env`:
- `PORT`: Server port (default: `8080`)
- `ALLOWED_ORIGINS`: CORS origins (default: `*`)
- `CACHE_DEFAULT_TTL`: Default cache expiration in seconds (default: `300`)

### 3. Run the Server
Run the application locally:
```bash
go run main.go
```

The server will start listening on the configured port (e.g. `http://localhost:8080`).

---

## API Endpoints & Response Examples

All endpoints (except health) are prefix-grouped under `/api/v1/`.

### 1. Health Check
* **Endpoint**: `GET /api/v1/health`
* **Rate Limit**: None
* **Response Example**:
```json
{
  "cache_items_count": 0,
  "endpoints": ["rates", "weather", "commodities", "market", "quakes"],
  "status": "ok",
  "uptime": "12s",
  "version": "1.0.0"
}
```

### 2. Exchange Rates
* **Endpoint**: `GET /api/v1/rates`
* **Rate Limit**: Strict Preset (20 req/min)
* **Response Example**:
```json
{
  "base": "IDR",
  "date": "2026-06-08",
  "rates": {
    "AUD": 0.000094,
    "EUR": 0.000057,
    "JPY": 0.0093,
    "MYR": 0.00028,
    "SAR": 0.00023,
    "SGD": 0.000082,
    "USD": 0.000062
  },
  "last_updated": "2026-06-08T22:20:00+07:00"
}
```

### 3. Weather Forecast
* **Endpoint**: `GET /api/v1/weather`
* **Rate Limit**: Normal Preset (30 req/min)
* **Response Example**:
```json
{
  "location": "Jakarta",
  "temperature": 29.5,
  "condition": "Cloudy",
  "humidity": 78,
  "wind_speed": 12.4,
  "last_updated": "2026-06-08T22:20:00+07:00"
}
```

### 4. Commodity Prices
* **Endpoint**: `GET /api/v1/commodities`
* **Rate Limit**: Strict Preset (20 req/min)
* **Response Example**:
```json
{
  "commodities": [
    {"name": "Gold", "price": 2350.40, "unit": "USD/oz"},
    {"name": "Crude Oil (WTI)", "price": 75.30, "unit": "USD/bbl"},
    {"name": "Coal (Newcastle)", "price": 134.15, "unit": "USD/ton"}
  ],
  "last_updated": "2026-06-08T22:20:00+07:00"
}
```

### 5. Market Indices
* **Endpoint**: `GET /api/v1/market`
* **Rate Limit**: Strict Preset (20 req/min)
* **Response Example**:
```json
{
  "indices": [
    {"name": "IHSG (JKSE)", "value": 6890.50, "change": 0.45},
    {"name": "Dow Jones", "value": 38780.20, "change": -0.12},
    {"name": "S&P 500", "value": 5350.80, "change": 0.08}
  ],
  "last_updated": "2026-06-08T22:20:00+07:00"
}
```

### 6. Recent Earthquakes
* **Endpoint**: `GET /api/v1/quakes`
* **Rate Limit**: Normal Preset (30 req/min)
* **Response Example**:
```json
{
  "earthquakes": [
    {
      "magnitude": 5.2,
      "depth": "10 km",
      "location": "92 km BaratDaya KEBUMEN-JATENG",
      "time": "2026-06-08T20:15:30+07:00"
    }
  ],
  "last_updated": "2026-06-08T22:20:00+07:00"
}
```

---

## Rate Limit Response (HTTP 429)

When rate limits are exceeded on any endpoint, the service returns a standard `429 Too Many Requests` status code with the following headers and payload:

**Response Headers**:
```http
HTTP/1.1 429 Too Many Requests
Retry-After: 60
X-Ratelimit-Limit: 20  (or 30)
X-Ratelimit-Remaining: 0
Content-Type: application/json
```

**Response Body**:
```json
{
  "error": "too many requests",
  "retry_after": "60s"
}
```
