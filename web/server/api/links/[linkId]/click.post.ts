export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig(event);
  const linkId = getRouterParam(event, "linkId");
  return $fetch(
    `${config.apiBase}/api/v1/nav/links/${encodeURIComponent(linkId || "")}/click`,
    { method: "POST" },
  );
});
