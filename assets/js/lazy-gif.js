function loadGifs() {
  const gifImages = document.querySelectorAll('.thumbnail__animated[data-src], .widget__animated[data-src]');
  
  gifImages.forEach(img => {
    // Set srcset first if it exists
    if (img.dataset.srcset) {
      img.srcset = img.dataset.srcset;
      img.removeAttribute('data-srcset');
    }
    // Then set src
    if (img.dataset.src) {
      img.src = img.dataset.src;
      img.removeAttribute('data-src');
    }
    // Add loaded class when image loads
    img.onload = function() {
      img.classList.add('loaded');
    };
    // Also handle error case
    img.onerror = function() {
      console.warn('Failed to load animated image:', img.dataset.src || img.src);
    };
  });
}

// Load immediately if DOM is ready, otherwise wait
if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', function() {
    setTimeout(loadGifs, 100);
  });
} else {
  // DOM already loaded
  setTimeout(loadGifs, 100);
}
