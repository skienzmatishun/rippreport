// CORS handling utilities for cross-origin requests

/**
 * Handle CORS preflight requests
 * @param {Request} request - The incoming request
 * @returns {Response} CORS preflight response
 */
export function handleCORS(request) {
  const origin = request.headers.get('Origin');
  const allowedOrigins = [
    'http://localhost:1313',
    'https://rippreport.com',
    'https://www.rippreport.com'
  ];

  const corsHeaders = {
    'Access-Control-Allow-Methods': 'GET, POST, OPTIONS',
    'Access-Control-Allow-Headers': 'Content-Type, Authorization',
    'Access-Control-Max-Age': '86400',
  };

  // Check if origin is allowed
  if (origin && allowedOrigins.includes(origin)) {
    corsHeaders['Access-Control-Allow-Origin'] = origin;
  }

  return new Response(null, {
    status: 204,
    headers: corsHeaders,
  });
}

/**
 * Get CORS headers for regular responses
 * @param {Request} request - The incoming request
 * @returns {Object} CORS headers object
 */
export function getCORSHeaders(request) {
  const origin = request.headers.get('Origin');
  const allowedOrigins = [
    'http://localhost:1313',
    'https://rippreport.com',
    'https://www.rippreport.com'
  ];

  const corsHeaders = {
    'Access-Control-Allow-Methods': 'GET, POST, OPTIONS',
    'Access-Control-Allow-Headers': 'Content-Type, Authorization',
    'Access-Control-Max-Age': '86400',
  };

  // Check if origin is allowed
  if (origin && allowedOrigins.includes(origin)) {
    corsHeaders['Access-Control-Allow-Origin'] = origin;
  }

  return corsHeaders;
}