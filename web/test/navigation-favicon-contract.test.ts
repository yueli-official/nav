import { readFileSync } from "node:fs";
import { expect, test } from "vitest";

const source = readFileSync(
  new URL("../app/components/NavigationFavicon.vue", import.meta.url),
  "utf8",
);

test("navigation favicons use a versioned local cache or stable text fallback", () => {
  expect(source).toContain("revision?: string");
  expect(source).toContain("const faviconSrc = computed");
  expect(source).toContain("?v=${encodeURIComponent(props.revision)}");
  expect(source).toContain("data-navigation-favicon-image");
  expect(source).toContain(":loading=\"props.eager ? 'eager' : 'lazy'\"");
  expect(source).toContain('v-if="faviconSrc && !failed"');
  expect(source).toContain('v-if="!loaded"');
  expect(source).toContain('@load="loaded = true"');
  expect(source).toContain('@error="failed = true"');
  expect(source).toContain("data-navigation-favicon-fallback");
  expect(source).not.toMatch(/google|duckduckgo|faviconkit/i);
});
