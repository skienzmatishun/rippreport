// Comment validation and spam protection utilities

/**
 * Validate comment content and check for spam
 * @param {Object} commentData - Comment data to validate
 * @param {Object} env - Environment variables
 * @returns {Object} Validation result with isValid and errors
 */
export async function validateComment(commentData, env) {
  const { content, name, userFingerprint } = commentData;
  const errors = [];

  // Content length validation (3-2000 characters)
  if (!content || typeof content !== 'string') {
    errors.push('Comment content is required');
  } else if (content.trim().length < 3) {
    errors.push('Comment must be at least 3 characters long');
  } else if (content.length > 2000) {
    errors.push('Comment must be no more than 2000 characters long');
  }

  // Name validation
  if (name && name.length > 50) {
    errors.push('Name must be no more than 50 characters long');
  }

  // Link detection and limiting (max 2 links per comment)
  if (content) {
    const linkCount = countLinks(content);
    if (linkCount > 2) {
      errors.push('Comments cannot contain more than 2 links');
    }
  }

  // Rate limiting check
  if (userFingerprint) {
    const rateLimitResult = await checkRateLimit(userFingerprint, env);
    if (!rateLimitResult.allowed) {
      errors.push(`Rate limit exceeded. Please wait ${rateLimitResult.waitTime} seconds before commenting again`);
    }
  }

  // Spam keyword detection
  if (content) {
    const spamResult = checkSpamKeywords(content, env);
    if (!spamResult.isValid) {
      errors.push('Comment contains prohibited content');
    }
  }

  return {
    isValid: errors.length === 0,
    errors
  };
}

/**
 * Count the number of links in content
 * @param {string} content - Comment content
 * @returns {number} Number of links found
 */
function countLinks(content) {
  // Match HTTP/HTTPS URLs and www. domains
  const urlRegex = /(https?:\/\/[^\s]+|www\.[^\s]+|[^\s]+\.[a-z]{2,}[^\s]*)/gi;
  const matches = content.match(urlRegex);
  return matches ? matches.length : 0;
}

/**
 * Check rate limiting for a user
 * @param {string} userFingerprint - User's fingerprint
 * @param {Object} env - Environment variables
 * @returns {Object} Rate limit result
 */
async function checkRateLimit(userFingerprint, env) {
  const rateLimitKey = `user:${userFingerprint}:rate`;
  const rateLimitWindow = 60; // 60 seconds
  const maxCommentsPerWindow = 3;
  
  try {
    const rateLimitData = await env.COMMENTS_KV.get(rateLimitKey);
    const now = Date.now();
    
    if (!rateLimitData) {
      // First comment from this user
      const newRateData = {
        count: 1,
        windowStart: now,
        lastComment: now
      };
      await env.COMMENTS_KV.put(rateLimitKey, JSON.stringify(newRateData), {
        expirationTtl: rateLimitWindow * 2 // Keep data for 2 windows
      });
      return { allowed: true };
    }

    const rateData = JSON.parse(rateLimitData);
    const windowElapsed = (now - rateData.windowStart) / 1000;

    if (windowElapsed >= rateLimitWindow) {
      // New window, reset count
      const newRateData = {
        count: 1,
        windowStart: now,
        lastComment: now
      };
      await env.COMMENTS_KV.put(rateLimitKey, JSON.stringify(newRateData), {
        expirationTtl: rateLimitWindow * 2
      });
      return { allowed: true };
    }

    if (rateData.count >= maxCommentsPerWindow) {
      // Rate limit exceeded
      const waitTime = Math.ceil(rateLimitWindow - windowElapsed);
      return { 
        allowed: false, 
        waitTime 
      };
    }

    // Update count
    rateData.count += 1;
    rateData.lastComment = now;
    await env.COMMENTS_KV.put(rateLimitKey, JSON.stringify(rateData), {
      expirationTtl: rateLimitWindow * 2
    });

    return { allowed: true };
  } catch (error) {
    console.error('Rate limit check error:', error);
    // Allow comment on error to avoid blocking legitimate users
    return { allowed: true };
  }
}

/**
 * Check for spam keywords in content
 * @param {string} content - Comment content
 * @param {Object} env - Environment variables
 * @returns {Object} Spam check result
 */
function checkSpamKeywords(content, env) {
  // Default spam keywords - can be extended via environment variables
  const defaultSpamKeywords = [
    'viagra',
    'cialis',
    'casino',
    'poker',
    'lottery',
    'winner',
    'congratulations',
    'click here',
    'free money',
    'make money fast',
    'work from home',
    'lose weight fast'
  ];

  // Get additional keywords from environment if available
  let spamKeywords = defaultSpamKeywords;
  try {
    if (env.SPAM_KEYWORDS) {
      const envKeywords = JSON.parse(env.SPAM_KEYWORDS);
      spamKeywords = [...defaultSpamKeywords, ...envKeywords];
    }
  } catch (error) {
    console.warn('Failed to parse SPAM_KEYWORDS environment variable:', error);
  }

  const lowerContent = content.toLowerCase();
  
  for (const keyword of spamKeywords) {
    if (lowerContent.includes(keyword.toLowerCase())) {
      return { isValid: false, keyword };
    }
  }

  return { isValid: true };
}

/**
 * Sanitize comment content to prevent XSS
 * @param {string} content - Raw comment content
 * @returns {string} Sanitized content
 */
export function sanitizeContent(content) {
  if (!content || typeof content !== 'string') {
    return '';
  }

  // Basic HTML entity encoding to prevent XSS
  return content
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#x27;')
    .replace(/\//g, '&#x2F;');
}

/**
 * Sanitize name field
 * @param {string} name - Raw name
 * @returns {string} Sanitized name
 */
export function sanitizeName(name) {
  if (!name || typeof name !== 'string') {
    return 'Anonymous';
  }

  // Trim and sanitize name
  const sanitized = name.trim().replace(/[<>]/g, '');
  return sanitized || 'Anonymous';
}