// Function to generate a unique light color based on the username
function getColorForUsername(username) {
    let hash = 0;
    for (let i = 0; i < username.length; i++) {
        hash = username.charCodeAt(i) + ((hash << 5) - hash);
    }
    const color = Math.abs(hash % 360); // Get a hue value from 0 to 359
    return `hsl(${color}, 20%, 90%)`; // Light color closer to white
}

// Function to apply colors to comments
function applyCommentColors() {
    const comments = document.querySelectorAll('.cactus-comment');
    let allStyled = true; // Flag to check if all comments were styled

    comments.forEach(comment => {
        const displayNameLink = comment.querySelector('.cactus-comment-displayname');
        const username = displayNameLink.href.split('/').pop(); // Extract username from link

        // Check if the username exists before applying the color
        if (username) {
            // Generate and set the background color based on the username
            const backgroundColor = getColorForUsername(username);
            comment.style.backgroundColor = backgroundColor;
        } else {
            allStyled = false; // Mark as false if any comment is not styled
        }
    });

    return allStyled; // Return true if all comments have been styled
}

// Function to try applying colors with a delay
function tryApplyColorsWithDelay(attempts = 0) {
    const allStyled = applyCommentColors(); // Attempt to apply colors

    if (!allStyled && attempts < 3) { // Retry if not all are styled, limit to 3 attempts
        setTimeout(() => tryApplyColorsWithDelay(attempts + 1), 25000); // Try again after 10 seconds
    }
}

// Create a Mutation Observer to watch for changes in the comments container
const commentsContainer = document.querySelector('.cactus-comments-list');

const observer = new MutationObserver(mutations => {
    for (const mutation of mutations) {
        if (mutation.addedNodes.length) {
            // New comments have been added, apply colors
            applyCommentColors();
        }
    }
});

// Start observing the comments container for child node changes
if (commentsContainer) {
    observer.observe(commentsContainer, { childList: true, subtree: true });
}

// Execute colors on DOMContentLoaded
document.addEventListener('DOMContentLoaded', () => {
    tryApplyColorsWithDelay(); // Start the color application process
});
