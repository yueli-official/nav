import { readFileSync } from "node:fs";
import { expect, test } from "vitest";

const page = readFileSync(
  new URL("../app/pages/manage/checks.vue", import.meta.url),
  "utf8",
);
const types = readFileSync(
  new URL("../app/types/navigation.ts", import.meta.url),
  "utf8",
);

test("runs link checks as a polled background job with visible progress", () => {
  expect(types).toContain("NavigationCheckJobStatus");
  expect(types).toContain('status: NavigationCheckJobStatus');
  expect(page).toContain("/api/v1/admin/nav/checks/jobs/${jobId}");
  expect(page).toContain("<UProgress");
  expect(page).toContain(":max=\"activeJob.total\"");
  expect(page).toContain("检查在后台并发进行");
});

test("does not make the result collection inert while a job runs", () => {
  expect(page).not.toContain(':inert="checking"');
  expect(page).toContain("nav-check-job-id");
});

test("offers complete result filters and a reversible manual exemption", () => {
  expect(types).toContain("healthCheckExempt?: boolean");
  expect(types).toContain("checkableTotal: number");
  expect(types).toContain("exempt: number");
  expect(page).toContain('label: `正常 · ${counts.value.healthy}`');
  expect(page).toContain('label: `失效 · ${counts.value.broken}`');
  expect(page).toContain('label: `超时 · ${counts.value.timeout}`');
  expect(page).toContain('label: `免检 · ${counts.value.exempt}`');
  expect(page).toContain("/api/v1/admin/nav/checks/${link.id}/exemption");
  expect(page).toContain("恢复检查");
  expect(page).toContain("设为免检");
});
