import type { NavMeView } from "~/types/navigation";
import {
  hasNavManagementAccess,
  resolveNavAccessStatus,
} from "~/utils/nav-access";

interface RequestableRoleView {
  key: string;
  displayName: string;
}

export type NavAccessStatus =
  | "signed_out"
  | "resolving"
  | "ready_manage"
  | "ready_requestable"
  | "ready_no_access"
  | "suspended"
  | "error";

interface NavAccessError {
  message: string;
  correlationId?: string;
}

interface NavAccessSnapshot {
  subject: string | null;
  status: NavAccessStatus;
  me: NavMeView | null;
  requestableRoles: RequestableRoleView[];
  error: NavAccessError | null;
  requestId: number;
}

function toAccessError(failure: unknown): NavAccessError {
  const response = failure as {
    data?: { message?: string; correlationId?: string; correlation_id?: string };
  };
  return {
    message:
      response.data?.message || "本站权限状态暂时不可用，请稍后重试。",
    correlationId:
      response.data?.correlationId || response.data?.correlation_id || undefined,
  };
}

// Effective access comes from this Nav instance, never Identity role claims
// or a browser-visible operator subject list.
export function useMe() {
  const { call } = useApi();
  const { user } = useAuth();
  const access = useState<NavAccessSnapshot>("nav-access-v3", () => ({
    subject: null,
    status: "signed_out",
    me: null,
    requestableRoles: [],
    error: null,
    requestId: 0,
  }));
  const me = computed(() => access.value.me);
  const accessStatus = computed(() => access.value.status);
  const accessError = computed(() => access.value.error);
  const accessSubject = computed(() => access.value.subject);
  const requestableRoles = computed(() => access.value.requestableRoles);

  function clear() {
    access.value = {
      subject: null,
      status: "signed_out",
      me: null,
      requestableRoles: [],
      error: null,
      requestId: access.value.requestId + 1,
    };
  }

  async function refresh(): Promise<NavMeView | null> {
    const expectedSubject = user.value?.sub?.trim() || null;
    const expectedUserKey = user.value?.userKey?.trim() || expectedSubject;
    if (!expectedSubject || !expectedUserKey) {
      clear();
      return null;
    }
    const requestId = access.value.requestId + 1;
    access.value = {
      subject: expectedSubject,
      status: "resolving",
      me: null,
      requestableRoles: [],
      error: null,
      requestId,
    };
    try {
      const resolved = (await call<{ me: NavMeView }>("/api/v1/me")).me;
      if (
        access.value.requestId !== requestId ||
        user.value?.sub !== expectedSubject
      ) {
        return null;
      }
      if (
        !resolved.authenticated ||
        resolved.sub !== expectedSubject ||
        resolved.userKey !== expectedUserKey
      ) {
        throw new Error("Nav access subject does not match the active session");
      }

      let roles: RequestableRoleView[] = [];
      if (
        resolved.membership?.status !== "suspended" &&
        !hasNavManagementAccess(resolved)
      ) {
        roles = (
          await call<{ items: RequestableRoleView[] }>(
            "/api/v1/authorization/requestable-roles",
          )
        ).items;
        if (
          access.value.requestId !== requestId ||
          user.value?.sub !== expectedSubject
        ) {
          return null;
        }
      }
      access.value = {
        subject: expectedSubject,
        status:
          resolved.membership?.status === "suspended"
            ? "suspended"
            : resolveNavAccessStatus(resolved, roles.length),
        me: resolved,
        requestableRoles: roles,
        error: null,
        requestId,
      };
      return resolved;
    } catch (failure) {
      if (
        access.value.requestId !== requestId ||
        user.value?.sub !== expectedSubject
      ) {
        return null;
      }
      access.value = {
        subject: expectedSubject,
        status: "error",
        me: null,
        requestableRoles: [],
        error: toAccessError(failure),
        requestId,
      };
      return null;
    }
  }

  const isAdministrator = computed(
    () => me.value?.isAdministrator ?? false,
  );
  const can = (capability: string) =>
    me.value?.capabilities.includes(capability) ?? false;
  const canManage = computed(() => hasNavManagementAccess(me.value));

  return {
    me,
    isAdministrator,
    can,
    canManage,
    accessStatus,
    accessError,
    accessSubject,
    requestableRoles,
    refresh,
    clear,
  };
}
