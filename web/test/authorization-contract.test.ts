import { readFileSync } from "node:fs";
import { resolve } from "node:path";
import { describe, expect, test } from "vitest";

const read = (path: string) =>
  readFileSync(resolve(process.cwd(), path), "utf8");

describe("instance-local authorization contracts", () => {
  test("does not derive management access from browser operator lists", () => {
    expect(read("nuxt.config.ts")).not.toContain("operatorSubs");
    for (const path of [
      "app/layouts/manage.vue",
      "app/pages/manage/index.vue",
      "app/pages/manage/categories.vue",
      "app/pages/manage/tags.vue",
      "app/pages/manage/checks.vue",
      "app/pages/manage/settings.vue",
      "app/pages/manage/authorization.vue",
    ]) {
      const source = read(path);
      expect(source).not.toMatch(/\bisAdmin\b/);
      expect(source).not.toContain("ManageShell");
    }
  });

  test("uses the public admin shell and capability-specific pages", () => {
    expect(read("app/layouts/manage.vue")).toContain("<YAdminShell");
    const expectations = new Map([
      ["app/pages/manage/index.vue", "nav.link.update"],
      ["app/pages/manage/categories.vue", "nav.structure.manage"],
      ["app/pages/manage/tags.vue", "nav.structure.manage"],
      ["app/pages/manage/checks.vue", "nav.health_check.run"],
      ["app/pages/manage/settings.vue", "nav.settings.manage"],
      ["app/pages/manage/authorization.vue", "isAdministrator"],
    ]);
    for (const [path, capability] of expectations) {
      const source = read(path);
      expect(source).toContain("<YAdminPage");
      expect(source).toContain(capability);
    }
  });
});
