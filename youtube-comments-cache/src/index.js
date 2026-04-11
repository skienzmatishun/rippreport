addEventListener('fetch', event => {
  event.respondWith(handleRequest(event.request))
})

async function handleRequest(request) {
  const url = new URL(request.url)
  const videoId = url.searchParams.get('videoId')
  const maxResults = url.searchParams.get('maxResults') || '20'
  const order = url.searchParams.get('order') || 'relevance'

  if (!videoId) {
    return new Response(JSON.stringify({ error: 'videoId parameter is required' }), { 
      status: 400,
      headers: { 'Content-Type': 'application/json' }
    })
  }

  // Create cache key
  const kvKey = `${videoId}_${maxResults}_${order}`
  
  // Try to get from KV cache
  let data = await youtube_comments.get(kvKey, { type: 'json' })

  if (!data) {
    // Fetch from YouTube API
    const apiKey = YOUTUBE_API_KEY
    
    if (!apiKey) {
      return new Response(JSON.stringify({ error: 'YouTube API key not configured' }), {
        status: 500,
        headers: { 'Content-Type': 'application/json' }
      })
    }

    const apiUrl = `https://www.googleapis.com/youtube/v3/commentThreads?part=snippet&videoId=${videoId}&maxResults=${maxResults}&order=${order}&key=${apiKey}`
    
    try {
      const response = await fetch(apiUrl)
      
      if (!response.ok) {
        const errorData = await response.json()
        return new Response(JSON.stringify({ 
          error: 'YouTube API error',
          details: errorData 
        }), {
          status: response.status,
          headers: { 'Content-Type': 'application/json' }
        })
      }
      
      data = await response.json()

      // Store in KV cache with 24 hour expiration
      await youtube_comments.put(kvKey, JSON.stringify(data), {
        expirationTtl: 86400 // 24 hours
      })
    } catch (error) {
      return new Response(JSON.stringify({ 
        error: 'Failed to fetch from YouTube API',
        message: error.message 
      }), {
        status: 500,
        headers: { 'Content-Type': 'application/json' }
      })
    }
  }

  return new Response(JSON.stringify(data), {
    headers: { 
      'Content-Type': 'application/json',
      'Access-Control-Allow-Origin': '*',
      'Cache-Control': 'public, max-age=3600'
    }
  })
}
