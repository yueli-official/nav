import { readFileSync } from "node:fs";
import { fileURLToPath } from "node:url";
import { describe, expect, it } from "vitest";

const appRoot = fileURLToPath(new URL("../app/", import.meta.url));

describe("Nav consumer account contract", () => {
  it.each(["components/AppHeader.vue", "layouts/manage.vue"])(
    "delegates authenticated identity rendering in %s",
    (relativePath) => {
      const source = readFileSync(
        new URL(relativePath, `file:///${appRoot}`),
        "utf8",
      );
      expect(source).toMatch(/<Consumer(?:Manage)?AccountControl/);
      expect(source).not.toMatch(/<PlatformUserMenu(?:\s|>)/);
      expect(source).not.toMatch(/<ManageUserMenu(?:\s|>)/);
    },
  );
});
