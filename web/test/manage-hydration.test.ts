import { readFileSync } from "node:fs";
import { resolve } from "node:path";
import { describe, expect, test } from "vitest";

const clientManagedPages = [
  "index.vue",
  "categories.vue",
  "tags.vue",
  "checks.vue",
  "settings.vue",
];

describe("manage hydration contracts", () => {
  test.each(clientManagedPages)(
    "%s keeps client-only async state behind a hydration boundary",
    (page) => {
      const source = readFileSync(
        resolve(process.cwd(), "app/pages/manage", page),
        "utf8",
      );

      expect(source).toContain("server: false");
      expect(source).toContain("ManageClientBoundary");
      expect(source).toContain("<ManageClientBoundary");
    },
  );
});
