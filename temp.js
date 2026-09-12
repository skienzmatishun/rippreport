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

  // CSS/js inject
  const mobileStyles = `
    <style>@media (max-width:768px){.title-logo[style*=margin-left],.title-logo[style*=margin-right]{margin:0!important}.header{position:fixed;top:0;left:0;width:100%;background-color:#000;z-index:1000;transition:transform .3s ease}.logo__link{display:flex;align-items:center;gap:4px;width:auto;font-size:25px!important}.logo__imagebox{order:-1}.logo__img{width:45px!important;height:35px!important;position:absolute;right:30px;top:0}.site-titles{display:flex;flex-direction:row!important;gap:4px!important}.post-title{display:none;white-space:nowrap;overflow:hidden;text-overflow:ellipsis;color:#fff;font-size:17px!important;width:260px!important;text-align:left;padding:4px}body{margin-top:60px}}.widget-recent-comments{display:none}</style>`;
  const mobileScript = `<script>document.addEventListener("DOMContentLoaded",function(){if(window.innerWidth>768)return;const e=document.querySelector(".logo__link");if(!e)return;let t=e.querySelector(".site-titles"),n=e.querySelector(".post-title");if(!t){t=document.createElement("div"),t.className="site-titles";e.querySelectorAll(".title-logo").forEach(e=>{t.appendChild(e.cloneNode(!0)),e.remove()}),e.appendChild(t)}if(!n){n=document.createElement("span"),n.className="post-title";"/"!==window.location.pathname&&document.querySelector("h1")&&(n.textContent=document.querySelector("h1").textContent),e.appendChild(n)}window.addEventListener("scroll",()=>{const e=window.pageYOffset;"/"!==window.location.pathname&&(e>100?(t.style.display="none",n.style.display="block"):(t.style.display="flex",n.style.display="none"))});const o=document.querySelector(".widget-recent-comments"),l=document.querySelector(".widget-recent");if(o)if(l){o.style.display="none";const e=new IntersectionObserver(t=>{t.forEach(t=>{t.isIntersecting&&(o.style.display="block",window.dispatchEvent(new CustomEvent("loadRecentComments")),e.disconnect())})},{rootMargin:"200px"});e.observe(l)}else o.style.display="block",window.dispatchEvent(new CustomEvent("loadRecentComments"))});</script>`;

  // Desktop lazy loading script
  const desktopScript = `<script>document.addEventListener("DOMContentLoaded",function(){if(console.log("[Desktop] Width:",window.innerWidth,"isMobile:",window.innerWidth<=768),window.innerWidth<=768)return;const e=document.querySelector(".widget-recent-comments"),o=document.querySelectorAll(".list__item");if(console.log("[Desktop] Widget found:",!!e,"Posts:",o.length),!e)return void console.log("[Desktop] No widget found, exiting");if(0===o.length||o.length<5)return console.log("[Desktop] Loading immediately (posts:",o.length,")"),e.style.display="block",void setTimeout(()=>{window.dispatchEvent(new CustomEvent("loadRecentComments"))},100);console.log("[Desktop] Setting up lazy load for 5th post");const t=o[4],n=new IntersectionObserver(o=>{o.forEach(o=>{o.isIntersecting&&(console.log("[Desktop] 5th post visible, loading"),e.style.display="block",window.dispatchEvent(new CustomEvent("loadRecentComments")),n.disconnect())})},{rootMargin:"300px"});n.observe(t)}); </script>`;

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
