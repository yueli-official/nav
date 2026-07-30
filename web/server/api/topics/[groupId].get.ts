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
  let response: Envelope<NavigationGroupResponse>;
  try {
    response = await $fetch<Envelope<NavigationGroupResponse>>(
      `${config.apiBase}/api/v1/nav/groups/${encodeURIComponent(groupId || "")}`,
      { query: { page: query.page, size: query.size, sort: query.sort } },
    );
  } catch (error) {
    const failure = error as { status?: number; statusCode?: number };
    if ((failure.statusCode ?? failure.status) === 404) {
      return { missing: true as const };
    }
    throw createError({ statusCode: 502, message: "导航服务暂时不可用" });
  }
  if (response.code !== "ok") {
    throw createError({ statusCode: 502, message: response.message });
  }
  return response.data;
});
