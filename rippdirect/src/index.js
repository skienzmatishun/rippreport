/**
 * RippDirect Worker - URL Preview and Link Processing
 * Fetches webpage metadata for social previews
 */

export default {
  async fetch(request, env, ctx) {
    const url = new URL(request.url);
    
    // Handle CORS preflight requests
    if (request.method === 'OPTIONS') {
      return new Response(null, {
        headers: {
          'Access-Control-Allow-Origin': '*',
          'Access-Control-Allow-Methods': 'GET, OPTIONS',
          'Access-Control-Allow-Headers': 'Content-Type',
        },
      });
    }

    // Only allow GET requests
    if (request.method !== 'GET') {
      return new Response('Method not allowed', { status: 405 });
    }

    try {
      // Get the URL to preview from query parameter
      const targetUrl = url.searchParams.get('q');
      
      if (!targetUrl) {
        return new Response('Missing URL parameter "q"', { 
          status: 400,
          headers: { 'Content-Type': 'text/plain' }
        });
      }

      // Validate URL
      let parsedUrl;
      try {
        parsedUrl = new URL(targetUrl);
      } catch (e) {
        return new Response('Invalid URL', { 
          status: 400,
          headers: { 'Content-Type': 'text/plain' }
        });
      }

      // Security: Only allow HTTP/HTTPS
      if (!['http:', 'https:'].includes(parsedUrl.protocol)) {
        return new Response('Invalid protocol', { 
          status: 400,
          headers: { 'Content-Type': 'text/plain' }
        });
      }

      // Fetch the webpage
      const response = await fetch(targetUrl, {
        headers: {
          'User-Agent': 'Mozilla/5.0 (compatible; RippReport-LinkBot/1.0)',
          'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8',
        },
        // Add timeout and size limits
        cf: {
          timeout: 10000, // 10 second timeout
        }
      });

      if (!response.ok) {
        return new Response(JSON.stringify({
          error: `[Error] ${targetUrl} returned ${response.status}`,
          url: targetUrl
        }), {
          status: 200,
          headers: {
            'Content-Type': 'application/json',
            'Access-Control-Allow-Origin': '*',
          }
        });
      }

      // Get the HTML content
      const html = await response.text();
      
      // Extract metadata
      const metadata = extractMetadata(html, targetUrl);
      
      return new Response(JSON.stringify(metadata), {
        headers: {
          'Content-Type': 'application/json',
          'Cache-Control': 'public, max-age=3600', // Cache for 1 hour
          'Access-Control-Allow-Origin': '*',
        },
      });

    } catch (error) {
      console.error('Error processing request:', error);
      return new Response(JSON.stringify({
        error: `[Error] Failed to process ${url.searchParams.get('q')}`,
        url: url.searchParams.get('q')
      }), {
        status: 200,
        headers: {
          'Content-Type': 'application/json',
          'Access-Control-Allow-Origin': '*',
        }
      });
    }
  },
};

/**
 * Extract metadata from HTML content
 */
function extractMetadata(html, url) {
  // Basic regex patterns for extracting metadata
  const titleMatch = html.match(/<title[^>]*>([^<]+)<\/title>/i);
  const ogTitleMatch = html.match(/<meta[^>]*property=["\']og:title["\'][^>]*content=["\']([^"']+)["\'][^>]*>/i);
  const ogDescriptionMatch = html.match(/<meta[^>]*property=["\']og:description["\'][^>]*content=["\']([^"']+)["\'][^>]*>/i);
  const ogImageMatch = html.match(/<meta[^>]*property=["\']og:image["\'][^>]*content=["\']([^"']+)["\'][^>]*>/i);
  const descriptionMatch = html.match(/<meta[^>]*name=["\']description["\'][^>]*content=["\']([^"']+)["\'][^>]*>/i);
  const charsetMatch = html.match(/<meta[^>]*charset=["\']?([^"'\s>]+)["\']?[^>]*>/i) || 
                      html.match(/charset=([^"'\s>]+)/i);

  // Extract title (prefer og:title, fallback to title tag)
  let title = ogTitleMatch ? ogTitleMatch[1] : (titleMatch ? titleMatch[1] : 'Title not available');
  title = decodeHtmlEntities(title.trim());

  // Extract description (prefer og:description, fallback to meta description)
  let description = ogDescriptionMatch ? ogDescriptionMatch[1] : (descriptionMatch ? descriptionMatch[1] : 'Description not available');
  description = decodeHtmlEntities(description.trim());

  // Extract image (og:image)
  let image = ogImageMatch ? ogImageMatch[1] : null;
  if (image) {
    // Convert relative URLs to absolute
    if (image.startsWith('/')) {
      const baseUrl = new URL(url);
      image = `${baseUrl.protocol}//${baseUrl.host}${image}`;
    } else if (!image.startsWith('http')) {
      const baseUrl = new URL(url);
      image = `${baseUrl.protocol}//${baseUrl.host}/${image}`;
    }
  }

  // Extract charset
  const charset = charsetMatch ? charsetMatch[1].toLowerCase() : 'utf-8';

  return {
    title,
    description,
    image,
    url,
    charset
  };
}

/**
 * Decode HTML entities
 */
function decodeHtmlEntities(text) {
  const entities = {
    '&amp;': '&',
    '&lt;': '<',
    '&gt;': '>',
    '&quot;': '"',
    '&#39;': "'",
    '&apos;': "'",
    '&nbsp;': ' ',
    '&hellip;': '...',
    '&mdash;': '—',
    '&ndash;': '–',
    '&rsquo;': "'",
    '&lsquo;': "'",
    '&ldquo;': '"',
    '&rdquo;': '"'
  };
  
  return text.replace(/&[#\w]+;/g, (entity) => {
    return entities[entity] || entity;
  });
}