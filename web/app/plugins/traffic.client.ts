import { createBrowserUUID } from "@yueli/nuxt-runtime/runtime";

const visitorStorageKey = "nav-visitor-v1";

function visitorID() {
  try {
    const saved = localStorage.getItem(visitorStorageKey);
    if (saved) return saved;
    const created = createBrowserUUID();
    localStorage.setItem(visitorStorageKey, created);
    return created;
  } catch {
    return undefined;
  }
}

export default defineNuxtPlugin((app) => {
  const router = useRouter();
  const { call } = useApi();
  const error = useError();
  let previousPath = "";

  function record() {
    const path = router.currentRoute.value.path.replace(/\/$/, "") || "/";
    if (path === previousPath) return;
    const referrer = previousPath ? "internal" : document.referrer ? new URL(document.referrer).hostname : "";
    previousPath = path;
    if (error.value || (path !== "/" && !/^\/categories\/[^/]+$/.test(path))) return;
    const privacy = navigator as Navigator & { globalPrivacyControl?: boolean };
    if (privacy.doNotTrack === "1" || privacy.globalPrivacyControl) return;
    void call("/api/v1/traffic/views", {
      method: "POST",
      body: {
        eventId: createBrowserUUID(),
        path,
        referrer,
        visitorId: visitorID(),
        occurredAt: new Date().toISOString(),
      },
    }).catch(() => undefined);
  }

  app.hook("app:mounted", record);
  app.hook("page:finish", record);
});
