# URL Shortener Service

This is a URL shortener service that allows you to convert long URLs into shorter, more manageable links. It provides a simple API endpoint that accepts POST requests to generate short URLs for the provided long URLs. The service also includes rate limiting based on IP address to prevent abuse.

## API Documentation

### Base URL

The base URL for accessing the API is: `http://your-domain.com/api/v1`

### Endpoint

#### POST `http://your-domain.com/api/v1`

This endpoint generates a short URL for a given long URL. Provide Body as sample Request payload given

##### Request Payload

The request payload should be a JSON object with the following fields:

- `URL`: The actual long URL that you want to shorten.
- `CustomShort` (optional): If you want to customize the short URL, you can provide a custom shortcode. If the custom shortcode is available, it will be used as the short URL. If it's already taken, a random short URL will be generated instead.
- `Expiry` (optional): The expiration time for the short URL in hours. If not provided, the default expiration time is set to 24 hours.

Example Request Payload:

```json
{
  "URL": "https://www.example.com/very/long/url",
  "CustomShort": "customshortcode",
  "Expiry": 48
}
```

##### Response

If the request is successful, the API will respond with a JSON object containing the generated short URL.

Example Response:

```json
{
    "url": "https://example.com/long-url",
    "shorturl": "https://your-domain.com/short-url",
    "expiry": "24h",
    "rate_limit": 8,
    "rate_limit_reset": "30m"
}
```
- `URL`: The original long URL provided by the user.
- `ShortURL`: The generated or custom short URL for the long URL.
- `Expiry`: The duration after which the short URL will expire (if specified).
- `XRateRemaining`: The number of remaining allowed requests within the rate limit.
- `XRateLimitReset`: The duration until the rate limit resets and allows new requests.

If there is an error during the request, the API will respond with an appropriate error message and status code.

## Rate Limiting

The URL shortener service implements rate limiting based on IP address to prevent abuse and ensure fair usage. Each IP address is limited to a maximum of 10 POST requests per 30 minutes. If the rate limit is exceeded, the API will respond with a `429 Too Many Requests` status code along with an error message indicating that the rate limit has been exceeded.

## Getting Started

To set up and use the URL shortener service, follow these steps:

1. Clone the repository or download the project files.
2. Navigate to the project directory.
3. Run the following command to start the services:

   ```
   docker-compose up -d
   ```

   This will start the URL shortener service and any required dependencies defined in the `docker-compose.yml` file.

4. Once the services are up and running, you can start using the API endpoint to generate short URLs by sending POST requests with the required payload.
