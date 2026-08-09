export default defineNuxtRouteMiddleware(async (to) => {
  if (!to.path.startsWith("/manage")) return;
  const { user, refresh, login } = useAuth();
  if (!user.value) await refresh();
  if (!user.value) return login(to.fullPath);
  if (import.meta.server) return;

  const {
    canManage,
    accessStatus,
    accessSubject,
    refresh: refreshMe,
  } = useMe();
  if (
    accessSubject.value !== user.value.sub ||
    !accessStatus.value.startsWith("ready_")
  ) {
    await refreshMe();
  }
  if (canManage.value) return;
  if (accessStatus.value === "error") {
    return showError({
      statusCode: 503,
      statusMessage: "本站权限状态暂时不可用，请重试。",
    });
  }
  if (accessStatus.value === "suspended") {
    return showError({
      statusCode: 403,
      statusMessage: "本站成员资格已暂停，无法进入控制台。",
    });
  }
  return showError({
    statusCode: 403,
    statusMessage: "当前账户没有本站控制台权限。",
  });
});
