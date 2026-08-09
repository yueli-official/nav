import { describe, expect, test } from "vitest";
import type { NavMeView } from "../app/types/navigation";
import {
  hasNavManagementAccess,
  resolveNavAccessStatus,
} from "../app/utils/nav-access";

const me = (capabilities: string[]): NavMeView => ({
  sub: "TestA123",
  authenticated: true,
  isAdministrator: false,
  capabilities,
});

describe("Nav access presentation", () => {
  test("recognizes instance-local management capabilities", () => {
    expect(hasNavManagementAccess(me(["nav.link.update"]))).toBe(true);
    expect(hasNavManagementAccess(me(["authorization.apply"]))).toBe(false);
    expect(hasNavManagementAccess(null)).toBe(false);
  });

  test("only presents an application when a role is requestable", () => {
    expect(resolveNavAccessStatus(me(["nav.link.submit"]), 0)).toBe(
      "ready_manage",
    );
    expect(resolveNavAccessStatus(me([]), 1)).toBe("ready_requestable");
    expect(resolveNavAccessStatus(me([]), 0)).toBe("ready_no_access");
  });
});
