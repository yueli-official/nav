import type { NavigationGroupResponse } from "../../../app/types/navigation";

interface Envelope<T> {
  code: string;
  data: T;
  message: string;
}

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig(event);
  const groupId = getRouterParam(event, "groupId");
  const query = getQuery(event);
  const response = await $fetch<Envelope<NavigationGroupResponse>>(
    `${config.apiBase}/api/v1/nav/groups/${encodeURIComponent(groupId || "")}`,
    { query: { page: query.page, size: query.size, sort: query.sort } },
  );
  if (response.code !== "ok") {
    throw createError({ statusCode: 502, statusMessage: response.message });
  }
  return response.data;
});
