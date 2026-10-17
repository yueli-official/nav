import { test, expect } from "@playwright/test";
import type { NavigationResponse } from "../../app/types/navigation";

test("首页搜索、分类卡片与可收藏的分类页面", async ({ page, request }) => {
  const catalog = await (await request.get("/api/navigation")).json() as NavigationResponse;
  const category = catalog.categories.find((entry) => entry.groups.some((group) => group.items.length));
  expect(category).toBeTruthy();
  if (!category) return;
  const group = category.groups.find((entry) => entry.items.length)!;
  const item = group.items[0]!;
  const categoryPath = `/categories/${encodeURIComponent(category.id)}`;
  const errors: string[] = [];
  page.on("pageerror", (error) => errors.push(error.message));
  page.on("console", (message) => {
    if (/hydration.*mismatch/i.test(message.text())) errors.push(message.text());
  });

  await page.goto("/");
  await expect(page.getByRole("heading", { name: "发现值得收藏的站点" })).toBeVisible();
  await expect(page.getByRole("navigation", { name: "导航分类" })).toHaveCount(0);
  const search = page.getByRole("textbox", { name: "搜索导航站点" });
  await search.fill(item.title);
  await expect(page.getByRole("heading", { name: "搜索结果", exact: true })).toBeVisible();
  await expect(page.getByRole("heading", { name: item.title, exact: true }).first()).toBeVisible();
  await search.fill("unmatched-站点-090923819");
  await expect(page.getByText("没有匹配的站点，试试其他名称或关键词。")).toBeVisible();
  await page.getByRole("button", { name: "清空搜索" }).click();
  await expect(page.getByRole("heading", { name: "本周值得逛" })).toBeVisible();
  const directoryLink = page.getByRole("navigation", { name: "主导航" }).getByRole("link", { name: "分类", exact: true });
  await directoryLink.focus();
  await page.keyboard.press("Enter");
  await expect(page).toHaveURL(/\/#categories$/);
  const categories = page.getByRole("region", { name: "分类", exact: true });
  await expect(categories.getByRole("link")).toHaveCount(catalog.categories.length);
  await expect(page.getByRole("heading", { name: "分类", exact: true })).toBeInViewport();
  await categories.getByRole("link").filter({ hasText: category.title }).first().click();
  await expect(page).toHaveURL(new RegExp(`${categoryPath}$`));
  await expect(page.getByRole("heading", { name: category.title, exact: true })).toBeVisible();
  await expect(page.getByRole("heading", { name: "本周值得逛" })).toHaveCount(0);
  const topics = page.getByRole("navigation", { name: "主题导航" });
  await expect(topics.getByRole("link")).toHaveCount(category.groups.length + 1);
  await expect(page.getByRole("group", { name: `筛选${category.title}中的主题` })).toHaveCount(0);
  await topics.getByRole("link", { name: group.title, exact: true }).click();
  await expect(page).toHaveURL(new RegExp(`group=${encodeURIComponent(group.id)}`));
  await page.reload();
  await expect(page.getByRole("heading", { name: group.title, exact: true })).toBeVisible();
  await page.goBack();
  await expect(page).toHaveURL(new RegExp(`${categoryPath}$`));
  await page.goForward();
  await expect(page).toHaveURL(new RegExp("group="));
  await page.goto(`/topics/${encodeURIComponent(group.id)}`);
  await page.getByRole("navigation", { name: "面包屑" }).getByRole("link", { name: category.title, exact: true }).click();
  await expect(page).toHaveURL(new RegExp(`${categoryPath}$`));
  await page.goto(`/?category=${encodeURIComponent(category.id)}`);
  await expect(page).toHaveURL(new RegExp(`${categoryPath}$`));
  await page.getByRole("navigation", { name: "面包屑" }).getByRole("link", { name: "分类", exact: true }).click();
  await expect(page).toHaveURL(/\/#categories$/);
  await page.goto("/navigation");
  await expect(page).toHaveURL(/\/#categories$/);
  expect((await page.goto("/categories/not-a-real-category"))?.status()).toBe(404);
  expect(errors).toEqual([]);
});

test("公开页面在窄屏、桌面与亮暗主题下可操作", async ({ page, request }, testInfo) => {
  const catalog = await (await request.get("/api/navigation")).json() as NavigationResponse;
  for (const width of [390, 640, 1024, 1440]) {
    await page.setViewportSize({ width, height: 960 });
    for (const theme of ["light", "dark"]) {
      await page.addInitScript((value) => localStorage.setItem("nuxt-color-mode", value), theme);
      for (const path of ["/", "/#categories", `/categories/${catalog.categories[0]!.id}`]) {
        await page.goto(path);
        await expect(page.locator("main h1")).toBeVisible();
        await expect(page.locator("html")).toHaveClass(new RegExp(`(^| )${theme}( |$)`));
        expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth + 1)).toBe(true);
        if (path === "/#categories") {
          await expect(page.getByRole("heading", { name: "分类", exact: true })).toBeInViewport();
        }
        const active = page.getByRole("navigation", { name: "主导航" }).locator('[aria-current="page"]');
        await expect(active).toHaveCount(1);
        await active.focus();
        await expect(active).toBeFocused();
        if (path.startsWith("/categories/") && width === 390 && theme === "light") {
          const group = catalog.categories[0]!.groups[0]!;
          const selector = page.getByRole("button", { name: "选择主题", exact: true });
          await selector.click();
          await page.getByRole("option", { name: group.title, exact: true }).click();
          await expect(page).toHaveURL(new RegExp(`group=${encodeURIComponent(group.id)}`));
          await expect(page.locator("main section h2")).toHaveCount(1);
          await selector.click();
          await page.getByRole("option", { name: "全部主题", exact: true }).click();
          await expect(page).not.toHaveURL(/group=/);
          await expect(page.getByRole("option", { name: "全部主题", exact: true })).toBeHidden();
        }
        if (width === 390 || width === 1440) await page.screenshot({ path: testInfo.outputPath(`${width}-${theme}-${path.replaceAll("/", "_") || "home"}.png`), fullPage: path === "/" });
      }
    }
  }
});

test("视觉增强保留焦点、聚光和减少动态效果", async ({ page }, testInfo) => {
  await page.goto("/");
  await expect(page.locator("html")).toHaveClass(/(^| )dark( |$)/);
  const input = page.getByRole("textbox", { name: "搜索导航站点" });
  await input.focus();
  await expect(input).toBeFocused();
  const card = page.locator(".nav-category-card").first();
  await card.scrollIntoViewIfNeeded();
  await card.hover();
  await expect.poll(() => card.evaluate((element) => element.style.getPropertyValue("--pointer-x"))).not.toBe("");
  await card.focus();
  await expect(card).toBeFocused();
  await page.screenshot({ path: testInfo.outputPath("dark-category-focus.png") });
  await page.emulateMedia({ reducedMotion: "reduce" });
  await expect.poll(() => card.evaluate((element) => getComputedStyle(element).transform)).toBe("none");
  await expect.poll(() => card.evaluate((element) => getComputedStyle(element, "::before").opacity)).toBe("0");
  await page.getByRole("navigation", { name: "主导航" }).getByRole("link", { name: "首页", exact: true }).click();
  await expect(page.getByRole("heading", { name: "发现值得收藏的站点" })).toBeVisible();
  expect(await page.locator(".nav-discovery").evaluate((element) => getComputedStyle(element).animationName)).toBe("none");
  await card.locator("h3").evaluate((element) => {
    element.textContent = "设计灵感、创意工具与值得反复收藏的精选资源站点".repeat(3);
  });
  for (const width of [390, 1440]) {
    await page.setViewportSize({ width, height: 960 });
    expect(await card.evaluate((element) => element.scrollWidth <= element.clientWidth + 1)).toBe(true);
    expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth + 1)).toBe(true);
  }
});
