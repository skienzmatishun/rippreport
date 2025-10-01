// Comment handling functions for storage and retrieval
import { getCORSHeaders } from '../utils/cors.js';
import { validateComment, sanitizeContent, sanitizeName } from '../utils/validation.js';

/**
 * Handle comment-related requests
 * @param {Request} request - The incoming request
 * @param {Object} env - Environment variables and bindings
 * @param {string} postId - The post identifier
 * @returns {Response} API response
 */
export async function handleComments(request, env, postId) {
  const method = request.method;
  const corsHeaders = getCORSHeaders(request);

  try {
    switch (method) {
      case 'GET':
        return await getComments(env, postId, corsHeaders);
      case 'POST':
        return await createComment(request, env, postId, corsHeaders);
      default:
        return new Response('Method Not Allowed', {
          status: 405,
          headers: {
            'Content-Type': 'text/plain',
            ...corsHeaders,
          },
        });
    }
  } catch (error) {
    console.error('Comment handler error:', error);
    return new Response('Internal Server Error', {
      status: 500,
      headers: {
        'Content-Type': 'application/json',
        ...corsHeaders,
      },
      body: JSON.stringify({ error: 'Internal server error' }),
    });
  }
}

/**
 * Retrieve comments for a specific post
 * @param {Object} env - Environment variables and bindings
 * @param {string} postId - The post identifier
 * @param {Object} corsHeaders - CORS headers
 * @returns {Response} Comments response
 */
async function getComments(env, postId, corsHeaders) {
  try {
    // Get all comments for this post stored as a single object
    const commentsKey = `post:${postId}:all-comments`;
    console.log('GET - commentsKey:', commentsKey);
    const commentsData = await env.COMMENTS_KV.get(commentsKey);
    console.log('GET - commentsData:', commentsData);
    
    if (!commentsData) {
      // No comments for this post yet
      console.log('GET - No comments found');
      return new Response(JSON.stringify([]), {
        status: 200,
        headers: {
          'Content-Type': 'application/json',
          ...corsHeaders,
        },
      });
    }

    const comments = JSON.parse(commentsData);
    console.log('GET - parsed comments:', comments);

    // Build nested comment structure
    const nestedComments = buildCommentTree(comments);
    console.log('GET - nested comments:', nestedComments);

    return new Response(JSON.stringify(nestedComments), {
      status: 200,
      headers: {
        'Content-Type': 'application/json',
        ...corsHeaders,
      },
    });
  } catch (error) {
    console.error('Error fetching comments:', error);
    return new Response(JSON.stringify({ error: 'Failed to fetch comments' }), {
      status: 500,
      headers: {
        'Content-Type': 'application/json',
        ...corsHeaders,
      },
    });
  }
}

/**
 * Create a new comment
 * @param {Request} request - The incoming request
 * @param {Object} env - Environment variables and bindings
 * @param {string} postId - The post identifier
 * @param {Object} corsHeaders - CORS headers
 * @returns {Response} Created comment response
 */
async function createComment(request, env, postId, corsHeaders) {
  try {
    const body = await request.json();
    const { content, name = 'Anonymous', parentId = null, userFingerprint } = body;

    // Validate comment data
    const validation = await validateComment({ content, name, userFingerprint }, env);
    if (!validation.isValid) {
      return new Response(JSON.stringify({ 
        error: 'Validation failed', 
        details: validation.errors 
      }), {
        status: 400,
        headers: {
          'Content-Type': 'application/json',
          ...corsHeaders,
        },
      });
    }

    // Sanitize input data
    const sanitizedContent = sanitizeContent(content);
    const sanitizedName = sanitizeName(name);

    // Generate unique comment ID
    const commentId = crypto.randomUUID();
    const timestamp = new Date().toISOString();

    // Calculate nesting level
    let level = 0;
    if (parentId) {
      // Get existing comments to find the parent
      const commentsKey = `post:${postId}:all-comments`;
      const existingCommentsData = await env.COMMENTS_KV.get(commentsKey);
      if (existingCommentsData) {
        const existingComments = JSON.parse(existingCommentsData);
        const parent = existingComments.find(c => c.id === parentId);
        if (parent) {
          level = Math.min(parent.level + 1, 2); // Flatten after level 2
        }
      }
    }

    // Create comment object
    const comment = {
      id: commentId,
      postId,
      content: sanitizedContent,
      name: sanitizedName,
      userFingerprint,
      parentId,
      timestamp,
      level,
      children: []
    };

    // Store comment in the all-comments object for this post
    const commentsKey = `post:${postId}:all-comments`;
    console.log('POST - commentsKey:', commentsKey);
    
    try {
      const existingCommentsData = await env.COMMENTS_KV.get(commentsKey);
      console.log('POST - existingCommentsData:', existingCommentsData);
      const existingComments = existingCommentsData ? JSON.parse(existingCommentsData) : [];
      console.log('POST - existingComments before:', existingComments);
      
      // Add the new comment to the list
      existingComments.push(comment);
      console.log('POST - existingComments after:', existingComments);
      
      const putResult = await env.COMMENTS_KV.put(commentsKey, JSON.stringify(existingComments));
      console.log('POST - KV put result:', putResult);
      console.log('POST - stored to KV successfully');
    } catch (kvError) {
      console.error('POST - KV operation failed:', kvError);
      throw kvError;
    }

    // Update comment count
    const countKey = `post:${postId}:count`;
    const currentCount = await env.COMMENTS_KV.get(countKey);
    const newCount = currentCount ? parseInt(currentCount) + 1 : 1;
    await env.COMMENTS_KV.put(countKey, newCount.toString());

    return new Response(JSON.stringify(comment), {
      status: 201,
      headers: {
        'Content-Type': 'application/json',
        ...corsHeaders,
      },
    });
  } catch (error) {
    console.error('Error creating comment:', error);
    return new Response(JSON.stringify({ error: 'Failed to create comment' }), {
      status: 500,
      headers: {
        'Content-Type': 'application/json',
        ...corsHeaders,
      },
    });
  }
}

/**
 * Build nested comment tree structure
 * @param {Array} comments - Flat array of comments
 * @returns {Array} Nested comment structure
 */
function buildCommentTree(comments) {
  if (!Array.isArray(comments)) {
    return [];
  }

  const commentMap = new Map();

  // Create a new node for each comment and store it in the map.
  // This avoids mutating the original input array.
  comments.forEach(comment => {
    commentMap.set(comment.id, { ...comment, children: [] });
  });

  const rootComments = [];
  
  // Build the tree structure using the nodes from our map.
  commentMap.forEach(commentNode => {
    if (commentNode.parentId && commentMap.has(commentNode.parentId)) {
      const parentNode = commentMap.get(commentNode.parentId);
      parentNode.children.push(commentNode);
    } else {
      rootComments.push(commentNode);
    }
  });

  // Sort comments by timestamp (oldest first)
  const sortByTimestamp = (a, b) => new Date(a.timestamp) - new Date(b.timestamp);
  
  rootComments.sort(sortByTimestamp);
  
  // Recursively sort children
  function sortChildren(comment) {
    comment.children.sort(sortByTimestamp);
    comment.children.forEach(sortChildren);
  }
  
  rootComments.forEach(sortChildren);

  return rootComments;
}