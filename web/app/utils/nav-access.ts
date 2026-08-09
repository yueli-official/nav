import type { NavMeView } from "../types/navigation";

export const navManagementCapabilities = [
  "nav.link.submit",
  "nav.link.update",
  "nav.link.moderate",
  "nav.structure.manage",
  "nav.health_check.run",
  "nav.settings.manage",
] as const;

export type NavResolvedAccessStatus =
  | "ready_manage"
  | "ready_requestable"
  | "ready_no_access";

export function hasNavManagementAccess(me: NavMeView | null): boolean {
  if (!me) return false;
  return navManagementCapabilities.some((capability) =>
    me.capabilities.includes(capability),
  );
}

export function resolveNavAccessStatus(
  me: NavMeView,
  requestableRoleCount: number,
): NavResolvedAccessStatus {
  if (hasNavManagementAccess(me)) return "ready_manage";
  return requestableRoleCount > 0 ? "ready_requestable" : "ready_no_access";
}
