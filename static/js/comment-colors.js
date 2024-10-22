function getColorForUsername(username) {
    // Use a simple hash function to derive a unique color based on the username
    let hash = 0;
    for (let i = 0; i < username.length; i++) {
        hash = username.charCodeAt(i) + ((hash << 5) - hash);
    }
    // Convert hash to a color
    const color = Math.abs(hash % 360); // Get a hue value from 0 to 359
    return `hsl(${color}, 60%, 70%)`; // Light color
}

document.querySelectorAll('.cactus-comment').forEach(comment => {
    const displayNameLink = comment.querySelector('.cactus-comment-displayname');
    const username = displayNameLink.href.split('/').pop(); // Extract username from link

    // Generate and set the background color based on the username
    const backgroundColor = getColorForUsername(username);
    comment.style.backgroundColor = backgroundColor;
});

document.querySelectorAll('.reply-button').forEach(button => {
    button.addEventListener('click', function() {
        const replyContainer = this.nextElementSibling; // Get the reply container
        replyContainer.style.display = replyContainer.style.display === 'none' ? 'block' : 'none';
    });
});
