// Fix YouTube embed sizing and cross-origin issues
document.addEventListener('DOMContentLoaded', function() {
    // Wrap all YouTube iframes in content with responsive containers
    const contentIframes = document.querySelectorAll('.content iframe[src*="youtube.com"]');
    
    contentIframes.forEach(function(iframe) {
        // Skip if already wrapped
        if (iframe.parentElement.classList.contains('video-container')) {
            return;
        }
        
        // Create wrapper
        const wrapper = document.createElement('div');
        wrapper.className = 'video-container';
        
        // Insert wrapper before iframe
        iframe.parentNode.insertBefore(wrapper, iframe);
        
        // Move iframe into wrapper
        wrapper.appendChild(iframe);
        
        // Remove fixed dimensions
        iframe.removeAttribute('width');
        iframe.removeAttribute('height');
        iframe.style.width = '100%';
        iframe.style.height = '100%';
    });
    
    // Handle cross-origin errors by suppressing them
    window.addEventListener('error', function(e) {
        if (e.message && e.message.includes('Blocked a frame with origin')) {
            e.preventDefault();
            return false;
        }
    });
});