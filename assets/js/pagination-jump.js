function jumpToPage() {
  const input = document.getElementById('page-jump');
  const pageNumber = parseInt(input.value);
  
  if (pageNumber && pageNumber > 0) {
    // Get the current URL and construct the new page URL
    const currentPath = window.location.pathname;
    const currentOrigin = window.location.origin;
    
    // Remove any existing /page/N from the path
    let baseUrl = currentPath.replace(/\/page\/\d+\/?$/, '');
    
    // Ensure baseUrl ends with / if it's not just /
    if (baseUrl !== '/' && !baseUrl.endsWith('/')) {
      baseUrl += '/';
    }
    
    if (pageNumber === 1) {
      // Page 1 is the base URL without /page/1
      window.location.href = currentOrigin + baseUrl;
    } else {
      // Other pages use /page/N format
      window.location.href = currentOrigin + baseUrl + 'page/' + pageNumber + '/';
    }
  }
}

// Allow Enter key to trigger the jump
document.addEventListener('DOMContentLoaded', function() {
  const input = document.getElementById('page-jump');
  if (input) {
    input.addEventListener('keypress', function(e) {
      if (e.key === 'Enter') {
        jumpToPage();
      }
    });
  }
});