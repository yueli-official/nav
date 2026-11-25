import { test, expect } from "@playwright/test";
import type { NavigationResponse } from "../../app/types/navigation";

test("网站卡片提供站点主页和快速入口，搜索直达站群功能", async ({ page }) => {
  await page.goto("/");
  const browse = page.getByRole("region", { name: "月离站群" });
  const software = browse.getByRole("article", { name: "月离软件" });
  await expect(software.getByRole("link", { name: /月离软件/ })).toHaveAttribute("href", "https://apps.yuelili.com/");
  await expect(software.getByRole("link", { name: "软件目录" })).toHaveAttribute("href", "https://apps.yuelili.com/software");
  const docs = browse.getByRole("article", { name: "月离文档" });
  await expect(docs.getByRole("link", { name: "Premiere 脚本开发" })).toHaveAttribute("href", "https://docs.yuelili.com/pr-scripting");

  const search = page.getByRole("textbox", { name: "搜索站群入口" });
  await search.fill("ＰＤＦ");
  const pdf = page.getByRole("region", { name: "搜索结果" }).getByRole("link", { name: /PDF 工具/ });
  await expect(pdf).toHaveAttribute("href", "https://apps.yuelili.com/software/pdfcraft");
  await search.dispatchEvent("keydown", { key: "Enter", isComposing: true });
  await expect(search).toBeFocused();
  await search.press("Enter");
  await expect(pdf).toBeFocused();
  await search.fill("不存在的入口abcdef");
  await expect(page.getByText(/没有匹配的入口/)).toBeVisible();
  await page.getByRole("button", { name: "清空搜索" }).click();
  await expect(search).toBeFocused();
  await expect(docs.getByRole("link", { name: "Premiere 脚本开发" })).toBeVisible();

  await page.keyboard.press("Control+k");
  const dialog = page.getByRole("dialog", { name: "搜索站群入口" });
  await expect(dialog).toBeVisible();
  await dialog.getByRole("textbox").fill("分享代码");
  await expect(dialog.getByRole("link", { name: /分享代码和文本/ })).toHaveAttribute("href", "https://paste.yuelili.com/");
  await page.route("https://paste.yuelili.com/", (route) => route.fulfill({ body: "Paste" }));
  await dialog.getByRole("textbox").press("Enter");
  await expect(page).toHaveURL("https://paste.yuelili.com/");
});

test("既有分类和主题直达、刷新及返回入口", async ({ page, request }) => {
  const catalog = await (await request.get("/api/navigation")).json() as NavigationResponse;
  const category = catalog.categories.find((entry) => entry.groups.length > 0);
  expect(category).toBeTruthy();
  if (!category) return;
  const group = category.groups[0]!;
  const path = `/categories/${encodeURIComponent(category.id)}`;
  await page.goto(path);
  await expect(page.getByRole("heading", { name: category.title, exact: true })).toBeVisible();
  await page.getByRole("navigation", { name: "主题导航" }).getByRole("link", { name: group.title, exact: true }).click();
  await expect(page).toHaveURL(/group=/);
  await page.reload();
  await expect(page.getByRole("heading", { name: group.title, exact: true })).toBeVisible();
  await page.goBack();
  await expect(page).toHaveURL(new RegExp(`${path}$`));
  await page.goForward();
  await expect(page).toHaveURL(/group=/);
  await page.goto(`/?category=${encodeURIComponent(category.id)}`);
  await expect(page).toHaveURL(new RegExp(`${path}$`));
  await page.getByRole("navigation", { name: "面包屑" }).getByRole("link", { name: "月离站群" }).click();
  await expect(page).toHaveURL(/\/#categories$/);
  await expect(page.getByRole("region", { name: "月离站群" })).toBeVisible();
  await page.goto("/navigation");
  await expect(page).toHaveURL(/\/#categories$/);
  expect((await page.goto("/categories/not-a-real-category"))?.status()).toBe(404);
});

test("网站卡片和九宫格在窄屏、桌面与明暗主题可用", async ({ page }, testInfo) => {
  const errors: string[] = [];
  page.on("pageerror", (error) => errors.push(error.message));
  for (const width of [390, 1440]) {
    await page.setViewportSize({ width, height: 900 });
    for (const theme of ["light", "dark"]) {
      await page.goto("/");
      await expect(page.getByRole("region", { name: "月离站群" })).toBeVisible();
      const html = page.locator("html");
      if (!(await html.evaluate((element, theme) => element.classList.contains(theme), theme))) {
        await page.getByRole("button", { name: "切换颜色模式" }).click();
      }
      await expect(html).toHaveClass(new RegExp(`(^| )${theme}( |$)`));
      expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth + 1)).toBe(true);
      await page.getByRole("region", { name: "月离站群" }).locator("img").evaluateAll(async (images) => {
        await Promise.race([
          Promise.allSettled(images.map((image) => image.decode())),
          new Promise((resolve) => setTimeout(resolve, 3000)),
        ]);
      });
      await page.screenshot({ path: testInfo.outputPath(`portal-${width}-${theme}.png`), fullPage: true, animations: "disabled" });
    }
  }
  await page.getByRole("button", { name: "月离网站" }).click();
  const launcher = page.locator("[data-site-launcher]");
  await expect(launcher.getByRole("link", { name: "文档", exact: true })).toHaveAttribute("href", "https://docs.yuelili.com/");
  await page.keyboard.press("Escape");
  await expect(launcher).toBeHidden();
  expect(errors).toEqual([]);
});
