import { defineConfig } from "@playwright/test";

export default defineConfig({
 testDir: ".", testMatch: "http-result.spec.ts",
 timeout: 90_000, expect: { timeout: 20_000 },
 workers: 1, reporter: "list",
 outputDir: "../../test-results/http-result",
 use: { baseURL: process.env.NAV_E2E_URL || "http://192.168.5.7:3006", headless: true, screenshot: "only-on-failure", trace: "retain-on-failure" },
});
