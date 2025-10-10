function loadGifs() {
  const gifImages = document.querySelectorAll('.thumbnail__animated[data-src], .widget__animated[data-src]');
  
  gifImages.forEach(img => {
    if (img.dataset.srcset) {
      img.srcset = img.dataset.srcset;
      img.removeAttribute('data-srcset');
    }
    if (img.dataset.src) {
      img.src = img.dataset.src;
      img.removeAttribute('data-src');
    }
    img.onload = function() {
      img.classList.add('loaded');
    };
  });
}

window.addEventListener('load', function() {
  setTimeout(loadGifs, 100);
});
