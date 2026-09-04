import type { NavigationResponse } from "../../app/types/navigation";

export default defineEventHandler(async (event): Promise<NavigationResponse> => {
  const config = useRuntimeConfig(event);
  const response: unknown = await $fetch<unknown>(
    `${config.apiBase}/api/v1/nav/catalog`,
  );
  try {
    return response as NavigationResponse;
  } catch {
    throw createError({
      statusCode: 502,
      statusMessage: "Navigation API request failed",
    });
  }
});
