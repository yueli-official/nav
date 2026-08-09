import { readFileSync } from "node:fs";
import { resolve } from "node:path";
import { expect, test } from "vitest";

test("defines a product dark-surface ladder instead of a single black plane", () => {
  const css = readFileSync(
    resolve(process.cwd(), "app/assets/css/main.css"),
    "utf8",
  );

  for (const token of [
    "--nav-dark-canvas",
    "--nav-dark-region",
    "--nav-dark-card",
    "--nav-dark-inset",
    "--nav-dark-overlay",
  ]) {
    expect(css).toContain(token);
  }
  expect(css).toContain("--yueli-surface-page: var(--nav-dark-canvas)");
  expect(css).toContain("--yueli-surface-card: var(--nav-dark-card)");
  expect(css).toContain("--ui-bg: var(--nav-dark-card)");
  expect(css).toContain("--ui-bg-elevated: var(--nav-dark-inset)");
  expect(css).toContain(".dark #manage-main");
});
