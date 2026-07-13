export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig(event);
  const linkId = getRouterParam(event, "linkId");
  return proxyRequest(
    event,
    `${config.apiBase}/api/v1/nav/links/${encodeURIComponent(linkId || "")}/favicon`,
  );
});
