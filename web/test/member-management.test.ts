import { describe, expect, it } from "vitest";
import { readFileSync } from "node:fs";
import { fileURLToPath } from "node:url";

const page = readFileSync(
  fileURLToPath(new URL("../app/pages/manage/members.vue", import.meta.url)),
  "utf8",
);
const layout = readFileSync(
  fileURLToPath(new URL("../app/layouts/manage.vue", import.meta.url)),
  "utf8",
);
const authorization = readFileSync(
  fileURLToPath(
    new URL("../app/pages/manage/authorization.vue", import.meta.url),
  ),
  "utf8",
);

describe("Nav member management surface", () => {
  it("uses the stable collection topology and keeps only rows scrollable", () => {
    expect(page).toContain("<CollectionPanel");
    expect(page).toContain("[&>[aria-live=polite]]:overflow-y-auto");
    expect(page).toContain("[&>footer]:shrink-0");
    expect(page).toContain("身份");
    expect(page).toContain("本站关系");
    expect(page).toContain("权限与活动");
  });

  it("keeps Identity, membership, and authorization copy separate", () => {
    expect(page).toContain("Identity 用户");
    expect(page).toContain("Nav 成员关系");
    expect(page).toContain("成员存在不代表拥有维护权限");
    expect(page).toContain("不修改用户中心账号");
  });

  it("gives members and policy distinct navigation and terminology", () => {
    expect(layout).toContain('label: "成员"');
    expect(layout).toContain('to: "/manage/members"');
    expect(layout).toContain('label: "权限策略"');
    expect(authorization).toContain("新导航成员自动获得内容维护者权限");
    expect(authorization).not.toContain("注册用户自动成为内容维护者");
  });

  it("makes suspension auditable and requires an operator reason", () => {
    expect(page).toContain('required\n            :error="suspensionReasonError"');
    expect(page).toContain(':disabled="!suspensionReason.trim()"');
    expect(page).toContain("selected.suspendedAt");
    expect(page).toContain("selected.suspendedBy");
    expect(page).toContain("暂停原因");
  });
});
