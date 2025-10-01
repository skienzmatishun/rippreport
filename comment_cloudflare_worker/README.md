# Cloudflare Comment System

Anonymous comment system built on Cloudflare Workers and KV storage for Hugo sites.

## Project Structure

```
comment_cloudflare_worker/
├── src/
│   ├── index.js              # Main worker entry point
│   ├── handlers/
│   │   └── comments.js       # Comment handling logic
│   └── utils/
│       ├── spam.js           # Spam detection utilities
│       └── cors.js           # CORS handling utilities
├── frontend/
│   ├── comments.js           # Client-side JavaScript
│   └── comments.css          # Comment system styles
├── tests/
│   └── comments.test.js      # Test files
├── wrangler.toml             # Cloudflare Worker configuration
├── package.json              # Node.js dependencies
└── README.md                 # This file
```

## Setup Instructions

1. Install dependencies:
   ```bash
   npm install
   ```

2. Configure KV namespace:
   - Create a KV namespace in Cloudflare dashboard
   - Update the namespace IDs in `wrangler.toml`

3. Update allowed origins:
   - Modify `ALLOWED_ORIGINS` in `wrangler.toml` to match your domain

4. Deploy to Cloudflare:
   ```bash
   npm run deploy:development  # For development
   npm run deploy:production   # For production
   ```

## Development

- `npm run dev` - Start local development server
- `npm run test` - Run tests
- `npm run lint` - Lint code
- `npm run format` - Format code

## Configuration

The system requires:
- Cloudflare KV namespace for comment storage
- Proper CORS configuration for your Hugo site domain
- Rate limiting and spam protection settings

See `wrangler.toml` for configuration options.