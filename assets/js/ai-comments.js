// Custom error class for comment system
class CommentError extends Error {
  constructor(code, message, status = 0, options = {}) {
    super(message);
    this.name = "CommentError";
    this.code = code;
    this.status = status;
    this.retryable = options.retryable || false;
    this.suggestions = options.suggestions || [];
    this.details = options.details;
    this.requestId = options.requestId;
    this.userMessage = options.userMessage || message;
    this.timestamp = new Date().toISOString();
  }
}

// API Client for comment system
class CommentAPIClient {
  constructor(options = {}) {
    this.apiBase = options.apiBase || "/api";
    this.retryAttempts = options.retryAttempts || 3;
    this.retryDelay = options.retryDelay || 1000;
    this.timeout = options.timeout || 10000;
  }

  async request(endpoint, options = {}) {
    const url = `${this.apiBase}${endpoint}`;
    const config = {
      timeout: this.timeout,
      headers: {
        "Content-Type": "application/json",
        ...options.headers,
      },
      ...options,
    };

    let lastError;

    for (let attempt = 1; attempt <= this.retryAttempts; attempt++) {
      try {
        // Check if we're offline
        if (!navigator.onLine) {
          throw new CommentError("OFFLINE", "You are currently offline", 0, {
            retryable: true,
            userMessage: "Please check your internet connection and try again.",
          });
        }

        const response = await this.fetchWithTimeout(url, config);

        if (!response.ok) {
          const errorData = await response.json().catch(() => ({}));
          const error = this.createCommentError(response, errorData);
          throw error;
        }

        return await response.json();
      } catch (error) {
        lastError = error;

        // Enhanced retry logic
        const shouldRetry = this.shouldRetryRequest(error, attempt);
        if (!shouldRetry) {
          throw this.enhanceError(error);
        }

        // Calculate delay with jitter to avoid thundering herd
        const baseDelay = this.retryDelay * Math.pow(2, attempt - 1);
        const jitter = Math.random() * 0.1 * baseDelay;
        const delay = baseDelay + jitter;

        await new Promise((resolve) => setTimeout(resolve, delay));

        console.warn(
          `API request failed (attempt ${attempt}/${this.retryAttempts}):`,
          error.message
        );
      }
    }

    throw this.enhanceError(lastError);
  }

  createCommentError(response, errorData) {
    const errorInfo = errorData.error || {};
    const error = new CommentError(
      errorInfo.code || "HTTP_ERROR",
      errorInfo.userMessage ||
        errorInfo.message ||
        `HTTP ${response.status}: ${response.statusText}`,
      response.status,
      {
        retryable: errorInfo.retryable || false,
        suggestions: errorInfo.suggestions || [],
        details: errorInfo.details,
        requestId: response.headers.get("X-Request-ID"),
      }
    );

    return error;
  }

  shouldRetryRequest(error, attempt) {
    // Don't retry if it's the last attempt
    if (attempt >= this.retryAttempts) return false;

    // Don't retry client errors (4xx) unless specifically marked as retryable
    if (error.status >= 400 && error.status < 500 && !error.retryable)
      return false;

    // Retry server errors (5xx)
    if (error.status >= 500) return true;

    // Retry network errors
    if (error.code === "NETWORK_ERROR" || error.code === "TIMEOUT") return true;

    // Retry rate limiting after delay
    if (error.code === "RATE_LIMIT_EXCEEDED") return true;

    // Check if error is explicitly marked as retryable
    return error.retryable || false;
  }

  enhanceError(error) {
    if (error instanceof CommentError) {
      return error;
    }

    // Convert generic errors to CommentError
    let code = "UNKNOWN_ERROR";
    let userMessage = "An unexpected error occurred. Please try again.";

    if (error.name === "TypeError" && error.message.includes("fetch")) {
      code = "NETWORK_ERROR";
      userMessage =
        "Network connection failed. Please check your internet connection.";
    } else if (
      error.name === "AbortError" ||
      error.message.includes("timeout")
    ) {
      code = "TIMEOUT";
      userMessage = "The request took too long. Please try again.";
    }

    return new CommentError(code, userMessage, 0, {
      retryable: true,
      originalError: error.message,
    });
  }

  async fetchWithTimeout(url, options) {
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), this.timeout);

    try {
      const response = await fetch(url, {
        ...options,
        signal: controller.signal,
      });
      clearTimeout(timeoutId);
      return response;
    } catch (error) {
      clearTimeout(timeoutId);
      if (error.name === "AbortError") {
        throw new Error("Request timeout");
      }
      throw error;
    }
  }

  // Comment API methods
  async getComments(pageId, orderingMode = "chronological") {
    const endpoint =
      orderingMode === "similarity"
        ? `/comments/similar/${encodeURIComponent(pageId)}`
        : `/comments/${encodeURIComponent(pageId)}`;

    return this.request(endpoint);
  }

  async createComment(commentData) {
    return this.request("/comments", {
      method: "POST",
      body: JSON.stringify(commentData),
    });
  }

  async replyToComment(parentId, replyData) {
    return this.request(`/comments/${parentId}/reply`, {
      method: "POST",
      body: JSON.stringify(replyData),
    });
  }

  // Health check method
  async checkHealth() {
    try {
      await this.request("/health");
      return true;
    } catch (error) {
      return false;
    }
  }
}

// Frontend JavaScript for comment system
class CommentSystem {
  constructor(options = {}) {
    this.pageId = options.pageId || window.location.pathname;
    this.pageTitle = options.pageTitle || ""; // Store page title for new comments
    this.container = null;
    this.comments = [];
    this.orderingMode = "similarity"; // Default to 'similarity' (Relevant) mode
    this.chronologicalOrder = "newest"; // 'newest' (recent) or 'oldest'
    this.isLoading = false;
    this.isOnline = navigator.onLine;
    this.updateInterval = null;
    this.lastUpdateTime = null;
    this.cachedOrderedComments = null;
    this.lastCommentCount = 0;
    this.orderingCache = new Map(); // Cache for different ordering modes

    // Load persistent cache from localStorage
    this.loadPersistentCache();

    // Initialize API client
    this.api = new CommentAPIClient({
      apiBase: options.apiBase || "/api",
      retryAttempts: options.retryAttempts || 3,
      retryDelay: options.retryDelay || 1000,
      timeout: options.timeout || 10000,
    });

    // Initialize the system
    this.init(options.container);
    this.setupNetworkHandling();
    this.setupAutoRefresh(options.autoRefresh);
    this.setupCacheControls();
  }

  init(containerSelector) {
    this.container =
      typeof containerSelector === "string"
        ? document.querySelector(containerSelector)
        : containerSelector;

    if (!this.container) {
      console.error("CommentSystem: Container not found");
      return;
    }

    this.render();
    this.loadComments();
  }

  render() {
    this.container.innerHTML = this.getHTML();
    this.attachEventListeners();
  }

  getHTML() {
    return `
      <div class="comment-system">
        <!-- Comment Form Section -->
        <div class="comment-form-section">
          <h3 class="section-title">Join the Discussion</h3>
          <div class="comment-form-container">
            ${this.getCommentFormHTML()}
          </div>
        </div>
        
        <div class="comments-container">
          <div class="comments-header">
            <h3 class="comments-title">
              Comments
              <span class="comments-count" id="comments-count">0</span>
            </h3>
            <div class="ordering-controls">
              <div class="ordering-toggle">
                <button class="ordering-button" data-mode="chronological">
                  Recent
                </button>
                <button class="ordering-button active" data-mode="similarity">
                  Relevant
                </button>
              </div>
            </div>
          </div>
          
          <div class="comments-list" id="comments-list">
            <!-- Comments will be rendered here -->
          </div>
        </div>
      </div>
    `;
  }

  getCommentFormHTML(parentComment = null) {
    const isReply = !!parentComment;
    const formId = isReply
      ? `reply-form-${parentComment.id}`
      : "main-comment-form";

    return `
      <form class="comment-form ${
        isReply ? "reply-form" : ""
      }" id="${formId}" data-parent-id="${parentComment?.id || ""}">
        <div class="form-header">
          <h4 class="form-title">${
            isReply ? "Reply to Comment" : "Leave a Comment"
          }</h4>
          <p class="form-subtitle">Share your thoughts and join the discussion</p>
        </div>
        
        ${isReply ? this.getReplyContextHTML(parentComment) : ""}
        
        <div class="anonymous-notice">
          <span class="anonymous-notice-icon">ℹ️</span>
          No account required. Leave the name field empty to post anonymously.
        </div>
        
        <div class="form-group">
          <label class="form-label" for="${formId}-name">Name (optional)</label>
          <input 
            type="text" 
            id="${formId}-name" 
            name="authorName"
            class="form-input" 
            placeholder="Your name or leave empty for Anonymous"
            maxlength="50"
          >
          <div class="form-help">Leave empty to post as "Anonymous"</div>
        </div>
        
        <div class="form-group">
          <label class="form-label" for="${formId}-content">Comment *</label>
          <textarea 
            id="${formId}-content" 
            name="content"
            class="form-textarea" 
            placeholder="Share your thoughts..."
            required
            maxlength="2000"
            rows="4"
          ></textarea>
          <div class="character-counter" id="${formId}-counter">0 / 2000</div>
          <div class="form-error" id="${formId}-content-error" style="display: none;"></div>
        </div>
        
        <div class="form-actions">
          ${
            isReply
              ? '<button type="button" class="form-btn form-btn-secondary cancel-reply-btn">Cancel</button>'
              : ""
          }
          <button type="submit" class="form-btn form-btn-primary">
            <span class="btn-text">${
              isReply ? "Post Reply" : "Post Comment"
            }</span>
          </button>
        </div>
        
        <div class="form-success" id="${formId}-success" style="display: none;">
          Your ${isReply ? "reply" : "comment"} has been posted successfully!
        </div>
      </form>
    `;
  }

  getReplyContextHTML(parentComment) {
    const truncatedContent =
      parentComment.content.length > 100
        ? parentComment.content.substring(0, 100) + "..."
        : parentComment.content;

    return `
      <div class="reply-context">
        <div class="reply-context-label">Replying to ${this.escapeHtml(
          parentComment.authorName || "Anonymous"
        )}:</div>
        <div class="reply-context-content">"${this.escapeHtml(
          truncatedContent
        )}"</div>
      </div>
    `;
  }

  attachEventListeners() {
    // Ordering toggle buttons
    const orderingButtons = this.container.querySelectorAll(".ordering-button");
    orderingButtons.forEach((button) => {
      button.addEventListener("click", (e) => {
        const mode = e.target.dataset.mode;
        console.log(
          "🔄 Button clicked:",
          mode,
          "current mode:",
          this.orderingMode,
          "isLoading:",
          this.isLoading
        );

        // Prevent race conditions by checking loading state first
        if (this.isLoading) {
          console.log("🔄 Button click ignored - already loading");
          return;
        }

        // Handle chronological toggle behavior
        if (mode === "chronological" && this.orderingMode === "chronological") {
          // Toggle between newest and oldest
          console.log("🔄 Toggling chronological order");
          this.toggleChronologicalOrder();
        } else if (this.orderingMode !== mode) {
          // Switch to different mode only if not already in that mode
          console.log("🔄 Switching to mode:", mode);
          console.log("🔄 About to call setOrderingMode");
          this.setOrderingMode(mode)
            .then(() => {
              console.log("🔄 setOrderingMode completed");
            })
            .catch((error) => {
              console.error("🔄 setOrderingMode failed:", error);
            });
        } else {
          console.log("🔄 Button click ignored - already in mode:", mode);
        }
      });
    });

    // Main comment form
    const mainForm = this.container.querySelector("#main-comment-form");
    if (mainForm) {
      this.attachFormListeners(mainForm);
    }

    // Reply buttons (will be attached when comments are rendered)
    this.attachReplyButtonListeners();
  }

  attachFormListeners(form) {
    const formId = form.id;
    const contentTextarea = form.querySelector('textarea[name="content"]');
    const counter = form.querySelector(`#${formId}-counter`);
    const submitBtn = form.querySelector('button[type="submit"]');
    const cancelBtn = form.querySelector(".cancel-reply-btn");

    // Character counter
    if (contentTextarea && counter) {
      contentTextarea.addEventListener("input", (e) => {
        this.updateCharacterCounter(e.target, counter);
      });
    }

    // Form submission
    form.addEventListener("submit", (e) => {
      e.preventDefault();
      this.handleFormSubmit(form);
    });

    // Cancel reply button
    if (cancelBtn) {
      cancelBtn.addEventListener("click", () => {
        this.cancelReply(form);
      });
    }

    // Real-time validation
    const nameInput = form.querySelector('input[name="authorName"]');
    if (nameInput) {
      nameInput.addEventListener("blur", () => {
        this.validateField(nameInput);
      });
    }

    if (contentTextarea) {
      contentTextarea.addEventListener("blur", () => {
        this.validateField(contentTextarea);
      });
    }
  }

  attachReplyButtonListeners() {
    // Use event delegation for reply buttons
    this.container.addEventListener("click", (e) => {
      if (
        e.target.classList.contains("reply-btn") ||
        e.target.closest(".reply-btn")
      ) {
        const btn = e.target.classList.contains("reply-btn")
          ? e.target
          : e.target.closest(".reply-btn");
        const commentId = btn.dataset.commentId;
        this.showReplyForm(commentId);
      }
    });
  }

  updateCharacterCounter(textarea, counter) {
    const length = textarea.value.length;
    const maxLength = parseInt(textarea.getAttribute("maxlength"));

    counter.textContent = `${length} / ${maxLength}`;

    // Update counter styling based on length
    counter.classList.remove("warning", "error");
    if (length > maxLength * 0.9) {
      counter.classList.add("warning");
    }
    if (length >= maxLength) {
      counter.classList.add("error");
    }
  }

  validateField(field) {
    const value = field.value.trim();
    const fieldName = field.name;
    let isValid = true;
    let errorMessage = "";

    // Clear previous errors
    field.classList.remove("error");
    const errorElement = field
      .closest(".form-group")
      .querySelector(".form-error");
    if (errorElement) {
      errorElement.style.display = "none";
    }

    // Validate content field
    if (fieldName === "content") {
      if (!value) {
        isValid = false;
        errorMessage = "Comment content is required";
      } else if (value.length < 1) {
        isValid = false;
        errorMessage = "Comment must be at least 1 character long";
      } else if (value.length > 2000) {
        isValid = false;
        errorMessage = "Comment must be less than 2000 characters";
      } else if (!/\S/.test(value)) {
        isValid = false;
        errorMessage = "Comment cannot be only whitespace";
      }
    }

    // Validate name field (optional but has constraints if provided)
    if (fieldName === "authorName" && value) {
      if (value.length > 50) {
        isValid = false;
        errorMessage = "Name must be less than 50 characters";
      }
    }

    // Show error if validation failed
    if (!isValid && errorElement) {
      field.classList.add("error");
      errorElement.textContent = errorMessage;
      errorElement.style.display = "block";
    }

    return isValid;
  }

  async handleFormSubmit(form) {
    const formData = new FormData(form);
    const submitBtn = form.querySelector('button[type="submit"]');
    const parentId = form.dataset.parentId || null;

    // Clear any previous errors
    this.clearFormErrors(form);

    // Validate all fields
    const nameField = form.querySelector('input[name="authorName"]');
    const contentField = form.querySelector('textarea[name="content"]');

    const isNameValid = this.validateField(nameField);
    const isContentValid = this.validateField(contentField);

    if (!isContentValid) {
      contentField.focus();
      return;
    }

    // Prepare comment data
    const commentData = {
      pageId: this.pageId,
      authorName: formData.get("authorName")?.trim() || "Anonymous",
      content: formData.get("content").trim(),
    };

    if (parentId) {
      commentData.parentId = parseInt(parentId);
    }

    // Set loading state
    this.setFormLoading(form, true);

    try {
      let result;
      if (parentId) {
        result = await this.api.replyToComment(parentId, commentData);
      } else {
        result = await this.api.createComment(commentData);
      }

      // Show success message
      this.showFormSuccess(
        form,
        parentId ? "Reply posted successfully!" : "Comment posted successfully!"
      );
      this.showNotification(
        parentId
          ? "Reply posted successfully!"
          : "Comment posted successfully!",
        "success"
      );

      // Reset form
      form.reset();
      this.updateCharacterCounter(
        contentField,
        form.querySelector(".character-counter")
      );

      // Clear frontend cache since new comment was added
      this.orderingCache.clear();
      console.log("🗑️ Cleared frontend cache after new comment posted");

      // Refresh comments to show the new comment
      setTimeout(() => {
        this.loadComments();

        // Hide reply form if it was a reply
        if (parentId) {
          this.cancelReply(form);
        }
      }, 1000);
    } catch (error) {
      console.error("Failed to post comment:", error);
      this.handleFormSubmissionError(form, error);
    } finally {
      this.setFormLoading(form, false);
    }
  }

  handleFormSubmissionError(form, error) {
    if (error instanceof CommentError) {
      // Handle structured comment errors
      this.showFormError(form, error.userMessage, {
        suggestions: error.suggestions,
        retryable: error.retryable,
        code: error.code,
      });

      // Handle specific error types
      switch (error.code) {
        case "VALIDATION_ERROR":
          this.highlightValidationErrors(form, error.details);
          break;
        case "RATE_LIMIT_EXCEEDED":
          this.showRateLimitMessage(form, error.details);
          break;
        case "SPAM_DETECTED":
          this.showSpamDetectionMessage(form, error.details);
          break;
        case "DUPLICATE_CONTENT":
          this.showDuplicateContentMessage(form, error.details);
          break;
        case "OFFLINE":
          this.showOfflineMessage(form);
          break;
      }
    } else {
      // Handle generic errors
      this.showFormError(
        form,
        "An unexpected error occurred. Please try again."
      );
    }
  }

  clearFormErrors(form) {
    // Clear field-level errors
    const errorElements = form.querySelectorAll(".form-error");
    errorElements.forEach((el) => {
      el.style.display = "none";
      el.textContent = "";
    });

    // Clear field error styling
    const errorFields = form.querySelectorAll(".error");
    errorFields.forEach((field) => field.classList.remove("error"));

    // Clear general form errors
    const generalError = form.querySelector(".form-error-general");
    if (generalError) {
      generalError.style.display = "none";
    }
  }

  highlightValidationErrors(form, validationDetails) {
    if (!validationDetails || !Array.isArray(validationDetails)) return;

    validationDetails.forEach((error) => {
      const fieldName = error.path?.[0];
      if (fieldName) {
        const field = form.querySelector(`[name="${fieldName}"]`);
        if (field) {
          field.classList.add("error");
          const errorElement = field
            .closest(".form-group")
            .querySelector(".form-error");
          if (errorElement) {
            errorElement.textContent = error.message;
            errorElement.style.display = "block";
          }
        }
      }
    });
  }

  showRateLimitMessage(form, details) {
    const message = details?.resetTime
      ? `Please wait until ${new Date(
          details.resetTime
        ).toLocaleTimeString()} before posting again.`
      : "You are posting too quickly. Please wait a moment before trying again.";

    this.showFormError(form, message, {
      type: "warning",
      suggestions: [
        "Wait a few minutes before posting again",
        "Avoid posting multiple comments quickly",
      ],
    });
  }

  showSpamDetectionMessage(form, details) {
    const suggestions = details?.reasons || [
      "Remove excessive links",
      "Write more natural content",
      "Avoid repetitive text",
    ];

    this.showFormError(
      form,
      "Your comment was flagged as potential spam. Please revise and try again.",
      {
        type: "warning",
        suggestions,
      }
    );
  }

  showDuplicateContentMessage(form, details) {
    this.showFormError(
      form,
      "This comment appears to be a duplicate. Please post something different.",
      {
        type: "info",
        suggestions: [
          "Modify your comment to make it unique",
          "Add additional thoughts or context",
        ],
      }
    );
  }

  showOfflineMessage(form) {
    this.showFormError(
      form,
      "You are currently offline. Your comment will be posted when your connection is restored.",
      {
        type: "warning",
        suggestions: [
          "Check your internet connection",
          "Try again when back online",
        ],
      }
    );
  }

  setFormLoading(form, loading) {
    const submitBtn = form.querySelector('button[type="submit"]');
    const cancelBtn = form.querySelector(".cancel-reply-btn");
    const inputs = form.querySelectorAll("input, textarea");
    const btnText = submitBtn.querySelector(".btn-text");

    // Update button state
    submitBtn.disabled = loading;
    submitBtn.classList.toggle("form-btn-loading", loading);

    // Update button text with loading indicator
    if (loading) {
      const originalText = btnText.textContent;
      submitBtn.dataset.originalText = originalText;
      btnText.innerHTML = `
        <span class="loading-spinner"></span>
        <span>Posting...</span>
      `;
    } else {
      const originalText = submitBtn.dataset.originalText;
      if (originalText) {
        btnText.textContent = originalText;
      }
    }

    // Disable other form elements
    if (cancelBtn) {
      cancelBtn.disabled = loading;
    }

    inputs.forEach((input) => {
      input.disabled = loading;
    });

    // Add loading overlay to form
    this.toggleFormLoadingOverlay(form, loading);
  }

  toggleFormLoadingOverlay(form, show) {
    let overlay = form.querySelector(".form-loading-overlay");

    if (show && !overlay) {
      overlay = document.createElement("div");
      overlay.className = "form-loading-overlay";
      overlay.innerHTML = `
        <div class="form-loading-content">
          <div class="loading-spinner"></div>
          <div class="loading-text">Posting your comment...</div>
        </div>
      `;
      form.appendChild(overlay);
    } else if (!show && overlay) {
      overlay.remove();
    }
  }

  showFormSuccess(form, message = "Comment posted successfully!") {
    const successElement = form.querySelector(".form-success");
    if (successElement) {
      successElement.textContent = message;
      successElement.style.display = "flex";

      // Add success animation
      successElement.classList.add("success-animation");

      setTimeout(() => {
        successElement.style.display = "none";
        successElement.classList.remove("success-animation");
      }, 3000);
    }
  }

  showFormError(form, message, options = {}) {
    // Create or update error message
    let errorElement = form.querySelector(".form-error-general");
    if (!errorElement) {
      errorElement = document.createElement("div");
      errorElement.className = "form-error form-error-general";
      errorElement.style.marginBottom = "1rem";
      form
        .querySelector(".form-actions")
        .insertAdjacentElement("beforebegin", errorElement);
    }

    // Set error type styling
    const errorType = options.type || "error";
    errorElement.className = `form-error form-error-general form-error-${errorType}`;

    // Create error content
    let errorHTML = `<div class="error-message">${this.escapeHtml(
      message
    )}</div>`;

    // Add suggestions if provided
    if (options.suggestions && options.suggestions.length > 0) {
      errorHTML += '<div class="error-suggestions">';
      errorHTML += '<div class="error-suggestions-title">Suggestions:</div>';
      errorHTML += '<ul class="error-suggestions-list">';
      options.suggestions.forEach((suggestion) => {
        errorHTML += `<li>${this.escapeHtml(suggestion)}</li>`;
      });
      errorHTML += "</ul></div>";
    }

    // Add retry button for retryable errors
    if (options.retryable) {
      errorHTML +=
        '<button type="button" class="error-retry-btn" onclick="this.closest(\'form\').dispatchEvent(new Event(\'submit\'))">Try Again</button>';
    }

    errorElement.innerHTML = errorHTML;
    errorElement.style.display = "block";

    // Auto-hide after delay (longer for errors with suggestions)
    const hideDelay = options.suggestions ? 10000 : 5000;
    setTimeout(() => {
      if (errorElement.style.display !== "none") {
        errorElement.style.display = "none";
      }
    }, hideDelay);
  }

  showReplyForm(commentId) {
    // Find the comment
    const comment = this.findCommentById(parseInt(commentId));
    if (!comment) return;

    // Check if reply form already exists
    const existingForm = this.container.querySelector(
      `#reply-form-${commentId}`
    );
    if (existingForm) {
      existingForm.scrollIntoView({ behavior: "smooth", block: "center" });
      existingForm.querySelector("textarea").focus();
      return;
    }

    // Find the comment element and add reply form
    const commentElement = this.container.querySelector(
      `[data-comment-id="${commentId}"]`
    );
    if (!commentElement) return;

    const replyFormHTML = this.getCommentFormHTML(comment);
    const replyContainer = document.createElement("div");
    replyContainer.innerHTML = replyFormHTML;

    // Insert after the comment actions
    const actionsElement = commentElement.querySelector(".comment-actions");
    actionsElement.insertAdjacentElement(
      "afterend",
      replyContainer.firstElementChild
    );

    // Attach event listeners to the new form
    const newForm = commentElement.querySelector(`#reply-form-${commentId}`);
    this.attachFormListeners(newForm);

    // Focus on the textarea
    const textarea = newForm.querySelector("textarea");
    textarea.focus();

    // Scroll into view
    newForm.scrollIntoView({ behavior: "smooth", block: "center" });
  }

  cancelReply(form) {
    form.remove();
  }

  findCommentById(id) {
    const findInComments = (comments) => {
      for (const comment of comments) {
        if (comment.id === id) return comment;
        if (comment.replies && comment.replies.length > 0) {
          const found = findInComments(comment.replies);
          if (found) return found;
        }
      }
      return null;
    };

    return findInComments(this.comments);
  }

  /**
   * Load persistent cache from localStorage
   */
  loadPersistentCache() {
    try {
      const cacheKey = `commentCache_${this.pageId}`;
      const cached = localStorage.getItem(cacheKey);
      if (cached) {
        const { cache, lastCommentCount, timestamp } = JSON.parse(cached);

        // Check if cache is not too old (3 months)
        const maxAge = 90 * 24 * 60 * 60 * 1000; // 90 days (3 months)
        if (Date.now() - timestamp < maxAge) {
          // Restore cache
          this.orderingCache = new Map(Object.entries(cache));
          this.lastCommentCount = lastCommentCount;
        }
      }
    } catch (error) {
      console.error("Failed to load persistent cache:", error);
    }
  }

  /**
   * Save cache to localStorage
   */
  savePersistentCache() {
    try {
      const cacheKey = `commentCache_${this.pageId}`;
      const cacheData = {
        cache: Object.fromEntries(this.orderingCache),
        lastCommentCount: this.lastCommentCount,
        timestamp: Date.now(),
      };
      localStorage.setItem(cacheKey, JSON.stringify(cacheData));
    } catch (error) {
      console.error("Failed to save persistent cache:", error);
    }
  }

  /**
   * Clear all comment cache (both memory and localStorage)
   */
  clearCache() {
    try {
      // Clear in-memory cache
      this.orderingCache.clear();

      // Clear localStorage cache
      const cacheKey = `commentCache_${this.pageId}`;
      localStorage.removeItem(cacheKey);

      // Reset comment count
      this.lastCommentCount = 0;

      console.log("🗑️ Comment cache cleared for page:", this.pageId);

      // Show notification if available
      if (this.showNotification) {
        this.showNotification("Comment cache cleared successfully", "success");
      }

      return true;
    } catch (error) {
      console.error("Failed to clear cache:", error);
      return false;
    }
  }

  /**
   * Clear cache and reload comments
   */
  async clearCacheAndReload() {
    this.clearCache();

    // Force reload comments from server
    await this.loadComments();

    console.log("🔄 Cache cleared and comments reloaded");
  }

  toggleChronologicalOrder() {
    // Toggle between newest and oldest
    this.chronologicalOrder =
      this.chronologicalOrder === "newest" ? "oldest" : "newest";

    // Update button text
    const chronologicalButton = this.container.querySelector(
      '[data-mode="chronological"]'
    );
    if (chronologicalButton) {
      chronologicalButton.textContent =
        this.chronologicalOrder === "newest" ? "Recent" : "Oldest";
    }

    // Add transition effect
    const commentsList = this.container.querySelector("#comments-list");
    commentsList.classList.add("transitioning");

    // Re-render comments with new order after transition
    setTimeout(() => {
      this.renderComments();
    }, 150);
  }

  async setOrderingMode(mode) {
    console.log(
      "🔄 setOrderingMode called with mode:",
      mode,
      "current mode:",
      this.orderingMode,
      "isLoading:",
      this.isLoading
    );

    if (this.orderingMode === mode || this.isLoading) {
      console.log("🔄 setOrderingMode: Early return - same mode or loading");
      return;
    }

    // Set loading state immediately to prevent race conditions
    this.setLoading(true);

    const previousMode = this.orderingMode;
    this.orderingMode = mode;
    console.log(
      "🔄 setOrderingMode: Mode changed from",
      previousMode,
      "to",
      mode
    );

    // Update button states
    const buttons = this.container.querySelectorAll(".ordering-button");
    buttons.forEach((btn) => {
      btn.classList.toggle("active", btn.dataset.mode === mode);
    });

    // Reset chronological button text when switching to similarity mode
    if (mode === "similarity") {
      this.chronologicalOrder = "newest"; // Reset to default
      const chronologicalButton = this.container.querySelector(
        '[data-mode="chronological"]'
      );
      if (chronologicalButton) {
        chronologicalButton.textContent = "Recent";
      }
    }

    // Add transition effect
    const commentsList = this.container.querySelector("#comments-list");
    if (commentsList) {
      commentsList.classList.add("transitioning");
      // Ensure comments are visible during transition
      commentsList.style.opacity = "0.5";
    }

    // Check if we have cached results for this mode
    const cacheKey =
      mode === "chronological"
        ? `${mode}_${this.chronologicalOrder}_${this.pageId}`
        : `${mode}_${this.pageId}`;

    if (this.orderingCache.has(cacheKey)) {
      // Use cached results instantly
      this.comments = this.orderingCache.get(cacheKey);
      this.setLoading(false);
      setTimeout(() => {
        this.renderComments();
      }, 150);
      return;
    }

    // For similarity mode, use transparent background analysis approach
    if (mode === "similarity") {
      console.log("🔄 [Transparent] Loading similarity comments from API");

      // Preserve current comments in case API fails
      const preservedComments = [...this.comments];

      this.setLoading(true);

      try {
        const data = await this.api.getComments(this.pageId, "similarity");
        console.log("🔄 [Transparent] API response:", data);

        if (
          data &&
          data.success &&
          data.comments &&
          Array.isArray(data.comments)
        ) {
          this.comments = data.comments;
          this.lastCommentCount = this.countAllComments(this.comments);

          // Cache the results
          this.orderingCache.set(cacheKey, [...this.comments]);
          this.savePersistentCache();

          // TRANSPARENT BEHAVIOR: No banners, no notifications
          // If chronological, start silent polling for cache completion
          if (data.orderingType === "chronological") {
            console.log("🔄 [Transparent] Showing chronological, will poll for relevance");
            this.startSilentPolling();
          } else if (data.orderingType === "similarity") {
            console.log("🔄 [Transparent] Showing relevance-sorted comments");
          }
        } else {
          console.log("🔄 [Transparent] Invalid response, keeping existing comments");
          this.comments = preservedComments;
        }

        this.renderComments();
      } catch (error) {
        console.error("🔄 [Transparent] API failed:", error);
        this.comments = preservedComments;
        this.orderingMode = previousMode;
        buttons.forEach((btn) => {
          btn.classList.toggle("active", btn.dataset.mode === previousMode);
        });
        this.renderComments();
      } finally {
        this.setLoading(false);
      }
    } else {
      // For chronological mode
      this.setLoading(false);
      setTimeout(() => {
        this.renderComments();
        setTimeout(() => {
          const commentsList = this.container.querySelector("#comments-list");
          if (
            commentsList &&
            commentsList.classList.contains("transitioning")
          ) {
            commentsList.classList.remove("transitioning", "loaded", "loading");
            commentsList.style.opacity = "";
          }
        }, 1000);
      }, 150);
    }
  }

  /**
   * Silent polling for cache completion (no UI indicators)
   * Polls 3 times at 30-second intervals
   */
  startSilentPolling() {
    let pollCount = 0;
    const maxPolls = 3;
    const pollInterval = 30000; // 30 seconds

    const poll = async () => {
      pollCount++;
      console.log(`🔄 [Transparent] Silent poll ${pollCount}/${maxPolls}`);

      try {
        const data = await this.api.getComments(this.pageId, "similarity");
        
        if (data && data.success && data.orderingType === "similarity") {
          // Cache is ready! Silently switch to relevant view
          console.log("🔄 [Transparent] Cache ready, switching to relevant view");
          this.comments = data.comments;
          this.lastCommentCount = this.countAllComments(this.comments);
          
          // Cache the results
          const cacheKey = `similarity_${this.pageId}`;
          this.orderingCache.set(cacheKey, [...this.comments]);
          this.savePersistentCache();
          
          // Silently re-render with relevant comments
          this.renderComments();
          return; // Stop polling
        }

        // Not ready yet, continue polling if under max
        if (pollCount < maxPolls) {
          setTimeout(poll, pollInterval);
        } else {
          console.log("🔄 [Transparent] Max polls reached, staying with current view");
        }
      } catch (error) {
        console.error("🔄 [Transparent] Poll failed:", error);
        // Stop polling on error
      }
    };

    // Start first poll after 30 seconds
    setTimeout(poll, pollInterval);
  }

  async loadComments() {
    this.setLoading(true);

    try {
      // Check if we're offline
      if (!this.isOnline) {
        throw new CommentError("OFFLINE", "You are currently offline", 0, {
          retryable: true,
          userMessage: "Comments will load when your connection is restored.",
        });
      }

      // TRANSPARENT BEHAVIOR: Always try to load with similarity (Relevant) mode first
      console.log("🔄 [Transparent] Initial load - requesting similarity mode");
      const data = await this.api.getComments(this.pageId, "similarity");
      const newComments = data.comments || [];

      // Check if we need to invalidate cache (new comments added)
      const currentCommentCount = this.countAllComments(newComments);
      const hasNewComments =
        this.lastCommentCount > 0 &&
        currentCommentCount !== this.lastCommentCount;

      if (hasNewComments) {
        // Clear cache when new comments are added (but not on initial load)
        this.orderingCache.clear();
      }

      // Always update the comment count
      this.lastCommentCount = currentCommentCount;

      this.comments = newComments;

      // Cache based on what we received
      if (data.orderingType === "similarity") {
        console.log("🔄 [Transparent] Received similarity-sorted comments");
        const similarityKey = `similarity_${this.pageId}`;
        this.orderingCache.set(similarityKey, [...this.comments]);
      } else {
        console.log("🔄 [Transparent] Received chronological comments, will poll for relevance");
        const chronologicalKey = `chronological_newest_${this.pageId}`;
        this.orderingCache.set(chronologicalKey, [...this.comments]);
        
        // Start silent polling if we got chronological but wanted similarity
        if (this.orderingMode === "similarity") {
          this.startSilentPolling();
        }
      }

      // Save to localStorage for persistence
      this.savePersistentCache();

      this.renderComments();
      this.updateCommentsCount();

      // Update last successful load time
      this.lastUpdateTime = new Date();

      // Clear any previous error states
      this.clearLoadingErrors();
    } catch (error) {
      console.error("Failed to load comments:", error);
      this.handleLoadingError(error);
    } finally {
      this.setLoading(false);
    }
  }

  handleLoadingError(error) {
    let errorMessage = "Failed to load comments. Please try again.";
    let showRetry = true;

    if (error instanceof CommentError) {
      errorMessage = error.userMessage;
      showRetry = error.retryable;

      // Handle specific error types
      switch (error.code) {
        case "OFFLINE":
          this.renderOfflineState();
          return;
        case "TIMEOUT":
          errorMessage =
            "Loading comments is taking longer than expected. Please try again.";
          break;
        case "NETWORK_ERROR":
          errorMessage =
            "Network connection failed. Please check your internet connection.";
          break;
      }
    }

    this.renderError(errorMessage, { showRetry, error });

    // Update AI status to error if similarity mode failed
    if (this.orderingMode === "similarity") {
      this.updateAIStatusIndicator("error");
      this.hideAIProcessingBanner();
    }
  }

  clearLoadingErrors() {
    const errorElement = this.container.querySelector(".comments-error");
    if (errorElement) {
      errorElement.remove();
    }

    const offlineElement = this.container.querySelector(".comments-offline");
    if (offlineElement) {
      offlineElement.remove();
    }
  }

  renderOfflineState() {
    const commentsList = this.container.querySelector("#comments-list");
    commentsList.innerHTML = `
      <div class="comments-offline">
        <div class="comments-offline-icon">📡</div>
        <div class="comments-offline-text">You're offline</div>
        <div class="comments-offline-subtext">Comments will load when your connection is restored</div>
      </div>
    `;
  }

  renderComments() {
    console.log(
      "🔄 renderComments called, mode:",
      this.orderingMode,
      "comments:",
      this.comments.length
    );

    const commentsList = this.container.querySelector("#comments-list");

    if (!this.comments || this.comments.length === 0) {
      console.log("🔄 No comments, showing empty state");
      commentsList.innerHTML = this.getEmptyStateHTML();
      this.completeTransition();
      return;
    }

    // Safety check: ensure comments are visible
    if (commentsList) {
      commentsList.style.opacity = "1";
      commentsList.style.transform = "translateY(0)";
    }

    let html = "";

    try {
      if (this.orderingMode === "similarity") {
        console.log("🔄 Rendering in similarity mode");
        
        // Backend already sorted with boost strategy (2 most recent at top, rest by relevance)
        // Just render in the order received from backend
        html = this.comments
          .map((comment) => this.renderCommentThread(comment))
          .join("");
      } else {
        console.log("🔄 Rendering in chronological mode");
        // Sort chronologically, keeping reply structure intact
        const sortedComments = this.sortCommentsChronologically([
          ...this.comments,
        ]);
        html = sortedComments
          .map((comment) => this.renderCommentThread(comment))
          .join("");
      }

      console.log("🔄 Final HTML length:", html.length);
      commentsList.innerHTML = html;
    } catch (error) {
      console.error("🔄 Error rendering comments:", error);
      commentsList.innerHTML = `<div class="comments-error">Error rendering comments: ${error.message}</div>`;
    }

    this.completeTransition();
  }

  sortCommentsChronologically(comments) {
    // Sort root comments by chronological order, preserving reply structure
    return comments.sort((a, b) => {
      const dateA = new Date(a.createdAt);
      const dateB = new Date(b.createdAt);

      if (this.chronologicalOrder === "newest") {
        return dateB - dateA; // Newest first (Recent)
      } else {
        return dateA - dateB; // Oldest first (Oldest)
      }
    });
  }

  renderSimilarityGroupedComments() {
    console.log("🎯 renderSimilarityGroupedComments called");
    console.log("🎯 Raw comments count:", this.comments.length);

    // Comments from backend are already properly threaded, just use them directly
    const threads = this.comments; // Already organized by backend
    console.log("🎯 Using pre-threaded comments count:", threads.length);

    // Order threads by relevance while preserving reply threading
    const sortedThreads = [...threads].sort((a, b) => {
      if (a.relevanceScore && b.relevanceScore) {
        return b.relevanceScore - a.relevanceScore; // Higher relevance first
      }
      // Fallback: simple heuristic based on content length and recency
      const scoreA =
        (a.content?.length || 0) * 0.1 +
        new Date(a.createdAt).getTime() / 1000000;
      const scoreB =
        (b.content?.length || 0) * 0.1 +
        new Date(b.createdAt).getTime() / 1000000;
      return scoreB - scoreA;
    });

    console.log("🎯 Sorted threads count:", sortedThreads.length);
    const html = sortedThreads
      .map((thread) => this.renderCommentThread(thread))
      .join("");
    console.log("🎯 Generated HTML length:", html.length);

    return html;
  }

  groupCommentsByTopic(threads) {
    // Simple topic grouping based on similarity scores
    // In a real implementation, this would use more sophisticated AI analysis

    const groups = [];
    const processedComments = new Set();

    threads.forEach((thread) => {
      if (processedComments.has(thread.id)) return;

      const similarComments = [thread];
      processedComments.add(thread.id);

      // Find similar comments (simplified logic)
      threads.forEach((otherThread) => {
        if (processedComments.has(otherThread.id)) return;

        if (
          this.calculateContentSimilarity(thread.content, otherThread.content) >
          0.3
        ) {
          similarComments.push(otherThread);
          processedComments.add(otherThread.id);
        }
      });

      // Generate topic title based on content
      const topicTitle = this.generateTopicTitle(similarComments);

      groups.push({
        title: topicTitle,
        comments: similarComments.sort((a, b) => {
          // Sort by similarity score if available, otherwise by date
          if (a.similarityScore && b.similarityScore) {
            return b.similarityScore - a.similarityScore;
          }
          return new Date(b.createdAt) - new Date(a.createdAt);
        }),
      });
    });

    return groups;
  }

  calculateContentSimilarity(content1, content2) {
    // Simple word-based similarity calculation
    const words1 = content1.toLowerCase().split(/\s+/);
    const words2 = content2.toLowerCase().split(/\s+/);

    const commonWords = words1.filter(
      (word) => words2.includes(word) && word.length > 3
    );

    return commonWords.length / Math.max(words1.length, words2.length);
  }

  /**
   * Try to get AI relevance scores without filtering comments
   */
  async tryGetAIRelevanceScores(cacheKey) {
    try {
      console.log("🤖 Trying to get AI relevance scores");

      // Only try if we have comments
      if (!this.comments || this.comments.length === 0) {
        console.log("🤖 No comments to analyze");
        return;
      }

      // Get page content for context
      const pageContext = this.extractPageContext();

      // Prepare comments for analysis (only root comments)
      const commentsForAnalysis = this.comments
        .filter((c) => !c.parentId)
        .map((c) => ({
          id: c.id,
          content: c.content,
          authorName: c.authorName,
          createdAt: c.createdAt,
        }));

      if (commentsForAnalysis.length <= 1) {
        console.log("🤖 Not enough root comments to analyze");
        return;
      }

      console.log("🤖 Analyzing", commentsForAnalysis.length, "root comments");

      // Call LLM relevance API
      const relevanceScores = await this.calculateLLMRelevance(
        commentsForAnalysis,
        pageContext
      );

      // Apply scores to existing comments (don't replace comments)
      this.comments.forEach((comment) => {
        const score = relevanceScores.find((s) => s.commentId === comment.id);
        if (score) {
          comment.relevanceScore = score.relevance;
          comment.topicRelevance = score.topicRelevance;
        }
      });

      console.log("🤖 Applied relevance scores, re-rendering");

      // Cache the results
      this.orderingCache.set(cacheKey, [...this.comments]);
      this.savePersistentCache();

      // Re-render with AI scores
      setTimeout(() => {
        this.renderComments();
      }, 100);
    } catch (error) {
      console.warn("🤖 AI relevance scoring failed:", error);
      // Don't throw - comments are already visible
    }
  }

  /**
   * Order comments by relevance using LLM-based analysis
   */
  async orderCommentsByRelevance() {
    try {
      // Get page content for context (if available)
      const pageContext = this.extractPageContext();

      // Prepare comments for LLM analysis
      const commentsForAnalysis = this.comments
        .filter((c) => !c.parentId) // Only analyze root comments
        .map((c) => ({
          id: c.id,
          content: c.content,
          authorName: c.authorName,
          createdAt: c.createdAt,
        }));

      if (commentsForAnalysis.length <= 1) {
        // Not enough comments to reorder
        return;
      }

      // Call LLM-based relevance API
      const relevanceScores = await this.calculateLLMRelevance(
        commentsForAnalysis,
        pageContext
      );

      // Apply relevance scores to comments
      this.comments.forEach((comment) => {
        const score = relevanceScores.find((s) => s.commentId === comment.id);
        if (score) {
          comment.relevanceScore = score.relevance;
          comment.topicRelevance = score.topicRelevance;
        }
      });

      // Comments are now ordered by relevance (caching handled by caller)
    } catch (error) {
      console.error("Failed to order comments by relevance:", error);
      // Fallback to chronological order - no need to throw
    }
  }

  /**
   * Calculate LLM-based relevance scores for comments (with server-side caching)
   */
  async calculateLLMRelevance(comments, pageContext) {
    try {
      // First, check server-side cache
      const cacheResponse = await this.api.request(
        `/comments/cache/${encodeURIComponent(this.pageId)}`
      );

      if (cacheResponse.cached) {
        console.log("Using server-side cached relevance analysis");
        return cacheResponse.relevanceScores;
      }

      // No cache found, perform analysis
      console.log("No cache found, performing LLM analysis");
      const response = await this.api.request("/comments/analyze-relevance", {
        method: "POST",
        body: JSON.stringify({
          pageId: this.pageId,
          pageContext: pageContext,
          comments: comments,
        }),
      });

      const relevanceScores = response.relevanceScores || [];

      // Cache the results on server for future visitors
      if (relevanceScores.length > 0) {
        try {
          await this.api.request(
            `/comments/cache/${encodeURIComponent(this.pageId)}`,
            {
              method: "POST",
              body: JSON.stringify({
                relevanceScores: relevanceScores,
                commentCount: this.lastCommentCount,
              }),
            }
          );
          console.log(
            "Cached relevance analysis on server for future visitors"
          );
        } catch (cacheError) {
          console.error("Failed to cache analysis on server:", cacheError);
        }
      }

      return relevanceScores;
    } catch (error) {
      console.error("LLM relevance analysis failed:", error);
      // Return default scores (chronological order)
      return comments.map((comment, index) => ({
        commentId: comment.id,
        relevance: 1.0 - index * 0.1, // Slight preference for newer comments
        topicRelevance: 0.5,
        reasoning: "Fallback scoring due to analysis failure",
      }));
    }
  }

  /**
   * Extract page context for relevance analysis
   */
  extractPageContext() {
    // Try to extract meaningful content from the page
    const title = document.title || "";
    const metaDescription =
      document.querySelector('meta[name="description"]')?.content || "";

    // Try to get main content
    let mainContent = "";
    const contentSelectors = [
      "article",
      ".post-content",
      ".entry-content",
      ".content",
      "main",
      ".main-content",
    ];

    for (const selector of contentSelectors) {
      const element = document.querySelector(selector);
      if (element) {
        mainContent = element.textContent?.substring(0, 1000) || "";
        break;
      }
    }

    // Extract post tags
    const tags = [];
    const tagElements = document.querySelectorAll('.post__tags .tags__link');
    tagElements.forEach(el => {
      const tagText = el.textContent.trim().toLowerCase();
      if (tagText) {
        tags.push(tagText);
      }
    });

    return {
      title: title,
      description: metaDescription,
      content: mainContent,
      url: window.location.pathname,
      tags: tags  // Add tags to context
    };
  }

  /**
   * Calculate tag-matching score for a comment
   * Returns a score between 0 and 1 based on how many tag words appear in the comment
   */
  calculateTagMatchScore(commentText, pageTags) {
    if (!pageTags || pageTags.length === 0) {
      return 0.5; // Neutral score if no tags
    }
    
    const commentLower = commentText.toLowerCase();
    let matchCount = 0;
    let totalTagWords = 0;
    
    // Break tags into individual words and check for partial matches
    pageTags.forEach(tag => {
      const tagWords = tag.split(/[\s-]+/); // Split on spaces and hyphens
      tagWords.forEach(word => {
        if (word.length <= 2) return; // Skip very short words
        totalTagWords++;
        
        // Check for partial match (word appears in comment)
        if (commentLower.includes(word)) {
          matchCount++;
        }
      });
    });
    
    if (totalTagWords === 0) {
      return 0.5; // Neutral score
    }
    
    // Calculate score: 0.5 base + 0.5 * (match ratio)
    // This ensures scores range from 0.5 (no matches) to 1.0 (all tags match)
    const matchRatio = matchCount / totalTagWords;
    return 0.5 + (matchRatio * 0.5);
  }

  /**
   * Count all comments including replies
   */
  countAllComments(comments) {
    let count = 0;
    const countRecursive = (commentList) => {
      commentList.forEach((comment) => {
        count++;
        if (comment.replies && comment.replies.length > 0) {
          countRecursive(comment.replies);
        }
      });
    };
    countRecursive(comments);
    return count;
  }

  completeTransition() {
    const commentsList = this.container.querySelector("#comments-list");

    if (!commentsList) {
      console.log("🔄 completeTransition: No comments list found");
      return;
    }

    console.log("🔄 completeTransition: Starting transition completion");

    // Remove skeleton loading
    const skeleton = commentsList.querySelector(".comments-skeleton");
    if (skeleton) {
      skeleton.remove();
    }

    // Handle transition animations
    if (commentsList.classList.contains("transitioning")) {
      console.log("🔄 completeTransition: Handling transition animation");

      // Simple transition: fade in
      setTimeout(() => {
        commentsList.classList.add("loaded");
        commentsList.style.opacity = "1";
        commentsList.style.transform = "translateY(0)";

        // Animate individual comments
        this.animateCommentsIn();

        setTimeout(() => {
          console.log("🔄 completeTransition: Removing transition classes");
          commentsList.classList.remove("transitioning", "loaded", "loading");
          commentsList.style.opacity = "";
          commentsList.style.transform = "";
        }, 300);
      }, 100);
    } else {
      console.log(
        "🔄 completeTransition: No transition, just animating comments in"
      );
      // No transition, just animate comments in
      this.animateCommentsIn();
    }
  }

  animateCommentsIn() {
    const comments = this.container.querySelectorAll(".comment-thread");
    comments.forEach((comment, index) => {
      comment.style.opacity = "0";
      comment.style.transform = "translateY(20px)";

      setTimeout(() => {
        comment.style.transition = "opacity 0.3s ease, transform 0.3s ease";
        comment.style.opacity = "1";
        comment.style.transform = "translateY(0)";

        // Clear inline styles after animation
        setTimeout(() => {
          comment.style.opacity = "";
          comment.style.transform = "";
          comment.style.transition = "";
        }, 300);
      }, index * 50); // Stagger animations
    });
  }

  // REMOVED: All banner and notification methods for transparent behavior
  // No showAIProcessingBanner, hideAIProcessingBanner, or updateAIStatusIndicator
  // System operates completely transparently with no UI indicators

  buildCommentThreads(comments) {
    // Comments from the backend are already properly threaded with replies nested
    // We just need to return them as-is since getCommentsWithReplies already organized them
    return comments.map((comment) => ({
      ...comment,
      replies: comment.replies || [],
    }));
  }

  renderCommentThread(comment, depth = 0) {
    if (!comment) return "";

    const threadClass = depth > 0 ? `thread-depth-${Math.min(depth, 4)}` : "";
    const isAnonymous =
      !comment.authorName || comment.authorName === "Anonymous";

    let html = `
      <div class="comment-thread">
        <div class="comment ${threadClass}" data-comment-id="${comment.id}">
          <div class="comment-header">
            <div class="comment-author ${isAnonymous ? "anonymous" : ""}">
              ${this.escapeHtml(comment.authorName || "Anonymous")}
            </div>
            <div class="comment-meta">
              <span class="comment-date">${this.formatDate(
                comment.createdAt
              )}</span>
            </div>
          </div>
          
          <div class="comment-content">
            ${
              this.processCommentContentSync
                ? this.processCommentContentSync(comment.content)
                : this.escapeHtml(comment.content)
            }
          </div>
          
          <div class="comment-actions">
            <button class="comment-action-btn reply-btn" data-comment-id="${
              comment.id
            }">
              <span>↳</span> Reply
            </button>
          </div>
        </div>
    `;

    // Render replies - this is the critical part for showing nested replies
    if (
      comment.replies &&
      Array.isArray(comment.replies) &&
      comment.replies.length > 0
    ) {
      console.log(
        `🔄 Rendering ${comment.replies.length} replies for comment ${comment.id}`
      );
      html += '<div class="comment-replies">';
      comment.replies.forEach((reply) => {
        html += this.renderCommentThread(reply, depth + 1);
      });
      html += "</div>";
    }

    html += "</div>";
    return html;
  }

  renderSimilarityIndicator(score) {
    const percentage = Math.round(score * 100);
    return `
      <span class="similarity-indicator" title="Topic similarity: ${percentage}%">
        <span class="similarity-dot"></span>
        ${percentage}% similar
        <div class="similarity-bar">
          <div class="similarity-fill" style="width: ${percentage}%"></div>
        </div>
      </span>
    `;
  }

  getEmptyStateHTML() {
    return `
      <div class="comments-empty">
        <div class="comments-empty-icon">💬</div>
        <div class="comments-empty-text">No comments yet</div>
        <div class="comments-empty-subtext">Be the first to share your thoughts!</div>
      </div>
    `;
  }

  renderError(message, options = {}) {
    const commentsList = this.container.querySelector("#comments-list");

    let errorHTML = `
      <div class="comments-error">
        <div class="comments-error-icon">⚠️</div>
        <div class="comments-error-text">${this.escapeHtml(message)}</div>
    `;

    if (options.showRetry) {
      errorHTML += `
        <button class="comments-error-retry" onclick="this.closest('.comment-system').commentSystem.loadComments()">
          Try Again
        </button>
      `;
    }

    // Add error details for debugging (in development)
    if (options.error && window.location.hostname === "localhost") {
      errorHTML += `
        <details class="comments-error-details">
          <summary>Error Details</summary>
          <pre>${this.escapeHtml(
            JSON.stringify(
              {
                code: options.error.code,
                message: options.error.message,
                status: options.error.status,
              },
              null,
              2
            )
          )}</pre>
        </details>
      `;
    }

    errorHTML += "</div>";

    commentsList.innerHTML = errorHTML;

    // Store reference to comment system for retry button
    this.container.commentSystem = this;
  }

  setLoading(loading) {
    this.isLoading = loading;
    const indicator = this.container.querySelector("#loading-indicator");
    const buttons = this.container.querySelectorAll(".ordering-button");
    const commentsList = this.container.querySelector("#comments-list");

    // Update loading indicator
    if (indicator) {
      indicator.style.display = loading ? "flex" : "none";

      // Update loading text based on mode
      const loadingText = indicator.querySelector("span");
      if (loadingText) {
        if (this.orderingMode === "similarity") {
          loadingText.textContent = "Analyzing comments...";
        } else {
          loadingText.textContent = "Loading comments...";
        }
      }
    }

    // Disable ordering buttons during loading
    buttons.forEach((btn) => {
      btn.disabled = loading;
      btn.classList.toggle("loading", loading);
    });

    // Add loading state to comments list
    if (commentsList) {
      commentsList.classList.toggle("loading", loading);

      // Add skeleton loading if no comments exist yet
      if (loading && commentsList.children.length === 0) {
        this.showSkeletonLoading(commentsList);
      }
    }
  }

  showSkeletonLoading(container) {
    const skeletonHTML = `
      <div class="comments-skeleton">
        ${Array(3)
          .fill(0)
          .map(
            () => `
          <div class="comment-skeleton">
            <div class="comment-skeleton-header">
              <div class="comment-skeleton-avatar"></div>
              <div class="comment-skeleton-meta">
                <div class="comment-skeleton-name"></div>
                <div class="comment-skeleton-date"></div>
              </div>
            </div>
            <div class="comment-skeleton-content">
              <div class="comment-skeleton-line"></div>
              <div class="comment-skeleton-line"></div>
              <div class="comment-skeleton-line short"></div>
            </div>
          </div>
        `
          )
          .join("")}
      </div>
    `;

    container.innerHTML = skeletonHTML;
  }

  updateCommentsCount() {
    const countElement = this.container.querySelector("#comments-count");
    if (countElement) {
      const totalCount = this.countAllComments(this.comments);
      countElement.textContent = totalCount;
    }
  }

  formatDate(dateString) {
    const date = new Date(dateString);
    const now = new Date();
    const diffMs = now - date;
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMs / 3600000);
    const diffDays = Math.floor(diffMs / 86400000);

    if (diffMins < 1) return "Just now";
    if (diffMins < 60) return `${diffMins}m ago`;
    if (diffHours < 24) return `${diffHours}h ago`;
    if (diffDays < 7) return `${diffDays}d ago`;

    return date.toLocaleDateString();
  }

  escapeHtml(text) {
    const div = document.createElement("div");
    div.textContent = text;
    return div.innerHTML;
  }

  processCommentContentSync(content) {
    // First escape HTML to prevent XSS
    let processedContent = this.escapeHtml(content);

    // URL regex pattern to match various URL formats
    const urlRegex = /(https?:\/\/[^\s<>"{}|\\^`[\]]+)/gi;

    // Replace URLs with clickable links
    processedContent = processedContent.replace(urlRegex, (url) => {
      // Clean up the URL (remove trailing punctuation that might not be part of the URL)
      const cleanUrl = url.replace(/[.,;:!?]+$/, "");
      const trailingPunct = url.slice(cleanUrl.length);

      // Create a safe display version (truncate very long URLs)
      let displayUrl = cleanUrl;
      if (displayUrl.length > 50) {
        displayUrl = displayUrl.substring(0, 47) + "...";
      }

      return `<a href="${cleanUrl}" target="_blank" rel="noopener noreferrer" class="comment-link">${displayUrl}</a>${trailingPunct}`;
    });

    return processedContent;
  }

  // Method to refresh comments (useful for after posting new comments)
  refresh() {
    this.loadComments();
  }

  setupNetworkHandling() {
    // Handle online/offline status
    this.handleOnline = () => {
      this.isOnline = true;
      console.log("Connection restored");
      this.showNetworkStatus("online");

      // Refresh comments when back online
      setTimeout(() => {
        this.loadComments();
      }, 1000);
    };

    this.handleOffline = () => {
      this.isOnline = false;
      console.log("Connection lost");
      this.showNetworkStatus("offline");
    };

    window.addEventListener("online", this.handleOnline);
    window.addEventListener("offline", this.handleOffline);

    // Initial network status check
    if (!navigator.onLine) {
      this.isOnline = false;
      this.showNetworkStatus("offline");
    }

    // Check connection health periodically
    setInterval(() => {
      this.checkConnectionHealth();
    }, 30000); // Check every 30 seconds
  }

  showNetworkStatus(status) {
    // Remove existing network status
    const existingStatus = this.container.querySelector(".network-status");
    if (existingStatus) {
      existingStatus.remove();
    }

    if (status === "offline") {
      const statusElement = document.createElement("div");
      statusElement.className = "network-status network-status-offline";
      statusElement.innerHTML = `
        <div class="network-status-icon">📡</div>
        <div class="network-status-message">
          <div class="network-status-title">You're offline</div>
          <div class="network-status-subtitle">Comments will be available when your connection is restored</div>
        </div>
      `;

      this.container.insertBefore(statusElement, this.container.firstChild);
    } else if (status === "online") {
      const statusElement = document.createElement("div");
      statusElement.className = "network-status network-status-online";
      statusElement.innerHTML = `
        <div class="network-status-icon">✅</div>
        <div class="network-status-message">
          <div class="network-status-title">Connection restored</div>
          <div class="network-status-subtitle">Refreshing comments...</div>
        </div>
      `;

      this.container.insertBefore(statusElement, this.container.firstChild);

      // Auto-remove after 3 seconds
      setTimeout(() => {
        if (statusElement.parentNode) {
          statusElement.remove();
        }
      }, 3000);
    }
  }

  async checkConnectionHealth() {
    if (!navigator.onLine) return;

    try {
      const isHealthy = await this.api.checkHealth();
      if (!isHealthy && this.isOnline) {
        this.showNetworkStatus("degraded");
      }
    } catch (error) {
      console.warn("Health check failed:", error);
    }
  }

  setupAutoRefresh(autoRefresh = true) {
    if (!autoRefresh) return;

    // Refresh comments every 2 minutes
    this.updateInterval = setInterval(() => {
      if (this.isOnline && !this.isLoading) {
        this.refreshComments();
      }
    }, 120000);
  }

  setupCacheControls() {
    // Make cache clearing available globally for debugging
    if (typeof window !== "undefined") {
      window.clearCommentCache = () => {
        return this.clearCacheAndReload();
      };

      // Add keyboard shortcut: Ctrl+Shift+C (or Cmd+Shift+C on Mac)
      document.addEventListener("keydown", (e) => {
        if ((e.ctrlKey || e.metaKey) && e.shiftKey && e.key === "C") {
          e.preventDefault();
          console.log("🗑️ Clearing comment cache via keyboard shortcut...");
          this.clearCacheAndReload();
        }
      });

      console.log("💡 Cache controls available:");
      console.log("   - clearCommentCache() - Clear cache and reload");
      console.log("   - Ctrl+Shift+C (Cmd+Shift+C) - Keyboard shortcut");
    }
  }

  async checkConnection() {
    const isHealthy = await this.api.checkHealth();
    if (isHealthy !== this.isOnline) {
      this.isOnline = isHealthy;
      this.showNetworkStatus(isHealthy ? "online" : "offline");
    }
  }

  showNetworkStatus(status) {
    // Remove existing status
    const existingStatus = this.container.querySelector(".network-status");
    if (existingStatus) {
      existingStatus.remove();
    }

    if (status === "offline") {
      const statusBar = document.createElement("div");
      statusBar.className = "network-status offline";
      statusBar.innerHTML = `
        <div style="background: #ea4335; color: white; padding: 0.5rem 1rem; text-align: center; font-size: 0.875rem;">
          ⚠️ You're offline. Comments will sync when connection is restored.
        </div>
      `;

      this.container.insertBefore(statusBar, this.container.firstChild);
    }
  }

  async refreshComments() {
    try {
      const data = await this.api.getComments(this.pageId, this.orderingMode);

      // Check if there are new comments
      const newCommentCount = data.comments.length;
      const currentCount = this.comments.length;

      if (newCommentCount > currentCount) {
        this.comments = data.comments;
        this.renderComments();
        this.updateCommentsCount();
        this.showNewCommentsNotification(newCommentCount - currentCount);
      }

      this.lastUpdateTime = new Date();
    } catch (error) {
      console.warn("Failed to refresh comments:", error.message);
    }
  }

  showNewCommentsNotification(count) {
    // Remove existing notification
    const existing = this.container.querySelector(".new-comments-notification");
    if (existing) {
      existing.remove();
    }

    const notification = document.createElement("div");
    notification.className = "new-comments-notification";
    notification.innerHTML = `
      <div style="background: #34a853; color: white; padding: 0.75rem 1rem; border-radius: 4px; margin-bottom: 1rem; display: flex; align-items: center; justify-content: space-between;">
        <span>🔔 ${count} new comment${count !== 1 ? "s" : ""} added</span>
        <button onclick="this.parentElement.parentElement.remove()" style="background: none; border: none; color: white; cursor: pointer; font-size: 1.2rem;">×</button>
      </div>
    `;

    const commentsContainer = this.container.querySelector(
      ".comments-container"
    );
    commentsContainer.insertBefore(notification, commentsContainer.firstChild);

    // Auto-hide after 5 seconds
    setTimeout(() => {
      if (notification.parentElement) {
        notification.remove();
      }
    }, 5000);
  }

  // Enhanced loadComments method using API client
  async loadComments() {
    this.setLoading(true);

    try {
      // Add extra delay for similarity mode to show AI processing
      if (this.orderingMode === "similarity") {
        await new Promise((resolve) => setTimeout(resolve, 800));
      }

      const data = await this.api.getComments(this.pageId, this.orderingMode);
      this.comments = data.comments || [];
      this.renderComments();
      this.updateCommentsCount();
      this.lastUpdateTime = new Date();

      // Hide AI processing banner after successful load
      if (this.orderingMode === "similarity") {
        setTimeout(() => {
          this.hideAIProcessingBanner();
        }, 500);
      }
    } catch (error) {
      console.error("Failed to load comments:", error);
      this.renderError(this.getErrorMessage(error));

      // Update AI status to error if similarity mode failed
      if (this.orderingMode === "similarity") {
        this.updateAIStatusIndicator("error");
        this.hideAIProcessingBanner();
      }
    } finally {
      this.setLoading(false);
    }
  }

  // Enhanced form submission using API client
  async handleFormSubmit(form) {
    const formData = new FormData(form);
    const submitBtn = form.querySelector('button[type="submit"]');
    const parentId = form.dataset.parentId || null;

    // Validate all fields
    const nameField = form.querySelector('input[name="authorName"]');
    const contentField = form.querySelector('textarea[name="content"]');

    const isNameValid = this.validateField(nameField);
    const isContentValid = this.validateField(contentField);

    if (!isContentValid) {
      contentField.focus();
      return;
    }

    // Check if online
    if (!this.isOnline) {
      this.showFormError(
        form,
        "You are currently offline. Please check your connection and try again."
      );
      return;
    }

    // Prepare comment data
    const commentData = {
      pageId: this.pageId,
      authorName: formData.get("authorName")?.trim() || "Anonymous",
      content: formData.get("content").trim(),
    };

    if (parentId) {
      commentData.parentId = parseInt(parentId);
    }

    // Set loading state
    this.setFormLoading(form, true);

    try {
      const result = parentId
        ? await this.api.replyToComment(parentId, commentData)
        : await this.api.createComment(commentData);

      // Show success message
      this.showFormSuccess(form);

      // Reset form
      form.reset();
      this.updateCharacterCounter(
        contentField,
        form.querySelector(".character-counter")
      );

      // Refresh comments to show the new comment
      setTimeout(() => {
        this.loadComments();

        // Hide reply form if it was a reply
        if (parentId) {
          this.cancelReply(form);
        }
      }, 1000);
    } catch (error) {
      console.error("Failed to post comment:", error);
      this.showFormError(form, this.getErrorMessage(error));
    } finally {
      this.setFormLoading(form, false);
    }
  }

  getErrorMessage(error) {
    if (!this.isOnline) {
      return "You are currently offline. Please check your connection and try again.";
    }

    if (error.status === 429) {
      return "You are posting too quickly. Please wait a moment before trying again.";
    }

    if (error.status === 400) {
      return (
        error.data?.error?.message ||
        "Invalid comment data. Please check your input."
      );
    }

    if (error.status >= 500) {
      return "Server error. Please try again in a moment.";
    }

    return error.message || "An unexpected error occurred. Please try again.";
  }

  // Method to add a new comment to the display without full reload
  addComment(comment) {
    this.comments.push(comment);
    this.renderComments();
    this.updateCommentsCount();
  }

  // Method to manually refresh comments
  refresh() {
    this.loadComments();
  }

  /**
   * Show AI processing banner (disabled - no banners)
   */
  showAIProcessingBanner() {
    // No-op: banners disabled
  }

  /**
   * Hide AI processing banner (disabled - no banners)
   */
  hideAIProcessingBanner() {
    // No-op: banners disabled
  }

  /**
   * Update AI status indicator (disabled - no visual indicators)
   */
  updateAIStatusIndicator(status) {
    // No-op: status indicators disabled
  }

  // Notification system for user feedback
  showNotification(message, type = "info", duration = 5000) {
    // Remove existing notifications
    const existingNotifications = document.querySelectorAll(
      ".comment-notification"
    );
    existingNotifications.forEach((notification) => notification.remove());

    const notification = document.createElement("div");
    notification.className = `comment-notification comment-notification-${type}`;

    const icons = {
      success: "✅",
      error: "❌",
      warning: "⚠️",
      info: "ℹ️",
    };

    notification.innerHTML = `
      <div class="notification-icon">${icons[type] || icons.info}</div>
      <div class="notification-message">${this.escapeHtml(message)}</div>
      <button class="notification-close" onclick="this.parentElement.remove()">×</button>
    `;

    // Position notification
    notification.style.position = "fixed";
    notification.style.top = "20px";
    notification.style.right = "20px";
    notification.style.zIndex = "10000";

    document.body.appendChild(notification);

    // Auto-remove after duration
    if (duration > 0) {
      setTimeout(() => {
        if (notification.parentNode) {
          notification.style.opacity = "0";
          setTimeout(() => notification.remove(), 300);
        }
      }, duration);
    }

    return notification;
  }

  // Progress indicator for long operations
  showProgress(message, progress = 0) {
    let progressElement = this.container.querySelector(".comment-progress");

    if (!progressElement) {
      progressElement = document.createElement("div");
      progressElement.className = "comment-progress";
      this.container.appendChild(progressElement);
    }

    progressElement.innerHTML = `
      <div class="progress-content">
        <div class="progress-message">${this.escapeHtml(message)}</div>
        <div class="progress-bar">
          <div class="progress-fill" style="width: ${Math.min(
            100,
            Math.max(0, progress)
          )}%"></div>
        </div>
        <div class="progress-percentage">${Math.round(progress)}%</div>
      </div>
    `;

    progressElement.style.display = "block";

    // Auto-hide when complete
    if (progress >= 100) {
      setTimeout(() => {
        progressElement.style.display = "none";
      }, 1000);
    }
  }

  hideProgress() {
    const progressElement = this.container.querySelector(".comment-progress");
    if (progressElement) {
      progressElement.style.display = "none";
    }
  }

  // Cleanup method
  destroy() {
    if (this.updateInterval) {
      clearInterval(this.updateInterval);
    }

    window.removeEventListener("online", this.handleOnline);
    window.removeEventListener("offline", this.handleOffline);

    // Remove notifications
    const notifications = document.querySelectorAll(".comment-notification");
    notifications.forEach((notification) => notification.remove());
  }

  /**
   * Show AI processing banner
   */
  showAIProcessingBanner() {
    // Remove existing banner
    this.hideAIProcessingBanner();

    const banner = document.createElement("div");
    banner.className = "ai-processing-banner";
    banner.innerHTML = `
      <div style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 1rem; border-radius: 8px; margin-bottom: 1rem; display: flex; align-items: center; gap: 1rem;">
        <div class="loading-spinner" style="width: 20px; height: 20px; border: 2px solid rgba(255,255,255,0.3); border-top-color: white; border-radius: 50%; animation: spin 1s linear infinite;"></div>
        <div style="flex: 1;">
          <div style="font-weight: 600; margin-bottom: 0.25rem;">AI Analysis in Progress</div>
          <div style="font-size: 0.875rem; opacity: 0.9;">Analyzing comments for relevance. Showing recent comments for now...</div>
        </div>
      </div>
    `;

    const commentsContainer = this.container.querySelector(".comments-container");
    if (commentsContainer) {
      commentsContainer.insertBefore(banner, commentsContainer.firstChild);
    }
  }

  /**
   * Hide AI processing banner
   */
  hideAIProcessingBanner() {
    const banner = this.container.querySelector(".ai-processing-banner");
    if (banner) {
      banner.remove();
    }
  }

  /**
   * Update AI status indicator (disabled - no visual indicators)
   */
  updateAIStatusIndicator(status) {
    // No-op: status indicators disabled
  }

  /**
   * Start polling for analysis completion
   */
  startAnalysisPolling() {
    // Clear any existing polling
    this.stopAnalysisPolling();

    console.log("🔄 Starting analysis polling for", this.pageId);

    let pollCount = 0;
    const maxPolls = 30; // Poll for up to 30 times (30 seconds with 1s interval)

    this.analysisPollingInterval = setInterval(async () => {
      pollCount++;

      try {
        console.log(`🔄 Polling attempt ${pollCount}/${maxPolls}`);

        const response = await this.api.request(
          `/comments/status/${encodeURIComponent(this.pageId)}`
        );

        if (response.success && response.status === "complete") {
          console.log("🔄 Analysis complete! Updating comments");

          // Stop polling
          this.stopAnalysisPolling();

          // Update comments with analyzed results
          if (response.comments && Array.isArray(response.comments)) {
            this.comments = response.comments;
            this.lastCommentCount = this.countAllComments(this.comments);

            // Update cache
            const cacheKey = `${this.orderingMode}_${this.pageId}`;
            this.orderingCache.set(cacheKey, [...this.comments]);
            this.savePersistentCache();

            // Update UI
            this.hideAIProcessingBanner();
            this.updateAIStatusIndicator("complete");
            this.renderComments();

            if (this.showNotification) {
              this.showNotification("Comments updated", "success");
            }
          }
        } else if (pollCount >= maxPolls) {
          console.log("🔄 Polling timeout, stopping");
          this.stopAnalysisPolling();

          // Update UI to show timeout
          this.hideAIProcessingBanner();
          this.updateAIStatusIndicator("error");

          if (this.showNotification) {
            this.showNotification(
              "Analysis is taking longer than expected. Comments remain in recent order.",
              "warning"
            );
          }
        }
      } catch (error) {
        console.error("🔄 Polling error:", error);

        if (pollCount >= maxPolls) {
          this.stopAnalysisPolling();
          this.hideAIProcessingBanner();
          this.updateAIStatusIndicator("error");
        }
      }
    }, 1000); // Poll every second
  }

  /**
   * Stop analysis polling
   */
  stopAnalysisPolling() {
    if (this.analysisPollingInterval) {
      clearInterval(this.analysisPollingInterval);
      this.analysisPollingInterval = null;
      console.log("🔄 Stopped analysis polling");
    }
  }
}

// Export for use in web pages
window.CommentSystem = CommentSystem;
window.CommentAPIClient = CommentAPIClient;
window.CommentError = CommentError;
