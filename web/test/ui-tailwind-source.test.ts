import { readFileSync } from "node:fs";
import { resolve } from "node:path";
import { expect, test } from "vitest";

test("includes responsive utilities declared by the shared UI package", () => {
  const css = readFileSync(
    resolve(process.cwd(), ".nuxt/yueli-styles.css"),
    "utf8",
  );
  expect(css).toMatch(/@import "[^"]*\/tailwindcss\/index\.css";/);
  expect(css).toMatch(/@import "[^"]*\/app\/assets\/css\/main\.css";/);
  expect(css).toMatch(/@source "[^"]*\/ui\/src";/);
});
