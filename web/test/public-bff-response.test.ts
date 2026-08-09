import { expect, test, vi } from "vitest";

test("public navigation BFF accepts the current bare API response", async () => {
  const catalog = {
    version: 1,
    site: { name: "月离导航" },
    categories: [],
    stats: { categoryCount: 0, groupCount: 0, linkCount: 0 },
  };

  vi.stubGlobal("defineEventHandler", <Handler>(handler: Handler) => handler);
  vi.stubGlobal("useRuntimeConfig", () => ({ apiBase: "http://nav-api.test" }));
  vi.stubGlobal("$fetch", vi.fn().mockResolvedValue(catalog));
  vi.stubGlobal("createError", (value: unknown) => value);

  const route = await import("../server/api/navigation.get");
  const handler = route.default as (event: unknown) => Promise<unknown>;

  await expect(handler({})).resolves.toEqual(catalog);
});
