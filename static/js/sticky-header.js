document.addEventListener('DOMContentLoaded', function() {
    // Only run on mobile devices
    if (window.innerWidth <= 768) {
        const header = document.querySelector('.header');
        const siteTitle = document.querySelector('.site-title');
        const postTitle = document.querySelector('.post-title');

        // Get the current post title if we're not on the homepage
        const currentPath = window.location.pathname;
        if (currentPath !== '/' && document.querySelector('h1')) {
            postTitle.textContent = document.querySelector('h1').textContent;
        }

        let lastScroll = 0;
        const scrollThreshold = 100;

        window.addEventListener('scroll', () => {
            const currentScroll = window.pageYOffset;

            if (currentPath !== '/') {
                if (currentScroll > scrollThreshold) {
                    siteTitle.style.display = 'none';
                    postTitle.style.display = 'block';
                } else {
                    siteTitle.style.display = 'flex';
                    postTitle.style.display = 'none';
                }
            }

            // Hide header when scrolling down, show when scrolling up
            if (currentScroll > lastScroll && currentScroll > 100) {
                header.style.transform = 'translateY(-100%)';
            } else {
                header.style.transform = 'translateY(0)';
            }

            lastScroll = currentScroll;
        });
    }
});
