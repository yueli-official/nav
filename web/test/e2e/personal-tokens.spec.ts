import {
  expect,
  request,
  test,
  type APIRequestContext,
  type Browser,
  type BrowserContext,
  type BrowserContextOptions,
  type Page,
} from "@playwright/test";

test.use({ trace: "off", video: "off", screenshot: "off" });

async function json(
  client: APIRequestContext,
  method: string,
  path: string,
  status: number,
  data?: unknown,
) {
  const response = await client.fetch(path, { method, data });
  expect(
    response.status(),
    `${method} ${path}: ${response.status() === status ? "" : await response.text()}`,
  ).toBe(status);
  if (status === 204) {
    expect(await response.text()).toBe("");
    return;
  }
  const body = await response.json();
  if (status >= 400) {
    expect(response.headers()["content-type"]).toContain(
      "application/problem+json",
    );
    expect(body.status).toBe(status);
    expect(body.traceId).toBeTruthy();
  }
  return body;
}

async function login(
  browser: Browser,
  account: string,
  nav: string,
  options: BrowserContextOptions,
): Promise<BrowserContext> {
  const context = await browser.newContext(options);
  const response = await context.request.post(`${account}/api/v1/auth/login`, {
    data: {
      email: process.env.NAV_E2E_EMAIL || "test@example.com",
      password:
        process.env.NAV_E2E_PASSWORD || "Yueli-local-development-2026",
    },
  });
  expect(
    response.ok(),
    `identity login failed with HTTP ${response.status()}`,
  ).toBeTruthy();
  const session = await context.request.get(`${nav}/auth/login?return_to=/`);
  expect(
    session.ok(),
    `Nav session bootstrap failed with HTTP ${session.status()}`,
  ).toBeTruthy();
  return context;
}

async function expectNoHorizontalOverflow(page: Page, label: string) {
  const widths = await page.evaluate(() => ({
    viewport: window.innerWidth,
    document: Math.max(
      document.documentElement.scrollWidth,
      document.body.scrollWidth,
    ),
  }));
  expect(
    widths.document,
    `${label}横向溢出：viewport=${widths.viewport}px document=${widths.document}px`,
  ).toBeLessThanOrEqual(widths.viewport + 1);
}

test("developer token catalog, Nav lifecycle, route boundaries, and revocation", async ({
  browser,
}, testInfo) => {
  test.setTimeout(180_000);
  const nav = process.env.NAV_E2E_URL || "http://nav.dev.yuelili.test:3006";
  const account =
    process.env.NAV_E2E_ACCOUNT_URL ||
    "http://account-nav.dev.yuelili.test:3100";
  const api = process.env.NAV_E2E_API_URL || "http://127.0.0.1:8090";
  const marker = `${Date.now().toString(36)}-${testInfo.project.name || "desktop"}`;
  const scope = (key: string) =>
    `site:${Buffer.from("nav-yueli-web").toString("base64url")}:${key}`;
  const capabilityKeys = [
    "nav.link.submit",
    "nav.link.update",
    "nav.link.moderate",
    "nav.structure.manage",
    "nav.health_check.run",
    "nav.settings.manage",
  ];
  const clients: APIRequestContext[] = [];
  let tokenID = 0;
  let token = "";
  let linkID = "";

  const desktop = await login(
    browser,
    account,
    nav,
    { viewport: { width: 1440, height: 1000 } },
  );
  const page = await desktop.newPage();

  try {
    const catalog = await json(
      desktop.request,
      "GET",
      `${account}/api/v1/pat/scopes`,
      200,
    );
    const navScopes = catalog.items.filter(
      (item: { site: string }) => item.site === "nav-yueli-web",
    );
    expect(navScopes.map((item: { key: string }) => item.key).sort()).toEqual(
      capabilityKeys.map(scope).sort(),
    );
    expect(catalog.unavailableSites).not.toContain("Nav");

    await page.goto(`${account}/developer-tokens`);
    const create = page.getByRole("button", { name: "创建令牌", exact: true });
    await expect(create).toBeVisible();
    await expect
      .poll(() =>
        create.evaluate((node) =>
          Boolean(
            (node as HTMLElement & { __vueParentComponent?: unknown })
              .__vueParentComponent,
          ),
        ),
      )
      .toBe(true);
    await create.click();
    const dialog = page.getByRole("dialog");
    await dialog
      .getByPlaceholder("例如：本地脚本")
      .fill(`Nav 完整令牌 ${marker}`);
    for (const permission of navScopes) {
      await dialog
        .getByRole("checkbox", { name: permission.label, exact: true })
        .check();
    }
    await expect(
      dialog.getByRole("checkbox", { name: "提交链接", exact: true }),
    ).toBeVisible();
    await expectNoHorizontalOverflow(page, "Nav PAT 桌面权限目录");
    await page.screenshot({
      path: testInfo.outputPath("developer-token-permissions-desktop.png"),
      fullPage: false,
    });
    const createdResponse = page.waitForResponse(
      (response) =>
        response.url().endsWith("/api/v1/pat") &&
        response.request().method() === "POST",
    );
    await dialog.getByRole("button", { name: "创建令牌", exact: true }).click();
    const response = await createdResponse;
    expect(response.status()).toBe(201);
    const created = await response.json();
    tokenID = created.id;
    token = created.token;

    const navClient = await request.newContext({
      baseURL: api,
      extraHTTPHeaders: { authorization: `Bearer ${token}` },
    });
    clients.push(navClient);

    await json(navClient, "GET", "/api/v1/nav/catalog", 200);
    await json(navClient, "GET", "/api/v1/me", 403);
    await json(navClient, "GET", "/api/v1/admin/nav/members", 403);
    await json(
      navClient,
      "GET",
      "/api/v1/authorization/manage/console",
      403,
    );
    await json(
      navClient,
      "GET",
      "/api/v1/internal/personal-token/permissions?userKey=TestA123",
      403,
    );

    const structure = await json(
      navClient,
      "GET",
      "/api/v1/admin/nav/structure",
      200,
    );
    const category = structure.categories.find(
      (item: { groups?: unknown[] }) => item.groups?.length,
    );
    if (!category) throw new Error("Nav fixture has no category with a group");
    const group = category.groups[0];
    expect(category?.id).toBeTruthy();
    expect(group?.id).toBeTruthy();
    const body = {
      categoryId: category.id,
      groupId: group.id,
      title: `PAT 链接 ${marker}`,
      url: `https://example.com/?nav-pat=${encodeURIComponent(marker)}`,
      description: "Nav PAT Playwright acceptance",
      tags: ["PAT"],
      keywords: ["developer-token"],
      kind: "tool",
      featured: false,
      status: "draft",
      sortOrder: 0,
    };
    const createdLink = await json(
      navClient,
      "POST",
      "/api/v1/admin/nav/links",
      201,
      body,
    );
    linkID = createdLink.link.id;
    const links = await json(
      navClient,
      "GET",
      `/api/v1/admin/nav/links?q=${encodeURIComponent(marker)}`,
      200,
    );
    expect(links.items.map((item: { id: string }) => item.id)).toContain(linkID);
    const updated = await json(
      navClient,
      "PATCH",
      `/api/v1/admin/nav/links/${linkID}`,
      200,
      { ...body, description: "Nav PAT updated acceptance" },
    );
    expect(updated.link.description).toBe("Nav PAT updated acceptance");
    await json(navClient, "GET", "/api/v1/admin/nav/checks?size=1", 200);
    await json(navClient, "GET", "/api/v1/admin/nav/settings", 200);

    const mobile = await login(
      browser,
      account,
      nav,
      { viewport: { width: 390, height: 844 }, isMobile: true },
    );
    try {
      const mobilePage = await mobile.newPage();
      await mobilePage.goto(`${account}/developer-tokens`);
      const mobileCreate = mobilePage.getByRole("button", {
        name: "创建令牌",
        exact: true,
      });
      await expect
        .poll(() =>
          mobileCreate.evaluate((node) =>
            Boolean(
              (node as HTMLElement & { __vueParentComponent?: unknown })
                .__vueParentComponent,
            ),
          ),
        )
        .toBe(true);
      await mobileCreate.click();
      const mobileDialog = mobilePage.getByRole("dialog");
      await expect(mobileDialog).toBeVisible();
      await mobileDialog.getByLabel("搜索分类或操作").fill("导航");
      const structurePermission = mobileDialog.getByRole("checkbox", {
        name: "管理导航结构",
        exact: true,
      });
      await expect(structurePermission).toBeVisible();
      await structurePermission.scrollIntoViewIfNeeded();
      await expectNoHorizontalOverflow(mobilePage, "Nav PAT 手机权限目录");
      await mobilePage.screenshot({
        path: testInfo.outputPath("developer-token-permissions-mobile.png"),
        fullPage: false,
      });
    } finally {
      await mobile.close();
    }

    await json(navClient, "DELETE", `/api/v1/admin/nav/links/${linkID}`, 204);
    linkID = "";
    await json(
      desktop.request,
      "DELETE",
      `${account}/api/v1/pat/${tokenID}`,
      204,
    );
    tokenID = 0;
    await json(navClient, "GET", "/api/v1/nav/catalog", 401);
  } finally {
    if (linkID && token) {
      const cleanup = await request.newContext({
        baseURL: api,
        extraHTTPHeaders: { authorization: `Bearer ${token}` },
      });
      await cleanup.delete(`/api/v1/admin/nav/links/${linkID}`);
      await cleanup.dispose();
    }
    if (tokenID) {
      await desktop.request.delete(`${account}/api/v1/pat/${tokenID}`);
    }
    for (const client of clients) await client.dispose();
    await desktop.close();
  }
});
