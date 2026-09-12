async function handleRequest(request) {
  const base = "https://rippreport.com";
  const statusCode = 301;

  const url = new URL(request.url);
  const { pathname, search } = url;

  // Never redirect the Hugo/Cactus comment cache.
  if (pathname.startsWith("/comments-cache/")) {
    return fetch(request, {
      cf: { cacheEverything: true, cacheTtl: 3600 },
    });
  }

  // Proxy for Rumble assets
  if (pathname.startsWith("/proxy/rumble/")) {
    const rumbleAssetPath = pathname.replace("/proxy/rumble/", "");
    const rumbleURL = `https://rumble.com/${rumbleAssetPath}`;
    return fetch(rumbleURL);
  }

  // Handle robots.txt explicitly
  if (pathname === "/robots.txt") {
    return fetch(request, {
      cf: { cacheEverything: true, cacheTtl: 3600 },
    });
  }

  // Do not redirect category/tag pages
  if (pathname.includes("/categories/") || pathname.includes("/tags/")) {
    return fetch(request, {
      cf: { cacheEverything: true, cacheTtl: 3600 },
    });
  }

  // Redirect legacy dated post URLs to /p/
  if (pathname.match(/^\/\d{4}\/\d{2}\/(.+)/i)) {
    const destinationURL = `${base}/p/${pathname.substring(12)}${search}`;
    return Response.redirect(destinationURL, statusCode);
  }

  // Do not rewrite static site files or special paths.
  const skipLegacyPRedirect =
    pathname === "/" ||
    pathname.startsWith("/p/") ||
    pathname.startsWith("/cdn-cgi/") ||
    pathname.match(
      /\.(jpe?g|png|gif|webp|svg|css|js|txt|ico|json|xml|map|woff2?|ttf|mp3|mp4|webm|webmanifest)$/i
    );

  // Legacy redirect: /whatever -> /p/whatever
  if (pathname !== "/" && !skipLegacyPRedirect) {
    const destinationURL = `${base}/p${pathname}${search}`;
    return Response.redirect(destinationURL, statusCode);
  }

  // Fetch original response
  const response = await fetch(request);

  const userAgent = request.headers.get("User-Agent") || "";
  const isMobile = /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(userAgent);

  // Only transform text/html responses
  const contentType = response.headers.get("Content-Type") || "";
  if (!contentType.includes("text/html")) {
    return new Response(response.body, {
      status: response.status,
      statusText: response.statusText,
      headers: {
        ...Object.fromEntries(response.headers),
        "Cache-Control": "public, max-age=3600",
      },
    });
  }

  let html = await response.text();

  // CSS to inject
  const mobileStyles = `
    <style>
      @media (max-width: 768px) {
        .title-logo[style*="margin-right"],
        .title-logo[style*="margin-left"] {
          margin: 0 !important;
        }

        .header {
          position: fixed;
          top: 0;
          left: 0;
          width: 100%;
          background-color: #000;
          z-index: 1000;
          transition: transform 0.3s ease;
        }

        .logo__link {
          display: flex;
          align-items: center;
          gap: 4px;
          width: initial;
          font-size: 25px !important;
        }

        .logo__imagebox {
          order: -1;
        }

        .logo__img {
          width: 45px !important;
          height: 35px !important;
          position: absolute;
          right: 30px;
          top: 0px;
        }

        .site-titles {
          display: flex;
          flex-direction: row !important;
          gap: 4px !important;
        }

        .post-title {
          display: none;
          white-space: nowrap;
          overflow: hidden;
          text-overflow: ellipsis;
          color: #fff;
          font-size: 17px !important;
          width: 260px !important;
          text-align: left;
          padding: 4px;
        }

        body {
          margin-top: 60px;
        }
      }
      
      /* Hide recent comments widget initially for lazy loading */
      .widget-recent-comments {
        display: none;
      }
    </style>
  `;

  // JavaScript to inject
  const mobileScript = `
    <script>
      document.addEventListener('DOMContentLoaded', function() {
        if (window.innerWidth > 768) return;

        const logo = document.querySelector('.logo__link');
        if (!logo) return;

        let siteTitles = logo.querySelector('.site-titles');
        let postTitle = logo.querySelector('.post-title');

        if (!siteTitles) {
          siteTitles = document.createElement('div');
          siteTitles.className = 'site-titles';

          const titleSpans = logo.querySelectorAll('.title-logo');
          titleSpans.forEach(span => {
            siteTitles.appendChild(span.cloneNode(true));
            span.remove();
          });

          logo.appendChild(siteTitles);
        }

        if (!postTitle) {
          postTitle = document.createElement('span');
          postTitle.className = 'post-title';

          const currentPath = window.location.pathname;
          if (currentPath !== '/' && document.querySelector('h1')) {
            postTitle.textContent = document.querySelector('h1').textContent;
          }

          logo.appendChild(postTitle);
        }

        const scrollThreshold = 100;
        window.addEventListener('scroll', () => {
          const currentScroll = window.pageYOffset;
          if (window.location.pathname !== '/') {
            if (currentScroll > scrollThreshold) {
              siteTitles.style.display = 'none';
              postTitle.style.display = 'block';
            } else {
              siteTitles.style.display = 'flex';
              postTitle.style.display = 'none';
            }
          }
        });

        // Lazy load recent comments widget on mobile
        const recentCommentsWidget = document.querySelector('.widget-recent-comments');
        const recentArticlesWidget = document.querySelector('.widget-recent');
        
        if (recentCommentsWidget) {
          if (recentArticlesWidget) {
            // Have the Recent widget to observe
            recentCommentsWidget.style.display = 'none';
            
            const observer = new IntersectionObserver((entries) => {
              entries.forEach(entry => {
                if (entry.isIntersecting) {
                  recentCommentsWidget.style.display = 'block';
                  window.dispatchEvent(new CustomEvent('loadRecentComments'));
                  observer.disconnect();
                }
              });
            }, { rootMargin: '200px' });
            
            observer.observe(recentArticlesWidget);
          } else {
            // No Recent widget (e.g., single post), load immediately
            recentCommentsWidget.style.display = 'block';
            window.dispatchEvent(new CustomEvent('loadRecentComments'));
          }
        }
      });
    </script>
  `;

  // Desktop lazy loading script
  const desktopScript = `
    <script>
      document.addEventListener('DOMContentLoaded', function() {
        console.log('[Desktop] Width:', window.innerWidth, 'isMobile:', window.innerWidth <= 768);
        if (window.innerWidth <= 768) return;

        const recentCommentsWidget = document.querySelector('.widget-recent-comments');
        const posts = document.querySelectorAll('.list__item');

        console.log('[Desktop] Widget found:', !!recentCommentsWidget, 'Posts:', posts.length);

        if (!recentCommentsWidget) {
          console.log('[Desktop] No widget found, exiting');
          return;
        }

        // Not a list page or fewer than 5 posts - load immediately
        if (posts.length === 0 || posts.length < 5) {
          console.log('[Desktop] Loading immediately (posts:', posts.length, ')');
          recentCommentsWidget.style.display = 'block';
          window.dispatchEvent(new CustomEvent('loadRecentComments'));
          return;
        }

        // List page with 5+ posts - lazy load at 5th post
        console.log('[Desktop] Setting up lazy load for 5th post');
        const fifthPost = posts[4];
        const observer = new IntersectionObserver((entries) => {
          entries.forEach(entry => {
            if (entry.isIntersecting) {
              console.log('[Desktop] 5th post visible, loading');
              recentCommentsWidget.style.display = 'block';
              window.dispatchEvent(new CustomEvent('loadRecentComments'));
              observer.disconnect();
            }
          });
        }, { rootMargin: '300px' });

        observer.observe(fifthPost);
      });
    </script>
  `;

  // Inject styles and scripts
  html = html
    .replace("</head>", `${mobileStyles}</head>`)
    .replace("</body>", `${mobileScript}${desktopScript}</body>`);

  return new Response(html, {
    status: response.status,
    statusText: response.statusText,
    headers: {
      ...Object.fromEntries(response.headers),
      "Content-Type": "text/html",
      "Cache-Control": "public, max-age=3600",
    },
  });
}

addEventListener("fetch", (event) => {
  event.respondWith(handleRequest(event.request));
});
