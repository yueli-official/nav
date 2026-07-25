import type { NavMeView } from "~/types/navigation";

// Effective access comes from this Nav instance, never Identity role claims
// or a browser-visible operator subject list.
export function useMe() {
  const { call } = useApi();
  const me = useState<NavMeView | null>("nav-me", () => null);

  async function refresh() {
    try {
      me.value = (await call<{ me: NavMeView }>("/api/v1/me")).me;
    } catch {
      me.value = null;
    }
  }

  const isAdministrator = computed(
    () => me.value?.isAdministrator ?? false,
  );
  const can = (capability: string) =>
    me.value?.capabilities.includes(capability) ?? false;
  const canManage = computed(() =>
    [
      "nav.link.submit",
      "nav.link.update",
      "nav.link.moderate",
      "nav.structure.manage",
      "nav.health_check.run",
      "nav.settings.manage",
    ].some(can),
  );

  return { me, isAdministrator, can, canManage, refresh };
}
