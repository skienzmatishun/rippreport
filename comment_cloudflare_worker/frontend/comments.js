// Frontend JavaScript for comment system
// To be implemented in later tasks

class CommentSystem {
  constructor(postId, apiEndpoint) {
    this.postId = postId;
    this.apiEndpoint = apiEndpoint;
    
    // Initialize user identity system
    this.userIdentity = this.getUserIdentity();
    this.userColor = this.assignUserColor(this.userIdentity.fingerprint);
  }

  /**
   * Load and display comments for the current post
   * Requirements: 1.5, 2.3, 2.4, 9.1
   */
  async loadComments() {
    try {
      const response = await fetch(`${this.apiEndpoint}/api/comments/${encodeURIComponent(this.postId)}`, {
        method: 'GET',
        headers: {
          'Accept': 'application/json',
        },
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const comments = await response.json();
      this.displayComments(comments);
      return comments;
    } catch (error) {
      console.error('Error loading comments:', error);
      this.displayError('Failed to load comments. Please try again later.');
      return [];
    }
  }

  /**
   * Display comments in the DOM with nested structure
   * Requirements: 2.3, 2.4, 9.1
   */
  displayComments(comments) {
    const container = document.getElementById('comments-list');
    if (!container) {
      console.error('Comments container not found');
      return;
    }

    // Clear existing comments
    container.innerHTML = '';

    // Update comment count
    const totalCount = this.countComments(comments);
    this.updateCommentsCount(totalCount);

    if (comments.length === 0) {
      this.showEmptyState();
      return;
    }

    // Render each top-level comment
    comments.forEach(comment => {
      const commentElement = this.renderComment(comment, 0);
      this.addCommentAnimations(commentElement);
      container.appendChild(commentElement);
    });
  }

  /**
   * Render a single comment with its children
   * Requirements: 2.3, 2.4, 3.3
   */
  renderComment(comment, level) {
    const commentDiv = document.createElement('div');
    commentDiv.className = `comment comment-level-${level}`;
    commentDiv.setAttribute('data-comment-id', comment.id);
    commentDiv.setAttribute('data-level', level);

    // Apply user color styling
    if (comment.userFingerprint) {
      this.applyUserColorStyling(commentDiv, comment.userFingerprint);
    }

    // Format timestamp
    const timestamp = new Date(comment.timestamp);
    const timeString = this.formatTimestamp(timestamp);

    // Create comment HTML structure
    commentDiv.innerHTML = `
      <div class="comment-header">
        <span class="comment-author">${this.escapeHtml(comment.name)}</span>
        <span class="comment-timestamp" title="${timestamp.toLocaleString()}">${timeString}</span>
      </div>
      <div class="comment-content">
        ${this.formatCommentContent(comment.content)}
      </div>
      <div class="comment-actions">
        <button class="reply-button" onclick="commentSystem.showReplyForm('${comment.id}', ${level})">
          Reply
        </button>
      </div>
      <div class="reply-form-container" id="reply-form-${comment.id}" style="display: none;"></div>
      <div class="comment-children"></div>
    `;

    // Render children comments
    if (comment.children && comment.children.length > 0) {
      const childrenContainer = commentDiv.querySelector('.comment-children');
      comment.children.forEach(child => {
        // Implement level flattening for deep threads (3+ levels)
        const childLevel = level >= 2 ? 2 : level + 1;
        const childElement = this.renderComment(child, childLevel);
        childrenContainer.appendChild(childElement);
      });
    }

    return commentDiv;
  }

  /**
   * Format comment content with basic HTML processing
   * Requirements: 1.5, 2.3
   */
  formatCommentContent(content) {
    // Escape HTML to prevent XSS
    let formatted = this.escapeHtml(content);
    
    // Convert line breaks to <br> tags
    formatted = formatted.replace(/\n/g, '<br>');
    
    // Convert URLs to clickable links (basic implementation)
    const urlRegex = /(https?:\/\/[^\s]+)/g;
    formatted = formatted.replace(urlRegex, '<a href="$1" target="_blank" rel="noopener noreferrer">$1</a>');
    
    return formatted;
  }

  /**
   * Escape HTML to prevent XSS attacks
   */
  escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
  }

  /**
   * Format timestamp for display
   */
  formatTimestamp(date) {
    const now = new Date();
    const diffMs = now - date;
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMs / 3600000);
    const diffDays = Math.floor(diffMs / 86400000);

    if (diffMins < 1) {
      return 'just now';
    } else if (diffMins < 60) {
      return `${diffMins} minute${diffMins === 1 ? '' : 's'} ago`;
    } else if (diffHours < 24) {
      return `${diffHours} hour${diffHours === 1 ? '' : 's'} ago`;
    } else if (diffDays < 7) {
      return `${diffDays} day${diffDays === 1 ? '' : 's'} ago`;
    } else {
      return date.toLocaleDateString();
    }
  }

  /**
   * Display error message to user
   */
  displayError(message) {
    this.showErrorState(message, true);
  }

  /**
   * Display loading state
   */
  displayLoading() {
    this.showLoadingState('Loading comments...');
  }

  /**
   * Refresh comments after a new comment is posted
   * Requirements: 1.5, 2.4
   */
  async refreshComments() {
    this.displayLoading();
    await this.loadComments();
  }

  /**
   * Submit a new comment
   * Requirements: 1.1, 1.2, 1.3, 2.1
   */
  async submitComment(content, name = 'Anonymous', parentId = null) {
    try {
      // Validate input
      const validation = this.validateCommentInput(content, name);
      if (!validation.isValid) {
        this.displayValidationErrors(validation.errors);
        return null;
      }

      // Show loading state
      this.setSubmissionLoading(true, parentId);

      // Prepare comment data
      const commentData = {
        content: content.trim(),
        name: name.trim() || 'Anonymous',
        parentId: parentId,
        userFingerprint: this.userIdentity.fingerprint
      };

      // Submit to API
      const response = await fetch(`${this.apiEndpoint}/api/comments/${encodeURIComponent(this.postId)}`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
        },
        body: JSON.stringify(commentData)
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.error || `HTTP error! status: ${response.status}`);
      }

      const newComment = await response.json();
      
      // Clear form and hide loading state
      this.clearCommentForm(parentId);
      this.setSubmissionLoading(false, parentId);
      
      // Refresh comments to show the new one
      await this.refreshComments();
      
      return newComment;
    } catch (error) {
      console.error('Error submitting comment:', error);
      this.setSubmissionLoading(false, parentId);
      this.displaySubmissionError(error.message, parentId);
      return null;
    }
  }

  /**
   * Validate comment input
   * Requirements: 4.1, 4.2, 4.3
   */
  validateCommentInput(content, name) {
    const errors = [];
    
    // Content validation
    if (!content || content.trim().length === 0) {
      errors.push('Comment content is required');
    } else if (content.trim().length < 3) {
      errors.push('Comment must be at least 3 characters long');
    } else if (content.trim().length > 2000) {
      errors.push('Comment must be no more than 2000 characters long');
    }
    
    // Link validation (max 2 links)
    const linkRegex = /https?:\/\/[^\s]+/g;
    const links = content.match(linkRegex) || [];
    if (links.length > 2) {
      errors.push('Comments can contain at most 2 links');
    }
    
    // Name validation
    if (name && name.trim().length > 50) {
      errors.push('Name must be no more than 50 characters long');
    }
    
    return {
      isValid: errors.length === 0,
      errors: errors
    };
  }

  /**
   * Display validation errors
   */
  displayValidationErrors(errors, parentId = null) {
    const errorContainer = this.getErrorContainer(parentId);
    if (errorContainer) {
      errorContainer.innerHTML = `
        <div class="validation-errors">
          ${errors.map(error => `<div class="error-item">${this.escapeHtml(error)}</div>`).join('')}
        </div>
      `;
      errorContainer.style.display = 'block';
    }
  }

  /**
   * Display submission error
   */
  displaySubmissionError(message, parentId = null) {
    const errorContainer = this.getErrorContainer(parentId);
    if (errorContainer) {
      errorContainer.innerHTML = `
        <div class="submission-error">
          Failed to submit comment: ${this.escapeHtml(message)}
        </div>
      `;
      errorContainer.style.display = 'block';
    }
  }

  /**
   * Get error container for form
   */
  getErrorContainer(parentId = null) {
    const formId = parentId ? `reply-form-${parentId}` : 'main-comment-form';
    const form = document.getElementById(formId);
    if (!form) return null;
    
    let errorContainer = form.querySelector('.error-container');
    if (!errorContainer) {
      errorContainer = document.createElement('div');
      errorContainer.className = 'error-container';
      form.insertBefore(errorContainer, form.firstChild);
    }
    
    return errorContainer;
  }

  /**
   * Clear error messages
   */
  clearErrors(parentId = null) {
    const errorContainer = this.getErrorContainer(parentId);
    if (errorContainer) {
      errorContainer.style.display = 'none';
      errorContainer.innerHTML = '';
    }
  }

  /**
   * Set loading state for submission
   */
  setSubmissionLoading(isLoading, parentId = null) {
    const formId = parentId ? `reply-form-${parentId}` : 'main-comment-form';
    const form = document.getElementById(formId);
    if (!form) return;
    
    const submitButton = form.querySelector('.submit-button');
    const textarea = form.querySelector('textarea');
    const nameInput = form.querySelector('input[type="text"]');
    
    if (submitButton) {
      submitButton.disabled = isLoading;
      submitButton.textContent = isLoading ? 'Submitting...' : (parentId ? 'Reply' : 'Post Comment');
    }
    
    if (textarea) textarea.disabled = isLoading;
    if (nameInput) nameInput.disabled = isLoading;
  }

  /**
   * Clear comment form
   */
  clearCommentForm(parentId = null) {
    const formId = parentId ? `reply-form-${parentId}` : 'main-comment-form';
    const form = document.getElementById(formId);
    if (!form) return;
    
    const textarea = form.querySelector('textarea');
    const nameInput = form.querySelector('input[type="text"]');
    
    if (textarea) {
      textarea.value = '';
      // Trigger input event to update character counter
      textarea.dispatchEvent(new Event('input'));
    }
    if (nameInput) nameInput.value = '';
    
    // Clear errors
    this.clearErrors(parentId);
    
    // Hide reply form if it's a reply
    if (parentId) {
      this.cancelReply(parentId);
    }
  }

  /**
   * Generate a unique user fingerprint combining browser characteristics
   * Requirements: 6.1, 6.3
   */
  generateUserFingerprint() {
    const canvas = document.createElement('canvas');
    const ctx = canvas.getContext('2d');
    ctx.textBaseline = 'top';
    ctx.font = '14px Arial';
    ctx.fillText('Browser fingerprint', 2, 2);
    
    const fingerprint = {
      screen: `${screen.width}x${screen.height}x${screen.colorDepth}`,
      timezone: Intl.DateTimeFormat().resolvedOptions().timeZone,
      language: navigator.language,
      platform: navigator.platform,
      userAgent: navigator.userAgent.substring(0, 100), // Truncate for privacy
      canvas: canvas.toDataURL(),
      cookiesEnabled: navigator.cookieEnabled,
      doNotTrack: navigator.doNotTrack,
      hardwareConcurrency: navigator.hardwareConcurrency || 0,
      deviceMemory: navigator.deviceMemory || 0,
      pixelRatio: window.devicePixelRatio || 1
    };

    // Create hash from fingerprint data
    const fingerprintString = JSON.stringify(fingerprint);
    return this.hashString(fingerprintString);
  }

  /**
   * Simple hash function for creating consistent identifiers
   */
  hashString(str) {
    let hash = 0;
    for (let i = 0; i < str.length; i++) {
      const char = str.charCodeAt(i);
      hash = ((hash << 5) - hash) + char;
      hash = hash & hash; // Convert to 32-bit integer
    }
    return Math.abs(hash).toString(36);
  }

  /**
   * Get or create user identity with cookie persistence
   * Requirements: 6.1, 6.2, 6.4, 6.5
   */
  getUserIdentity() {
    const cookieName = 'comment_user_id';
    const sessionKey = 'comment_session_id';
    
    // Try to get existing identity from cookie first
    let userIdentity = this.getCookie(cookieName);
    
    if (userIdentity) {
      try {
        return JSON.parse(userIdentity);
      } catch (e) {
        // Invalid cookie data, generate new identity
        console.warn('Invalid user identity cookie, generating new identity');
      }
    }

    // Fallback to session storage if cookies are disabled
    if (!navigator.cookieEnabled) {
      const sessionIdentity = sessionStorage.getItem(sessionKey);
      if (sessionIdentity) {
        try {
          return JSON.parse(sessionIdentity);
        } catch (e) {
          console.warn('Invalid session identity, generating new identity');
        }
      }
    }

    // Generate new identity
    const browserFingerprint = this.generateUserFingerprint();
    const randomComponent = Math.random().toString(36).substring(2, 15);
    const combinedFingerprint = this.hashString(browserFingerprint + randomComponent);
    
    const identity = {
      fingerprint: combinedFingerprint,
      created: new Date().toISOString(),
      version: '1.0'
    };

    // Store identity with 1-year expiration
    if (navigator.cookieEnabled) {
      this.setCookie(cookieName, JSON.stringify(identity), 365);
    } else {
      // Fallback to session storage
      sessionStorage.setItem(sessionKey, JSON.stringify(identity));
    }

    return identity;
  }

  /**
   * Set cookie with expiration
   * Requirements: 6.2, 6.4
   */
  setCookie(name, value, days) {
    const expires = new Date();
    expires.setTime(expires.getTime() + (days * 24 * 60 * 60 * 1000));
    const expiresString = expires.toUTCString();
    
    document.cookie = `${name}=${encodeURIComponent(value)}; expires=${expiresString}; path=/; SameSite=Lax; Secure=${location.protocol === 'https:'}`;
  }

  /**
   * Get cookie value
   * Requirements: 6.2
   */
  getCookie(name) {
    const nameEQ = name + "=";
    const ca = document.cookie.split(';');
    
    for (let i = 0; i < ca.length; i++) {
      let c = ca[i];
      while (c.charAt(0) === ' ') c = c.substring(1, c.length);
      if (c.indexOf(nameEQ) === 0) {
        return decodeURIComponent(c.substring(nameEQ.length, c.length));
      }
    }
    return null;
  }

  /**
   * Generate deterministic color from user fingerprint
   * Requirements: 3.1, 3.2, 3.4, 3.5
   */
  assignUserColor(fingerprint) {
    // Use fingerprint to generate consistent HSL color
    const hash = this.hashString(fingerprint);
    const hashNum = parseInt(hash.substring(0, 8), 36);
    
    // Generate hue (0-360) from hash
    const hue = hashNum % 360;
    
    // Use fixed saturation and lightness for good contrast
    // Saturation: 65-85% for vibrant but not overwhelming colors
    // Lightness: 45-65% for good contrast on both light and dark backgrounds
    const saturation = 65 + (hashNum % 20); // 65-85%
    const lightness = 45 + (hashNum % 20);  // 45-65%
    
    const color = {
      hsl: `hsl(${hue}, ${saturation}%, ${lightness}%)`,
      hex: this.hslToHex(hue, saturation, lightness),
      hue: hue,
      saturation: saturation,
      lightness: lightness
    };

    // Ensure sufficient contrast by adjusting lightness if needed
    if (this.getContrastRatio(color.hex, '#ffffff') < 3) {
      // Too light, darken it
      color.lightness = Math.max(30, lightness - 20);
      color.hsl = `hsl(${hue}, ${saturation}%, ${color.lightness}%)`;
      color.hex = this.hslToHex(hue, saturation, color.lightness);
    } else if (this.getContrastRatio(color.hex, '#000000') < 3) {
      // Too dark, lighten it
      color.lightness = Math.min(70, lightness + 20);
      color.hsl = `hsl(${hue}, ${saturation}%, ${color.lightness}%)`;
      color.hex = this.hslToHex(hue, saturation, color.lightness);
    }

    return color;
  }

  /**
   * Convert HSL to Hex color
   */
  hslToHex(h, s, l) {
    l /= 100;
    const a = s * Math.min(l, 1 - l) / 100;
    const f = n => {
      const k = (n + h / 30) % 12;
      const color = l - a * Math.max(Math.min(k - 3, 9 - k, 1), -1);
      return Math.round(255 * color).toString(16).padStart(2, '0');
    };
    return `#${f(0)}${f(8)}${f(4)}`;
  }

  /**
   * Calculate contrast ratio between two colors
   */
  getContrastRatio(color1, color2) {
    const lum1 = this.getLuminance(color1);
    const lum2 = this.getLuminance(color2);
    const brightest = Math.max(lum1, lum2);
    const darkest = Math.min(lum1, lum2);
    return (brightest + 0.05) / (darkest + 0.05);
  }

  /**
   * Get relative luminance of a color
   */
  getLuminance(hex) {
    const rgb = this.hexToRgb(hex);
    const rsRGB = rgb.r / 255;
    const gsRGB = rgb.g / 255;
    const bsRGB = rgb.b / 255;

    const r = rsRGB <= 0.03928 ? rsRGB / 12.92 : Math.pow((rsRGB + 0.055) / 1.055, 2.4);
    const g = gsRGB <= 0.03928 ? gsRGB / 12.92 : Math.pow((gsRGB + 0.055) / 1.055, 2.4);
    const b = bsRGB <= 0.03928 ? bsRGB / 12.92 : Math.pow((bsRGB + 0.055) / 1.055, 2.4);

    return 0.2126 * r + 0.7152 * g + 0.0722 * b;
  }

  /**
   * Convert hex color to RGB
   */
  hexToRgb(hex) {
    const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
    return result ? {
      r: parseInt(result[1], 16),
      g: parseInt(result[2], 16),
      b: parseInt(result[3], 16)
    } : null;
  }

  /**
   * Get user color for a specific fingerprint
   * Requirements: 3.1, 3.2, 3.4, 3.5
   */
  getUserColor(fingerprint) {
    return this.assignUserColor(fingerprint);
  }

  /**
   * Apply user color styling to a comment element
   * Requirements: 3.3, 3.4, 3.5
   */
  applyUserColorStyling(element, fingerprint) {
    const color = this.getUserColor(fingerprint);
    
    // Apply color as CSS custom property for flexible styling
    element.style.setProperty('--user-color', color.hsl);
    element.style.setProperty('--user-color-hex', color.hex);
    
    // Apply color as left border and subtle background
    element.style.borderLeft = `4px solid ${color.hsl}`;
    
    // Create a very subtle background tint (8% opacity for better readability)
    const rgbColor = this.hexToRgb(color.hex);
    const backgroundTint = `rgba(${rgbColor.r}, ${rgbColor.g}, ${rgbColor.b}, 0.08)`;
    element.style.backgroundColor = backgroundTint;
    
    // Add data attributes for CSS targeting and accessibility
    element.setAttribute('data-user-color', color.hex);
    element.setAttribute('data-user-fingerprint', fingerprint);
    element.setAttribute('aria-label', `Comment by user with color ${color.hex}`);
    
    // Add a subtle hover effect
    element.addEventListener('mouseenter', () => {
      const hoverTint = `rgba(${rgbColor.r}, ${rgbColor.g}, ${rgbColor.b}, 0.15)`;
      element.style.backgroundColor = hoverTint;
    });
    
    element.addEventListener('mouseleave', () => {
      element.style.backgroundColor = backgroundTint;
    });
    
    return color;
  }

  /**
   * Enhanced loading state management
   * Requirements: 5.6, 8.5
   */
  showLoadingState(message = 'Loading...') {
    const container = document.getElementById('comments-list');
    if (container) {
      container.innerHTML = `
        <div class="loading-state">
          <div class="loading-spinner"></div>
          <div class="loading-text">${this.escapeHtml(message)}</div>
        </div>
      `;
    }
  }

  /**
   * Enhanced error state management
   * Requirements: 5.6
   */
  showErrorState(message, canRetry = true) {
    const container = document.getElementById('comments-list');
    if (container) {
      container.innerHTML = `
        <div class="error-state">
          <div class="error-icon">⚠️</div>
          <div class="error-text">${this.escapeHtml(message)}</div>
          ${canRetry ? `
            <button class="retry-button" onclick="commentSystem.loadComments()">
              Try Again
            </button>
          ` : ''}
        </div>
      `;
    }
  }

  /**
   * Show empty state with call-to-action
   * Requirements: 8.5
   */
  showEmptyState() {
    const container = document.getElementById('comments-list');
    if (container) {
      container.innerHTML = `
        <div class="empty-state">
          <div class="empty-icon">💬</div>
          <div class="empty-title">No comments yet</div>
          <div class="empty-text">Be the first to share your thoughts!</div>
          <button class="cta-button" onclick="document.getElementById('main-comment-content').focus()">
            Write a Comment
          </button>
        </div>
      `;
    }
  }

  /**
   * Add smooth animations for comment interactions
   * Requirements: 8.5
   */
  addCommentAnimations(element) {
    // Fade in animation for new comments
    element.style.opacity = '0';
    element.style.transform = 'translateY(20px)';
    element.style.transition = 'opacity 0.3s ease, transform 0.3s ease';
    
    // Trigger animation after element is in DOM
    setTimeout(() => {
      element.style.opacity = '1';
      element.style.transform = 'translateY(0)';
    }, 10);
    
    // Add subtle scale animation on hover
    element.addEventListener('mouseenter', () => {
      if (!element.classList.contains('animating')) {
        element.style.transform = 'scale(1.02)';
        element.style.transition = 'transform 0.2s ease';
      }
    });
    
    element.addEventListener('mouseleave', () => {
      if (!element.classList.contains('animating')) {
        element.style.transform = 'scale(1)';
      }
    });
  }

  /**
   * Enhanced responsive behavior
   * Requirements: 8.5
   */
  updateResponsiveLayout() {
    const container = document.getElementById('comments-container');
    if (!container) return;
    
    const isMobile = window.innerWidth <= 768;
    const isSmallMobile = window.innerWidth <= 480;
    
    // Adjust comment nesting on mobile
    const comments = container.querySelectorAll('.comment');
    comments.forEach(comment => {
      const level = parseInt(comment.getAttribute('data-level') || '0');
      
      if (isSmallMobile) {
        // Reduce nesting on small screens
        comment.style.marginLeft = level > 0 ? '15px' : '0';
      } else if (isMobile) {
        // Moderate nesting on mobile
        comment.style.marginLeft = level > 0 ? `${Math.min(level * 20, 40)}px` : '0';
      } else {
        // Full nesting on desktop
        comment.style.marginLeft = level > 0 ? `${level * 30}px` : '0';
      }
    });
    
    // Adjust form layout
    const formRows = container.querySelectorAll('.form-row');
    formRows.forEach(row => {
      if (isMobile) {
        row.style.flexDirection = 'column';
        row.style.alignItems = 'stretch';
      } else {
        row.style.flexDirection = 'row';
        row.style.alignItems = 'flex-end';
      }
    });
  }

  /**
   * Generate a set of predefined accessible colors as fallback
   * Requirements: 3.4, 3.5
   */
  getAccessibleColorPalette() {
    return [
      { name: 'blue', hsl: 'hsl(210, 70%, 50%)', hex: '#3182ce' },
      { name: 'green', hsl: 'hsl(120, 60%, 45%)', hex: '#38a169' },
      { name: 'purple', hsl: 'hsl(270, 65%, 55%)', hex: '#805ad5' },
      { name: 'orange', hsl: 'hsl(25, 75%, 55%)', hex: '#dd6b20' },
      { name: 'teal', hsl: 'hsl(180, 60%, 45%)', hex: '#319795' },
      { name: 'pink', hsl: 'hsl(320, 65%, 55%)', hex: '#d53f8c' },
      { name: 'indigo', hsl: 'hsl(240, 70%, 60%)', hex: '#5a67d8' },
      { name: 'red', hsl: 'hsl(0, 65%, 55%)', hex: '#e53e3e' }
    ];
  }

  /**
   * Get fallback color if fingerprint-based color generation fails
   * Requirements: 3.4, 3.5
   */
  getFallbackColor(fingerprint) {
    const palette = this.getAccessibleColorPalette();
    const hash = this.hashString(fingerprint || 'anonymous');
    const index = parseInt(hash.substring(0, 2), 36) % palette.length;
    return palette[index];
  }

  /**
   * Validate that a color meets accessibility requirements
   * Requirements: 3.4, 3.5
   */
  isColorAccessible(color, backgroundColor = '#ffffff') {
    const contrastRatio = this.getContrastRatio(color, backgroundColor);
    return contrastRatio >= 3; // WCAG AA standard for large text
  }

  /**
   * Get the current user's color
   * Requirements: 3.1, 3.2, 3.4
   */
  getCurrentUserColor() {
    return this.userColor;
  }

  /**
   * Initialize the comment system
   * Requirements: 1.5, 2.3, 9.1
   */
  async init() {
    // Create the comments container if it doesn't exist
    this.createCommentsContainer();
    
    // Load existing comments
    await this.loadComments();
    
    // Set up event listeners
    this.setupEventListeners();
    
    // Set up responsive handling
    this.setupResponsiveHandling();
  }

  /**
   * Create the comments container HTML structure
   */
  createCommentsContainer() {
    let container = document.getElementById('comments-container');
    if (!container) {
      container = document.createElement('div');
      container.id = 'comments-container';
      container.className = 'comments-container';
      
      // Find a good place to insert the comments (after main content)
      const main = document.querySelector('main, article, .content, .post-content');
      if (main) {
        main.appendChild(container);
      } else {
        document.body.appendChild(container);
      }
    }

    container.innerHTML = `
      <div class="comments-header">
        <h3 class="comments-title">Comments</h3>
        <div class="comments-count" id="comments-count"></div>
      </div>
      <div id="comments-list"></div>
      <div class="comment-form-container" id="main-comment-form">
        <div class="main-comment-form">
          <h4>Leave a Comment</h4>
          <div class="error-container" style="display: none;"></div>
          <textarea 
            id="main-comment-content" 
            placeholder="Share your thoughts..." 
            rows="4"
            maxlength="2000"
            required
          ></textarea>
          <div class="form-row">
            <input 
              type="text" 
              id="main-comment-name" 
              placeholder="Your name (optional)" 
              maxlength="50"
            />
            <div class="form-actions">
              <button onclick="commentSystem.submitMainComment()" class="submit-button">
                Post Comment
              </button>
            </div>
          </div>
          <div class="character-count">
            <span id="main-char-count">0</span>/2000
          </div>
        </div>
      </div>
    `;
  }

  /**
   * Set up event listeners for the comment system
   */
  setupEventListeners() {
    // Global click handler for reply buttons (using event delegation)
    document.addEventListener('click', (e) => {
      if (e.target.classList.contains('reply-button')) {
        e.preventDefault();
        const commentId = e.target.getAttribute('onclick')?.match(/'([^']+)'/)?.[1];
        const level = parseInt(e.target.getAttribute('onclick')?.match(/(\d+)\)/)?.[1] || '0');
        if (commentId) {
          this.showReplyForm(commentId, level);
        }
      }
    });

    // Set up character counter for main comment form
    setTimeout(() => {
      const mainTextarea = document.getElementById('main-comment-content');
      const mainCounter = document.getElementById('main-char-count');
      if (mainTextarea && mainCounter) {
        mainTextarea.addEventListener('input', () => {
          mainCounter.textContent = mainTextarea.value.length;
          mainCounter.style.color = mainTextarea.value.length > 1900 ? '#e53e3e' : '#718096';
        });

        // Clear errors on input
        mainTextarea.addEventListener('input', () => {
          this.clearErrors();
        });
      }

      // Set up name input clear errors
      const mainNameInput = document.getElementById('main-comment-name');
      if (mainNameInput) {
        mainNameInput.addEventListener('input', () => {
          this.clearErrors();
        });
      }
    }, 100);
  }

  /**
   * Submit main comment form
   * Requirements: 1.1, 1.2, 1.3
   */
  async submitMainComment() {
    const contentTextarea = document.getElementById('main-comment-content');
    const nameInput = document.getElementById('main-comment-name');
    
    if (!contentTextarea) {
      console.error('Main comment form not found');
      return;
    }
    
    const content = contentTextarea.value;
    const name = nameInput ? nameInput.value : '';
    
    await this.submitComment(content, name, null);
  }

  /**
   * Set up responsive handling
   * Requirements: 8.5
   */
  setupResponsiveHandling() {
    // Initial layout update
    this.updateResponsiveLayout();
    
    // Handle window resize
    let resizeTimeout;
    window.addEventListener('resize', () => {
      clearTimeout(resizeTimeout);
      resizeTimeout = setTimeout(() => {
        this.updateResponsiveLayout();
      }, 150);
    });
    
    // Handle orientation change on mobile
    window.addEventListener('orientationchange', () => {
      setTimeout(() => {
        this.updateResponsiveLayout();
      }, 300);
    });
  }

  /**
   * Show reply form for a specific comment
   * Requirements: 2.1, 2.5
   */
  showReplyForm(commentId, level) {
    // Hide any existing reply forms
    document.querySelectorAll('.reply-form-container').forEach(container => {
      if (container.id !== `reply-form-${commentId}`) {
        container.style.display = 'none';
        container.innerHTML = '';
      }
    });

    const container = document.getElementById(`reply-form-${commentId}`);
    if (!container) return;

    if (container.style.display === 'none') {
      container.style.display = 'block';
      container.innerHTML = `
        <div class="reply-form">
          <h4>Reply to comment</h4>
          <textarea 
            id="reply-content-${commentId}" 
            placeholder="Write your reply..." 
            rows="3"
            maxlength="2000"
            required
          ></textarea>
          <div class="form-row">
            <input 
              type="text" 
              id="reply-name-${commentId}" 
              placeholder="Your name (optional)" 
              maxlength="50"
            />
            <div class="form-actions">
              <button onclick="commentSystem.cancelReply('${commentId}')" class="cancel-button">
                Cancel
              </button>
              <button onclick="commentSystem.submitReply('${commentId}', ${level})" class="submit-button">
                Reply
              </button>
            </div>
          </div>
          <div class="character-count">
            <span id="reply-char-count-${commentId}">0</span>/2000
          </div>
        </div>
      `;

      // Add character counter
      const textarea = document.getElementById(`reply-content-${commentId}`);
      const counter = document.getElementById(`reply-char-count-${commentId}`);
      textarea.addEventListener('input', () => {
        counter.textContent = textarea.value.length;
        counter.style.color = textarea.value.length > 1900 ? '#e53e3e' : '#718096';
      });
    } else {
      container.style.display = 'none';
      container.innerHTML = '';
    }
  }

  /**
   * Cancel reply form
   */
  cancelReply(commentId) {
    const container = document.getElementById(`reply-form-${commentId}`);
    if (container) {
      container.style.display = 'none';
      container.innerHTML = '';
    }
  }

  /**
   * Submit a reply to a specific comment
   * Requirements: 1.1, 1.2, 1.3, 2.1, 2.5
   */
  async submitReply(commentId, level) {
    const contentTextarea = document.getElementById(`reply-content-${commentId}`);
    const nameInput = document.getElementById(`reply-name-${commentId}`);
    
    if (!contentTextarea) {
      console.error('Reply form not found');
      return;
    }
    
    const content = contentTextarea.value;
    const name = nameInput ? nameInput.value : '';
    
    // Determine parent ID based on level flattening (3+ levels)
    // Requirements: 2.5 - Implement level flattening for deep comment threads
    let parentId = commentId;
    if (level >= 2) {
      // For level 2 and above, find the root comment to flatten the thread
      const commentElement = document.querySelector(`[data-comment-id="${commentId}"]`);
      if (commentElement) {
        // Find the root comment (level 0) by traversing up
        let currentElement = commentElement;
        while (currentElement && parseInt(currentElement.getAttribute('data-level')) > 0) {
          currentElement = currentElement.parentElement.closest('.comment');
        }
        if (currentElement) {
          parentId = currentElement.getAttribute('data-comment-id');
        }
      }
    }
    
    await this.submitComment(content, name, parentId);
  }

  /**
   * Update comments count display
   */
  updateCommentsCount(count) {
    const countElement = document.getElementById('comments-count');
    if (countElement) {
      if (count === 0) {
        countElement.textContent = '';
      } else {
        countElement.textContent = `${count} comment${count === 1 ? '' : 's'}`;
      }
    }
  }

  /**
   * Count total comments including nested ones
   */
  countComments(comments) {
    let count = 0;
    comments.forEach(comment => {
      count++; // Count this comment
      if (comment.children && comment.children.length > 0) {
        count += this.countComments(comment.children); // Count children recursively
      }
    });
    return count;
  }

  /**
   * Test the user identity and color system
   * For development and debugging purposes
   */
  testUserIdentitySystem() {
    console.log('Testing User Identity System:');
    console.log('User Identity:', this.userIdentity);
    console.log('User Color:', this.userColor);
    
    // Test color consistency
    const testFingerprint = 'test-fingerprint-123';
    const color1 = this.assignUserColor(testFingerprint);
    const color2 = this.assignUserColor(testFingerprint);
    
    console.log('Color Consistency Test:');
    console.log('Color 1:', color1);
    console.log('Color 2:', color2);
    console.log('Colors Match:', JSON.stringify(color1) === JSON.stringify(color2));
    
    // Test accessibility
    console.log('Accessibility Test:');
    console.log('White Background Contrast:', this.getContrastRatio(this.userColor.hex, '#ffffff'));
    console.log('Black Background Contrast:', this.getContrastRatio(this.userColor.hex, '#000000'));
    console.log('Is Accessible (White BG):', this.isColorAccessible(this.userColor.hex, '#ffffff'));
    console.log('Is Accessible (Black BG):', this.isColorAccessible(this.userColor.hex, '#000000'));
    
    return {
      identity: this.userIdentity,
      color: this.userColor,
      consistent: JSON.stringify(color1) === JSON.stringify(color2),
      accessible: this.isColorAccessible(this.userColor.hex)
    };
  }
}