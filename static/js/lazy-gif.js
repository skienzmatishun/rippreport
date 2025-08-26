document.addEventListener('DOMContentLoaded', function() {
  const gifImages = document.querySelectorAll('.thumbnail__gif[data-src]');
  
  if ('IntersectionObserver' in window) {
    const imageObserver = new IntersectionObserver((entries, observer) => {
      entries.forEach(entry => {
        if (entry.isIntersecting) {
          const img = entry.target;
          img.src = img.dataset.src;
          img.removeAttribute('data-src');
          
          img.onload = function() {
            img.classList.add('loaded');
          };
          
          observer.unobserve(img);
        }
      });
    });

    gifImages.forEach(img => imageObserver.observe(img));
  } else {
    // Fallback for browsers without IntersectionObserver
    gifImages.forEach(img => {
      img.src = img.dataset.src;
      img.removeAttribute('data-src');
      img.classList.add('loaded');
    });
  }
});