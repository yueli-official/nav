export default defineNuxtRouteMiddleware(async (to) => {
  if (!to.path.startsWith("/manage")) return;
  const { user, refresh, login } = useAuth();
  if (!user.value) await refresh();
  if (!user.value) return login(to.fullPath);
  if (import.meta.server) return;

  const { me, canManage, refresh: refreshMe } = useMe();
  if (!me.value) await refreshMe();
  if (canManage.value) return;
  return navigateTo("/");
});
