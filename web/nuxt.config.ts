const siteBrand = process.env.NUXT_PUBLIC_SITE_BRAND || "月离导航";
const cookieSecure =
  process.env.NUXT_COOKIE_SECURE === undefined
    ? process.env.NODE_ENV === "production"
    : process.env.NUXT_COOKIE_SECURE === "true";

export default defineNuxtConfig({
  extends: ["@yueli/identity-nuxt"],
  modules: ["@nuxt/ui", "@yueli/ui", "@yueli/nuxt-runtime"],
  icon: {
    serverBundle: { collections: ["tabler"] },
    clientBundle: {
      scan: {
        globInclude: [
          "app/**/*.{vue,ts}",
          "node_modules/@yueli/**/*.{vue,js,mjs,ts}",
        ],
        globExclude: ["test/**", "tests/**", ".*"],
      },
      sizeLimitKb: 256,
    },
  },
  yueliRuntime: {
    defaultTarget: "platform",
    targets: {
      platform: {
        path: "/",
        ssr: {
          cookies: ["rs_session"],
          headers: ["accept-language", "user-agent"],
        },
      },
      identity: {
        path: "/identity-api",
        ssr: {
          cookies: ["rs_session"],
          headers: ["accept-language", "user-agent"],
        },
      },
    },
  },
  css: ["~/assets/css/main.css"],
  app: {
    head: {
      htmlAttrs: { lang: "zh-CN" },
      meta: [
        { property: "og:site_name", content: siteBrand },
        { property: "og:type", content: "website" },
        { name: "twitter:card", content: "summary" },
      ],
    },
  },
  fonts: {
    providers: {
      google: false,
      googleicons: false,
      bunny: false,
      fontshare: false,
      fontsource: false,
    },
  },
  nitro: {
    esbuild: {
      options: {
        exclude: /node_modules(?!.*(?:@yueli\+|@yueli[\\/]))/,
      },
    },
  },
  buildDir: process.env.NUXT_BUILD_DIR || ".nuxt",
  devServer: {
    host: "127.0.0.1",
    port: Number(process.env.NUXT_DEV_PORT || "3006"),
  },
  runtimeConfig: {
    apiBase: process.env.NUXT_API_BASE || "http://127.0.0.1:8090",
    downstreamBase: process.env.NUXT_DOWNSTREAM_BASE || "http://127.0.0.1:8090",
    identityBase:
      process.env.NUXT_IDENTITY_BASE || "http://127.0.0.1:8081",
    cookieSecure,
    authCookieSecure: cookieSecure,
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
      oidcPostLogoutRedirectUri:
        process.env.NUXT_PUBLIC_OIDC_POST_LOGOUT_REDIRECT_URI ||
        "http://localhost:3006/",
      oidcScopes:
        process.env.NUXT_PUBLIC_OIDC_SCOPES ||
        "openid profile email roles offline_access",
      accountUrl:
        process.env.NUXT_PUBLIC_ACCOUNT_URL || "http://localhost:3000",
      siteSlug: process.env.NUXT_PUBLIC_SITE_SLUG || "nav-yueli",
      siteBrand,
      siteDomain: process.env.NUXT_PUBLIC_SITE_DOMAIN || "nav.localhost",
    },
  },
  devtools: { enabled: true },
});
