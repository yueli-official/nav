import type { NavigationResponse } from "../../app/types/navigation";
import { decodeNavApiResponse } from "../../app/utils/apiCompat";

export default defineEventHandler(async (event): Promise<NavigationResponse> => {
  const config = useRuntimeConfig(event);
  const response: unknown = await $fetch<unknown>(
    `${config.apiBase}/api/v1/nav/catalog`,
  );
  try {
    return decodeNavApiResponse<NavigationResponse>(response);
  } catch {
    throw createError({
      statusCode: 502,
      statusMessage: "Navigation API request failed",
    });
  }
});
