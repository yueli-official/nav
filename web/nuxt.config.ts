const siteBrand = process.env.NUXT_PUBLIC_SITE_BRAND || "月离导航";

export default defineNuxtConfig({
  extends: ["@platform/auth", "@platform/site", "@platform/manage"],
  modules: ["@nuxt/ui"],
  css: ["~/assets/css/main.css"],
  app: {
    head: {
      meta: [
        { property: "og:site_name", content: siteBrand },
        { property: "og:type", content: "website" },
        { name: "twitter:card", content: "summary" },
      ],
    },
  },
  buildDir: process.env.NUXT_BUILD_DIR || ".nuxt",
  devServer: { port: Number(process.env.NUXT_DEV_PORT || "3006") },
  runtimeConfig: {
    apiBase: process.env.NUXT_API_BASE || "http://127.0.0.1:8090",
    downstreamBase: process.env.NUXT_DOWNSTREAM_BASE || "http://127.0.0.1:8090",
    sealSecret:
      process.env.NUXT_SEAL_SECRET ||
      "dev-nav-seal-secret-change-me-0123456789abcdef",
    public: {
      oidcIssuer:
        process.env.NUXT_PUBLIC_OIDC_ISSUER || "http://localhost:8081",
      oidcClientId: process.env.NUXT_PUBLIC_OIDC_CLIENT_ID || "nav-yueli-web",
      oidcRedirectUri:
        process.env.NUXT_PUBLIC_OIDC_REDIRECT_URI ||
        "http://localhost:3006/auth/callback",
      oidcScopes:
        process.env.NUXT_PUBLIC_OIDC_SCOPES ||
        "openid profile email roles offline_access",
      operatorSubs:
        process.env.NUXT_PUBLIC_OPERATOR_SUBS ||
        "ac73d232-ce55-487d-bb39-fd336f1a9806",
      accountUrl:
        process.env.NUXT_PUBLIC_ACCOUNT_URL || "http://localhost:3000",
      siteSlug: process.env.NUXT_PUBLIC_SITE_SLUG || "nav-yueli",
      siteBrand,
      siteDomain: process.env.NUXT_PUBLIC_SITE_DOMAIN || "nav.localhost",
    },
  },
  devtools: { enabled: true },
});
