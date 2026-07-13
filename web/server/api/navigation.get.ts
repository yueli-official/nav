import type { NavigationResponse } from "../../app/types/navigation";

interface Envelope<T> {
  code: string;
  data: T;
  message: string;
  traceId: string;
}

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig(event);
  const response = await $fetch<Envelope<NavigationResponse>>(
    `${config.apiBase}/api/v1/nav/catalog`,
  );
  if (response.code !== "ok") {
    throw createError({
      statusCode: 502,
      statusMessage: response.message || "Navigation API request failed",
    });
  }
  return response.data;
});
