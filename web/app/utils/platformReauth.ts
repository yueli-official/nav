export interface NavReauthOptions {
  requireLoggedIn?: boolean;
}

export type NavReauthHandler = (
  options?: NavReauthOptions,
) => Promise<boolean>;

export function getOptionalNavReauth(): NavReauthHandler | undefined {
  const nuxtApp = tryUseNuxtApp() as
    | (ReturnType<typeof useNuxtApp> & {
        $platformReauth?: NavReauthHandler;
      })
    | null;

  return nuxtApp?.$platformReauth;
}
