/**
 * Cactus Comments Reply System
 * Adds reply functionality to existing cactus comments using timestamps for targeting
 */

console.log("Reply system script loaded");

class CactusReplySystem {
  constructor() {
    console.log("CactusReplySystem constructor called");
    this.replyButtons = new Map();
    this.replyForms = new Map();
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
      this.handleReplySubmit(timestamp, textarea.value);
    });

    buttonContainer.appendChild(submitButton);
    buttonContainer.appendChild(cancelButton);

    form.appendChild(replyInfo);
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

  handleReplySubmit(timestamp, replyText) {
    if (!replyText.trim()) {
      alert("Please enter a reply message.");
      return;
    }

    // Format the reply with reference to original comment
    const formattedReply = this.formatReply(timestamp, replyText);

    // Insert the reply into the main comment textarea
    this.insertReplyIntoMainForm(formattedReply);

    // Hide the reply form
    this.hideReplyForm(timestamp);

    // Scroll to main comment form
    this.scrollToMainForm();
  }

  formatReply(timestamp, replyText) {
    const timeStr = this.formatTimestamp(timestamp);
    return `@Reply to comment from ${timeStr}:\n\n${replyText}`;
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
interceptSubmit() {
  const mainForm = document.querySelector(".cactus-editor, .cactus-comment-form, form");
  if (mainForm) {
    mainForm.addEventListener("submit", (event) => {
      // Find the username input field within the form
      const usernameInput = mainForm.querySelector("input[type=text][name=username]");
      
      // If the input exists and its value is blank, set it to "Anonymous"
      if (usernameInput && !usernameInput.value.trim()) {
        usernameInput.value = "Anonymous";
      }
    });
  }
}