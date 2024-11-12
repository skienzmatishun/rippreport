export async function onRequest(context) {
  const { request, env } = context;

  // Log the request URL for debugging
  console.log("Request URL:", request.url);

  const url = new URL(request.url);

  // Fetch the response
  const response = await fetch(request);
  const contentType = response.headers.get('Content-Type') || '';

  // Log the content type
  console.log("Content-Type:", contentType);

  // Only proceed if the content is HTML
  if (!contentType.includes('text/html')) {
    console.log("Non-HTML response, returning original response.");
    return response; // Return unchanged response if not HTML
  }

  // HTMLRewriter to remove `loading="lazy"` from the first `.thumbnail__image`
  let foundFirstThumbnail = false;

  const transformedResponse = new HTMLRewriter()
  .on("img.thumbnail__image", {
    element(img) {
      // Log each `.thumbnail__image` found
      console.log("Found .thumbnail__image element");

      // Check if this is the first match
      if (!foundFirstThumbnail) {
        console.log("Removing loading attribute from the first .thumbnail__image");
        img.removeAttribute("loading");
        foundFirstThumbnail = true; // Only modify the first match
      } else {
        console.log("Skipping additional .thumbnail__image elements");
      }
    },
  })
  .transform(response);

  console.log("Transformation complete, returning modified response.");

  return transformedResponse;
}
