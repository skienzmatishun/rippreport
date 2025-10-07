function loadGifs() {
  const gifImages = document.querySelectorAll('.thumbnail__animated[data-src], .widget__animated[data-src]');
  
  gifImages.forEach(img => {
    img.src = img.dataset.src;
    img.removeAttribute('data-src');
    
    img.onload = function() {
      img.classList.add('loaded');
    };
  });
}

// Wait for everything to finish loading before loading GIFs
window.addEventListener('load', function() {
  // Add a small delay to ensure all resources are truly loaded
  setTimeout(loadGifs, 100);
});