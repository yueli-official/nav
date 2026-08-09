import type { NavigationGroupResponse } from "../../../app/types/navigation";
import { decodeNavApiResponse } from "../../../app/utils/apiCompat";

export default defineEventHandler(
  async (
    event,
  ): Promise<NavigationGroupResponse | { missing: true }> => {
    const config = useRuntimeConfig(event);
    const groupId = getRouterParam(event, "groupId");
    const query = getQuery(event);
    let response: unknown;
    try {
      response = await $fetch<unknown>(
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
    try {
      return decodeNavApiResponse<NavigationGroupResponse>(response);
    } catch {
      throw createError({ statusCode: 502, message: "导航服务暂时不可用" });
    }
  },
);
