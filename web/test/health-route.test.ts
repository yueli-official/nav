import { expect, test, vi } from "vitest";

test("exposes the Nitro health endpoint used by lifecycle probes", async () => {
  vi.stubGlobal("defineEventHandler", <Handler>(handler: Handler) => handler);
  const route = await import("../server/routes/healthz.get");

  expect(route.default()).toEqual({ status: "up" });
});
