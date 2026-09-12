addEventListener('fetch', event => {
  event.respondWith(handleRequest(event.request))
})

// ---------------------------------------------------------------------------
// YouTube video ID extraction
// ---------------------------------------------------------------------------
// Matches:
//   https://www.youtube.com/watch?v=VIDEO_ID
//   https://youtube.com/watch?v=VIDEO_ID
//   https://m.youtube.com/watch?v=VIDEO_ID
//   https://youtu.be/VIDEO_ID
//   https://www.youtube.com/shorts/VIDEO_ID
//   https://www.youtube.com/embed/VIDEO_ID
// Also tolerates extra query params (e.g. &t=30s, &list=...).
function extractYouTubeVideoId(url) {
  if (!url) return null

  // youtu.be/VIDEO_ID  (short link)
  let match = url.match(
    /(?:youtu\.be\/)([A-Za-z0-9_-]{11})(?:[?&#]|$)/
  )
  if (match) return match[1]

  // youtube.com/watch?v=VIDEO_ID
  match = url.match(
    /(?:youtube\.com\/watch\?)(?:[^#]*&)?v=([A-Za-z0-9_-]{11})(?:[&#]|$)/
  )
  if (match) return match[1]

  // youtube.com/shorts/VIDEO_ID  or  youtube.com/embed/VIDEO_ID
  match = url.match(
    /(?:youtube\.com\/(?:shorts|embed|v)\/)([A-Za-z0-9_-]{11})(?:[?&#]|$)/
  )
  if (match) return match[1]

  return null
}

// Build a thumbnail URL directly from the video ID (YouTube's CDN).
function youtubeThumbnailUrl(videoId) {
  return `https://img.youtube.com/vi/${videoId}/maxresdefault.jpg`
}

// Returns true only if the comment text is exactly one URL (nothing else).
function isUrlOnly(text) {
  if (!text) return false
  return /^https?:\/\/[^\s]+$/i.test(text.trim())
}

// ---------------------------------------------------------------------------
// Attach a linkPreview object ONLY for comments that are just a YouTube URL.
// Every other comment is left completely untouched.
// ---------------------------------------------------------------------------
function attachYouTubePreviews(data) {
  if (!data || !Array.isArray(data.items)) return data

  const processSnippet = (snippet) => {
    if (!snippet || !isUrlOnly(snippet.textOriginal)) return

    const videoId = extractYouTubeVideoId(snippet.textOriginal.trim())
    if (!videoId) return // not a YouTube link — leave it alone

    snippet.linkPreview = {
      type: 'youtube',
      videoId,
      thumbnail: youtubeThumbnailUrl(videoId),
      fallbackThumbnail: `https://img.youtube.com/vi/${videoId}/hqdefault.jpg`,
      url: snippet.textOriginal.trim()
    }
  }

  for (const item of data.items) {
    // Top-level comment
    processSnippet(item?.snippet?.topLevelComment?.snippet)

    // Replies (defensive — commentThreads doesn't return these by default)
    const replies = item?.replies?.comments
    if (Array.isArray(replies)) {
      for (const reply of replies) {
        processSnippet(reply?.snippet)
      }
    }
  }

  return data
}

// ---------------------------------------------------------------------------
// Main request handler
// ---------------------------------------------------------------------------
async function handleRequest(request) {
  const url = new URL(request.url)
  const videoId = url.searchParams.get('videoId')
  const maxResults = url.searchParams.get('maxResults') || '20'
  const order = url.searchParams.get('order') || 'relevance'

  // Handle CORS preflight
  if (request.method === 'OPTIONS') {
    return new Response(null, {
      status: 204,
      headers: {
        'Access-Control-Allow-Origin': '*',
        'Access-Control-Allow-Methods': 'GET, OPTIONS',
        'Access-Control-Allow-Headers': 'Content-Type',
        'Access-Control-Max-Age': '86400'
      }
    })
  }

  if (!videoId) {
    return new Response(
      JSON.stringify({ error: 'videoId parameter is required' }),
      {
        status: 400,
        headers: {
          'Content-Type': 'application/json',
          'Access-Control-Allow-Origin': '*'
        }
      }
    )
  }

  // Cache key (bumped to _v3 to invalidate any previously cached data)
  const kvKey = `${videoId}_${maxResults}_${order}_v3`

  // Try KV cache first
  let data = await youtube_comments.get(kvKey, { type: 'json' })

  if (data) {
    return new Response(JSON.stringify(data), {
      headers: {
        'Content-Type': 'application/json',
        'Access-Control-Allow-Origin': '*',
        'Cache-Control': 'public, max-age=86400',
        'X-Cache': 'HIT'
      }
    })
  }

  const apiKey = YOUTUBE_API_KEY

  if (!apiKey) {
    return new Response(
      JSON.stringify({ error: 'YouTube API key not configured' }),
      {
        status: 500,
        headers: {
          'Content-Type': 'application/json',
          'Access-Control-Allow-Origin': '*'
        }
      }
    )
  }

  const apiUrl = `https://www.googleapis.com/youtube/v3/commentThreads?part=snippet&videoId=${videoId}&maxResults=${maxResults}&order=${order}&key=${apiKey}`

  try {
    const response = await fetch(apiUrl)

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}))

      const errorResponse = {
        error: 'YouTube API error',
        details: errorData
      }

      await youtube_comments.put(kvKey, JSON.stringify(errorResponse), {
        expirationTtl: 3600
      })

      return new Response(JSON.stringify(errorResponse), {
        status: response.status,
        headers: {
          'Content-Type': 'application/json',
          'Access-Control-Allow-Origin': '*',
          'X-Cache': 'MISS'
        }
      })
    }

    data = await response.json()

    // Attach YouTube previews ONLY where the comment is just a YouTube URL.
    // Every other comment is untouched.
    data = attachYouTubePreviews(data)

    await youtube_comments.put(kvKey, JSON.stringify(data))

    return new Response(JSON.stringify(data), {
      headers: {
        'Content-Type': 'application/json',
        'Access-Control-Allow-Origin': '*',
        'Cache-Control': 'public, max-age=86400',
        'X-Cache': 'MISS'
      }
    })
  } catch (error) {
    return new Response(
      JSON.stringify({
        error: 'Failed to fetch from YouTube API',
        message: error.message
      }),
      {
        status: 500,
        headers: {
          'Content-Type': 'application/json',
          'Access-Control-Allow-Origin': '*'
        }
      }
    )
  }
}