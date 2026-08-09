import { readFileSync } from "node:fs";
import { resolve } from "node:path";
import { expect, test } from "vitest";

test("includes responsive utilities declared by the shared UI package", () => {
  const css = readFileSync(
    resolve(process.cwd(), "app/assets/css/main.css"),
    "utf8",
  );
  expect(css).toContain('@import "@yueli/ui/tailwind.css";');
});
