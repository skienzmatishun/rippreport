export async function onRequest(context) {
  const { request } = context;

  const url = new URL(request.url);

  // Fetch the response
  const response = await fetch(request);
  const contentType = response.headers.get('Content-Type') || '';

  // Only proceed if the content is HTML
  if (!contentType.includes('text/html')) {
    return response; // Return unchanged response if not HTML
  }

  // HTMLRewriter logic to remove `loading="lazy"` from the first 
`.thumbnail__image`
  return new HTMLRewriter()
    .on("img.thumbnail__image", {
      element(img) {
        if (!this.foundFirstThumbnail) {
          img.removeAttribute("loading");
          this.foundFirstThumbnail = true; // Only modify the first match
        }
      },
    })
    .transform(response);
}
