export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig(event);
  const linkId = getRouterParam(event, "linkId");
  return proxyRequest(
    event,
    `${config.apiBase}/api/v1/nav/links/${encodeURIComponent(linkId || "")}/favicon`,
    {
      onResponse(proxyEvent, response) {
        if (response.status !== 404) return;
        setResponseStatus(proxyEvent, 204);
        setResponseHeader(proxyEvent, "cache-control", "public, max-age=300");
      },
    },
  );
});
