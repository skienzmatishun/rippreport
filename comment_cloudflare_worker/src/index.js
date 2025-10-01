// Main Cloudflare Worker entry point for comment system
import { handleComments } from './handlers/comments.js';
import { handleCORS, getCORSHeaders } from './utils/cors.js';

export default {
  async fetch(request, env, ctx) {
    const url = new URL(request.url);
    const { pathname, searchParams } = url;
    const method = request.method;

    // Handle CORS preflight requests
    if (method === 'OPTIONS') {
      return handleCORS(request);
    }

    try {
      // Debug endpoint to inspect KV store
      if (pathname === '/api/debug') {
        const postId = searchParams.get('postId') || 'debug-test-post';
        
        const commentsKey = `post:${postId}:all-comments`;
        const commentsData = await env.COMMENTS_KV.get(commentsKey);
        
        const debugInfo = {
          postId,
          commentsKey,
          rawData: commentsData,
          parsedData: commentsData ? JSON.parse(commentsData) : null
        };
        
        return new Response(JSON.stringify(debugInfo, null, 2), {
          status: 200,
          headers: {
            'Content-Type': 'application/json',
            ...getCORSHeaders(request),
          },
        });
      }
      
      // Route matching for comment API endpoints
      if (pathname.startsWith('/api/comments/')) {
        // Extract post ID from path: /api/comments/{postId}
        const pathParts = pathname.split('/');
        if (pathParts.length >= 4) {
          const postId = pathParts.slice(3).join('/'); // Support nested paths in post IDs
          return await handleComments(request, env, postId);
        }
      }

      // Default 404 response for unmatched routes
      return new Response('Not Found', {
        status: 404,
        headers: {
          'Content-Type': 'text/plain',
          ...getCORSHeaders(request),
        },
      });

    } catch (error) {
      console.error('Worker error:', error);
      return new Response('Internal Server Error', {
        status: 500,
        headers: {
          'Content-Type': 'text/plain',
          ...getCORSHeaders(request),
        },
      });
    }
  },
};

