/**
 * Cactus Comments Reply System - DISABLED
 * This file is disabled to prevent conflicts with the chat.html reply system
 */

(function() {
  console.log("Reply system script disabled - using chat.html version instead");
  return; // Exit early to prevent execution
})();

// The rest of this file is disabled
if (false) {
class CactusReplySystem {
  constructor() {
    console.log("CactusReplySystem constructor called");
    this.replyButtons = new Map();
    this.replyForms = new Map();
    this.pendingReply = null;
    this.monitorAttempts = 0;
    this.init();
  }

  init() {
    console.log("Initializing reply system");
    // Wait for cactus comments to load
    this.waitForComments();

    // Set up mutation observer to handle dynamically loaded comments
    this.setupMutationObserver();
  }

  waitForComments() {
    console.log("Waiting for comments to load");
    const checkForComments = () => {
      // Try multiple selectors to find comments
      let comments = document.querySelectorAll(".cactus-comment");

      if (comments.length === 0) {
        // Try other possible selectors
        comments = document.querySelectorAll('[class*="comment"]');
        console.log(
          `Found ${comments.length} elements with "comment" in class name`
        );
      }

      if (comments.length === 0) {
        // Look for any elements in the comment section
        const commentSection = document.querySelector("#comment-section");
        if (commentSection) {
          comments = commentSection.querySelectorAll("div, article, section");
          console.log(
            `Found ${comments.length} potential comment elements in comment section`
          );
        }
      }

      console.log(`Found ${comments.length} comments`);
      if (comments.length > 0) {
        this.addReplyButtons();
      } else {
        setTimeout(checkForComments, 1000);
      }
    };
    checkForComments();
  }

  setupMutationObserver() {
    const commentsContainer = document.querySelector("#comment-section");
    console.log("Comments container:", commentsContainer);
    if (!commentsContainer) return;

    const observer = new MutationObserver((mutations) => {
      let shouldUpdate = false;
      mutations.forEach((mutation) => {
        if (mutation.addedNodes.length > 0) {
          mutation.addedNodes.forEach((node) => {
            if (
              node.nodeType === 1 &&
              (node.classList?.contains("cactus-comment") ||
                node.querySelector?.(".cactus-comment"))
            ) {
              shouldUpdate = true;
            }
          });
        }
      });

      if (shouldUpdate) {
        console.log("New comments detected, adding reply buttons");
        setTimeout(() => this.addReplyButtons(), 500);
      }
    });

    observer.observe(commentsContainer, {
      childList: true,
      subtree: true,
    });
  }

  addReplyButtons() {
    const comments = document.querySelectorAll(".cactus-comment");
    console.log(`Adding reply buttons to ${comments.length} comments`);

    comments.forEach((comment, index) => {
      const timestamp = this.extractTimestamp(comment);
      console.log(`Comment ${index}: timestamp = ${timestamp}`);
      if (!timestamp || this.replyButtons.has(timestamp)) return;

      const replyButton = this.createReplyButton(timestamp);
      const replyContainer = this.createReplyContainer(timestamp);

      // Try multiple locations to insert the reply button
      let insertLocation = comment.querySelector(".cactus-comment-body");
      if (!insertLocation) {
        insertLocation = comment.querySelector(".cactus-comment-content");
      }
      if (!insertLocation) {
        insertLocation = comment.querySelector(".comment-body");
      }
      if (!insertLocation) {
        // Just append to the comment itself
        insertLocation = comment;
      }

      if (insertLocation) {
        insertLocation.appendChild(replyButton);
        insertLocation.appendChild(replyContainer);

        this.replyButtons.set(timestamp, replyButton);
        this.replyForms.set(timestamp, replyContainer);
        console.log(
          `Added reply button for comment with timestamp: ${timestamp}`
        );
      } else {
        console.log(
          "No suitable location found for reply button in comment:",
          comment
        );
      }
    });
  }

  extractTimestamp(comment) {
    // Try multiple methods to extract timestamp
    const timeElement = comment.querySelector(
      ".cactus-comment-time, .cactus-comment-timestamp, time"
    );
    if (timeElement) {
      // Check for datetime attribute first
      if (timeElement.dateTime) {
        return timeElement.dateTime;
      }
      // Check for data-timestamp attribute
      if (timeElement.dataset.timestamp) {
        return timeElement.dataset.timestamp;
      }
      // Use text content as fallback
      return timeElement.textContent.trim();
    }

    // Fallback: use comment ID or generate from content
    const commentId = comment.id || comment.dataset.id;
    if (commentId) return commentId;

    // Last resort: create hash from comment content
    const content = comment.textContent.trim();
    return this.generateHash(content);
  }

  generateHash(str) {
    let hash = 0;
    for (let i = 0; i < str.length; i++) {
      const char = str.charCodeAt(i);
      hash = (hash << 5) - hash + char;
      hash = hash & hash; // Convert to 32-bit integer
    }
    return Math.abs(hash).toString();
  }

  createReplyButton(timestamp) {
    const button = document.createElement("button");
    button.className = "cactus-reply-button";
    button.textContent = "Reply";
    button.style.cssText = `
            background: #007cba;
            color: white;
            border: none;
            padding: 4px 8px;
            border-radius: 3px;
            font-size: 12px;
            cursor: pointer;
            margin-top: 8px;
            margin-right: 8px;
            transition: background-color 0.2s;
        `;

    button.addEventListener("mouseenter", () => {
      button.style.backgroundColor = "#005a87";
    });

    button.addEventListener("mouseleave", () => {
      button.style.backgroundColor = "#007cba";
    });

    button.addEventListener("click", (e) => {
      e.preventDefault();
      this.toggleReplyForm(timestamp);
    });

    return button;
  }

  createReplyContainer(timestamp) {
    const container = document.createElement("div");
    container.className = "cactus-reply-container";
    container.style.cssText = `
            display: none;
            margin-top: 10px;
            padding: 10px;
            background: #f8f9fa;
            border-radius: 5px;
            border-left: 3px solid #007cba;
        `;

    const form = this.createReplyForm(timestamp);
    container.appendChild(form);

    return container;
  }

  createReplyForm(timestamp) {
    const form = document.createElement("form");
    form.className = "cactus-reply-form";
    form.autocomplete = "on"; // Make sure autocomplete is enabled

    const replyInfo = document.createElement("div");
    replyInfo.style.cssText = `
            font-size: 12px;
            color: #666;
            margin-bottom: 8px;
            font-style: italic;
        `;
    replyInfo.textContent = `Replying to comment from ${this.formatTimestamp(
      timestamp
    )}`;

    const textarea = document.createElement("textarea");
    textarea.className = "cactus-reply-textarea";
    textarea.placeholder = "Write your reply...";
    textarea.style.cssText = `
            width: 100%;
            min-height: 80px;
            padding: 8px;
            border: 1px solid #ddd;
            border-radius: 3px;
            font-family: inherit;
            font-size: 14px;
            resize: vertical;
            box-sizing: border-box;
        `;

    // New name input field
    const nameInput = document.createElement("input");
    nameInput.type = "text";
    nameInput.className = "cactus-reply-name";
    nameInput.placeholder = "Name";
    nameInput.value = "Anonymous"; // Set the default value
    nameInput.style.cssText = `
            width: 100%;
            padding: 8px;
            border: 1px solid #ddd;
            border-radius: 3px;
            font-family: inherit;
            font-size: 14px;
            box-sizing: border-box;
            margin-bottom: 8px;
        `;
    
    const buttonContainer = document.createElement("div");
    buttonContainer.style.cssText = `
            margin-top: 8px;
            display: flex;
            gap: 8px;
        `;

    const submitButton = document.createElement("button");
    submitButton.type = "submit";
    submitButton.textContent = "Post Reply";
    submitButton.style.cssText = `
            background: hsl(0, 0%, 25.1%);
            color: white;
            border: none;
            padding: 6px 12px;
            border-radius: 3px;
            font-size: 12px;
            cursor: pointer;
            transition: background-color 0.2s;
        `;

    const cancelButton = document.createElement("button");
    cancelButton.type = "button";
    cancelButton.textContent = "Cancel";
    cancelButton.style.cssText = `
            background: hsl(0, 0%, 25.1%);
            color: white;
            border: none;
            padding: 6px 12px;
            border-radius: 3px;
            font-size: 12px;
            cursor: pointer;
            transition: background-color 0.2s;
        `;

    // Event listeners
    submitButton.addEventListener("mouseenter", () => {
      submitButton.style.backgroundColor = "hsl(0, 0%, 25.1%)";
    });

    submitButton.addEventListener("mouseleave", () => {
      submitButton.style.backgroundColor = "hsl(0, 0%, 25.1%)";
    });

    cancelButton.addEventListener("mouseenter", () => {
      cancelButton.style.backgroundColor = "hsl(0, 0%, 25.1%)";
    });

    cancelButton.addEventListener("mouseleave", () => {
      cancelButton.style.backgroundColor = "hsl(0, 0%, 25.1%)";
    });

    cancelButton.addEventListener("click", () => {
      this.hideReplyForm(timestamp);
    });

    form.addEventListener("submit", (e) => {
      e.preventDefault();
      const username = nameInput.value.trim() || "Anonymous";
      this.handleReplySubmit(timestamp, textarea.value, username);
    });

    buttonContainer.appendChild(submitButton);
    buttonContainer.appendChild(cancelButton);

    form.appendChild(replyInfo);
    form.appendChild(nameInput); // Add the new name input field
    form.appendChild(textarea);
    form.appendChild(buttonContainer);

    return form;
  }

  toggleReplyForm(timestamp) {
    const container = this.replyForms.get(timestamp);
    if (!container) return;

    const isVisible = container.style.display !== "none";

    if (isVisible) {
      this.hideReplyForm(timestamp);
    } else {
      this.showReplyForm(timestamp);
    }
  }

  showReplyForm(timestamp) {
    // Hide all other reply forms first
    this.replyForms.forEach((container, ts) => {
      if (ts !== timestamp) {
        container.style.display = "none";
      }
    });

    const container = this.replyForms.get(timestamp);
    if (container) {
      container.style.display = "block";
      const textarea = container.querySelector(".cactus-reply-textarea");
      if (textarea) {
        textarea.focus();
      }
    }
  }

  hideReplyForm(timestamp) {
    const container = this.replyForms.get(timestamp);
    if (container) {
      container.style.display = "none";
      const textarea = container.querySelector(".cactus-reply-textarea");
      if (textarea) {
        textarea.value = "";
      }
    }
  }

  handleReplySubmit(timestamp, replyText, username) {
    if (!replyText.trim()) {
      alert("Please enter a reply message.");
      return;
    }

    // Store the reply data for post-submission processing
    this.pendingReply = {
      timestamp,
      replyText,
      username,
      originalComment: this.getOriginalCommentText(timestamp)
    };

    // Insert just the reply text into the main comment textarea
    this.insertReplyIntoMainForm(replyText);

    // Hide the reply form
    this.hideReplyForm(timestamp);

    // Scroll to main comment form
    this.scrollToMainForm();

    // Set up post-submission monitoring
    this.setupSubmissionInterceptor();
  }

  getOriginalCommentText(timestamp) {
    console.log("Getting original comment text for timestamp:", timestamp);
    
    // Find the original comment by timestamp and extract its text
    const comments = document.querySelectorAll(".cactus-comment");
    for (let comment of comments) {
      const commentTimestamp = this.extractTimestamp(comment);
      console.log("Checking comment with timestamp:", commentTimestamp);
      
      if (commentTimestamp === timestamp) {
        console.log("Found matching comment");
        
        const messageElement = comment.querySelector(".cactus-message-text");
        if (messageElement) {
          // Get only the paragraph text, not the reply button
          const paragraphs = messageElement.querySelectorAll("p");
          let text = "";
          
          if (paragraphs.length > 0) {
            text = Array.from(paragraphs).map(p => p.textContent.trim()).join(" ");
          } else {
            text = messageElement.textContent.trim();
          }
          
          // Remove any existing reply buttons from the text
          text = text.replace(/Reply$/, '').trim();
          
          console.log("Extracted original comment text:", text);
          return text.substring(0, 100); // First 100 characters
        }
      }
    }
    
    console.log("Could not find original comment text");
    return "";
  }

  monitorForNewComment() {
    if (!this.pendingReply) return;

    console.log("Starting to monitor for new comment with pending reply:", this.pendingReply);
    
    // Store the current comment count
    const initialCommentCount = document.querySelectorAll(".cactus-comment").length;
    console.log("Initial comment count:", initialCommentCount);

    const checkForNewComment = () => {
      const comments = document.querySelectorAll(".cactus-comment");
      console.log("Current comment count:", comments.length);
      
      // Check if we have a new comment
      if (comments.length > initialCommentCount) {
        const lastComment = comments[comments.length - 1];
        console.log("Found new comment, checking if it's ours");
        
        if (!lastComment.dataset.replyProcessed) {
          // Check if this is likely our new comment by looking for our reply text
          const messageElement = lastComment.querySelector(".cactus-message-text, .cactus-comment-body");
          if (messageElement) {
            const commentText = messageElement.textContent.trim();
            console.log("New comment text:", commentText);
            console.log("Looking for reply text:", this.pendingReply.replyText);
            
            if (commentText.includes(this.pendingReply.replyText.trim())) {
              console.log("Found matching comment, adding quote");
              this.addQuoteToComment(lastComment);
              lastComment.dataset.replyProcessed = "true";
              this.pendingReply = null;
              return;
            }
          }
        }
      }
      
      // Keep checking for up to 10 seconds
      if (this.monitorAttempts < 20) {
        this.monitorAttempts++;
        setTimeout(checkForNewComment, 500);
      } else {
        console.log("Monitoring timeout reached, clearing pending reply");
        this.pendingReply = null;
        this.monitorAttempts = 0;
      }
    };

    this.monitorAttempts = 0;
    setTimeout(checkForNewComment, 1000); // Wait a bit for the comment to appear
  }

  addQuoteToComment(commentElement) {
    if (!this.pendingReply) {
      console.log("No pending reply to add quote for");
      return;
    }

    console.log("Adding quote to comment with pending reply:", this.pendingReply);

    const messageElement = commentElement.querySelector(".cactus-message-text");
    if (!messageElement) {
      console.log("Could not find message element in comment");
      return;
    }

    console.log("Found message element:", messageElement);

    // Create the quote element
    const quoteElement = document.createElement("div");
    quoteElement.className = "reply-quote";
    quoteElement.style.cssText = `
      font-style: italic;
      color: #666;
      border-left: 3px solid #007cba;
      margin-bottom: 10px;
      font-size: 0.9em;
      background: rgba(0, 124, 186, 0.05);
      padding: 8px 10px;
      border-radius: 3px;
    `;
    
    const truncatedText = this.pendingReply.originalComment.length > 100 
      ? this.pendingReply.originalComment.substring(0, 100) + "..."
      : this.pendingReply.originalComment;
    
    quoteElement.textContent = `"${truncatedText}"`;
    
    console.log("Created quote element with text:", quoteElement.textContent);

    // Insert the quote at the beginning of the message element
    if (messageElement.firstChild) {
      messageElement.insertBefore(quoteElement, messageElement.firstChild);
    } else {
      messageElement.appendChild(quoteElement);
    }
    
    console.log("Quote element inserted into message");
  }

  setupSubmissionInterceptor() {
    if (!this.pendingReply) return;
    
    console.log("Setting up submission interceptor");
    
    // Find the send button and add a click listener
    const sendButton = document.querySelector(".cactus-send-button");
    if (sendButton && !sendButton.dataset.intercepted) {
      sendButton.dataset.intercepted = "true";
      
      const originalClickHandler = sendButton.onclick;
      sendButton.onclick = (e) => {
        console.log("Send button clicked, pending reply exists");
        
        // Call original handler if it exists
        if (originalClickHandler) {
          originalClickHandler.call(sendButton, e);
        }
        
        // Start monitoring after a short delay
        setTimeout(() => {
          this.monitorForNewComment();
        }, 500);
      };
      
      console.log("Submission interceptor set up");
    }
  }

  formatTimestamp(timestamp) {
    // Try to parse and format the timestamp nicely
    try {
      const date = new Date(timestamp);
      if (!isNaN(date.getTime())) {
        return date.toLocaleString();
      }
    } catch (e) {
      // If parsing fails, return the original timestamp
    }

    // If it's a hash or unrecognizable format, just return a shortened version
    if (timestamp.length > 20) {
      return timestamp.substring(0, 20) + "...";
    }

    return timestamp;
  }

  insertReplyIntoMainForm(replyText) {
    // Look for the main cactus comment textarea
    const mainTextarea = document.querySelector(
      '.cactus-editor-textarea, .cactus-textarea, textarea[placeholder*="comment"], textarea[placeholder*="Comment"]'
    );

    if (mainTextarea) {
      const currentValue = mainTextarea.value;
      const newValue = currentValue
        ? `${currentValue}\n\n${replyText}`
        : replyText;
      mainTextarea.value = newValue;

      // Trigger input event to notify cactus of the change
      const event = new Event("input", { bubbles: true });
      mainTextarea.dispatchEvent(event);
    } else {
      // Fallback: copy to clipboard
      this.copyToClipboard(replyText);
      alert("Reply copied to clipboard. Please paste it into the comment box.");
    }
  }

  scrollToMainForm() {
    const mainForm = document.querySelector(
      ".cactus-editor, .cactus-comment-form, form"
    );
    if (mainForm) {
      mainForm.scrollIntoView({ behavior: "smooth", block: "center" });
    }
  }

  copyToClipboard(text) {
    if (navigator.clipboard && window.isSecureContext) {
      navigator.clipboard.writeText(text);
    } else {
      // Fallback for older browsers
      const textArea = document.createElement("textarea");
      textArea.value = text;
      textArea.style.position = "fixed";
      textArea.style.left = "-999999px";
      textArea.style.top = "-999999px";
      document.body.appendChild(textArea);
      textArea.focus();
      textArea.select();

      try {
        document.execCommand("copy");
      } catch (err) {
        console.error("Failed to copy text: ", err);
      }

      document.body.removeChild(textArea);
    }
  }
}

// Initialize the reply system when DOM is ready
console.log("Setting up reply system initialization");
if (document.readyState === "loading") {
  document.addEventListener("DOMContentLoaded", () => {
    console.log("DOM loaded, creating CactusReplySystem");
    new CactusReplySystem();
  });
} else {
  console.log("DOM already loaded, creating CactusReplySystem immediately");
  new CactusReplySystem();
}

// Also initialize after a delay to catch late-loading comments
setTimeout(() => {
  if (!window.cactusReplySystem) {
    console.log("Creating delayed CactusReplySystem instance");
    window.cactusReplySystem = new CactusReplySystem();
  }
}, 3000);
} // End of disabled code block