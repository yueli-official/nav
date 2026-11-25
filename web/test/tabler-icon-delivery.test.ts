import { readFileSync } from "node:fs";
import { expect, test, vi } from "vitest";

const contractUrl = new URL(
  "../../api/contracts/iconcontract/category-tabler.v1.json",
  import.meta.url,
);

test("every persisted category icon is bundled by the local icon provider", async () => {
  const contract = JSON.parse(readFileSync(contractUrl, "utf8")) as {
    icons: string[];
  };
  const allowed = new Set(contract.icons);
  const migration = readFileSync(
    new URL("../../api/manifest/sql/migrations/0001_init.up.sql", import.meta.url),
    "utf8",
  );
  const picker = readFileSync(
    new URL("../app/components/ManageIconPicker.vue", import.meta.url),
    "utf8",
  );
  vi.stubGlobal("defineNuxtConfig", (config: unknown) => config);
  const config = await import("../nuxt.config").then(({ default: config }) => config);
  vi.unstubAllGlobals();

  const persisted = migration.match(/i-tabler-[a-z0-9-]+/g) ?? [];
  const selectable = picker.match(/i-tabler-[a-z0-9-]+/g) ?? [];
  for (const icon of [...persisted, ...selectable]) {
    expect(allowed, `${icon} must be bundled and accepted by the API`).toContain(
      icon,
    );
  }

  expect(config.icon?.clientBundle).toMatchObject({
    icons: expect.arrayContaining(contract.icons.map((icon) => icon.replace(/^i-tabler-/, "tabler:"))),
  });
  expect(config.icon).toMatchObject({ provider: "none", fallbackToApi: false });
});
