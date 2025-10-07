async function handleRequest(request) {
  const base = "https://rippreport.com";
  const statusCode = 301;
  const url = new URL(request.url);
  const { pathname, search } = url;

  // Proxy for Rumble assets
  if (pathname.startsWith("/proxy/rumble/")) {
    const rumbleAssetPath = pathname.replace("/proxy/rumble/", "");
    const rumbleURL = `https://rumble.com/${rumbleAssetPath}`;
    return fetch(rumbleURL);
  }

  // Handle robots.txt explicitly
  if (pathname === "/robots.txt") {
    return fetch(request);
  }

  // Redirection Logic
  if (pathname.includes("/categories/") || pathname.includes("/tags/")) {
    return fetch(request);
  }
  if (pathname.match(/^\/\d{4}\/\d{2}\/.+?\.(jpe?g|png|gif|pdf)$/i)) {
    const destinationURL = `https://storage.googleapis.com/stateless-rippreport-com${pathname}${search}`;
    return Response.redirect(destinationURL, statusCode);
  }
  if (pathname.match(/^\/\d{4}\/\d{2}\/(.+)/i)) {
    const destinationURL = `${base}/p/${pathname.substring(12)}${search}`;
    return Response.redirect(destinationURL, statusCode);
  }
  if (
    pathname !== "/" &&
    !pathname.startsWith("/p/") &&
    !pathname.match(/\.(jpe?g|png|gif|svg|css|js|txt|ico)$/i)
  ) {
    const destinationURL = `${base}/p${pathname}${search}`;
    return Response.redirect(destinationURL, statusCode);
  }

  // Fetch Response and Mobile Customization
  const response = await fetch(request);
  const userAgent = request.headers.get("User-Agent") || "";

  // Check if mobile device
  const isMobile =
    /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(
      userAgent
    );
  if (!isMobile) {
    return response;
  }

  // Only transform text/html responses
  const contentType = response.headers.get("Content-Type") || "";
  if (!contentType.includes("text/html")) {
    return response;
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
          top: -12px;
        }
        .logo__img svg {
          transform: scale(40%) !important;
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
    </style>
  `;

  // JavaScript to inject
  const mobileScript = `
  <script>
    document.addEventListener('DOMContentLoaded', function() {
      if (window.innerWidth > 768) return;

      const logo = document.querySelector('.logo__link');
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
          postTitle.textContent = 
document.querySelector('h1').textContent;
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
    });
  </script>
`;

  // Inject styles and scripts
  html = html
    .replace("</head>", `${mobileStyles}</head>`)
    .replace("</body>", `${mobileScript}</body>`);

  return new Response(html, {
    headers: {
      "Content-Type": "text/html",
      "Cache-Control": "public, max-age=3600",
      "Content-Security-Policy":
        response.headers.get("Content-Security-Policy") || "",
    },
  });
}

addEventListener("fetch", (event) => {
  event.respondWith(handleRequest(event.request));
});
