# RippDirect Worker

A Cloudflare Worker that provides URL preview functionality for social media previews and link processing.

## Features

- **URL Preview**: Fetches webpage metadata including title, description, and images
- **Open Graph Support**: Extracts Open Graph meta tags for rich previews
- **Error Handling**: Graceful error handling with informative responses
- **CORS Support**: Allows cross-origin requests from your Hugo site
- **Caching**: Caches responses for better performance

## API Usage

The worker accepts GET requests with a `q` parameter containing the URL to preview:

```
https://links.rippreport.com/?q=https://example.com
```

### Response Format

```json
{
  "title": "Page Title",
  "description": "Page description",
  "image": "https://example.com/image.jpg",
  "url": "https://example.com",
  "charset": "utf-8"
}
```

### Error Response

```json
{
  "error": "[Error] https://example.com returned 404",
  "url": "https://example.com"
}
```

## Setup

1. Install dependencies:
   ```bash
   npm install
   ```

2. Configure your Cloudflare account:
   ```bash
   npx wrangler login
   ```

3. Deploy to development:
   ```bash
   npm run dev
   ```

4. Deploy to production:
   ```bash
   npm run deploy:production
   ```

## Configuration

Update `wrangler.toml` to configure:
- Custom domains and routes
- Environment variables
- Worker settings

## Integration

This worker is used by the Hugo socialpreview shortcode to fetch webpage metadata for creating rich link previews.

## Security

- Only allows HTTP/HTTPS protocols
- Validates all input URLs
- Includes timeout protection
- CORS headers for secure cross-origin requests