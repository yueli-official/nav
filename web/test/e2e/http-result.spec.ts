import { expect, test, type Browser, type BrowserContext } from "@playwright/test";

const baseURL = process.env.NAV_E2E_URL || "http://192.168.5.7:3006";
async function authenticated(browser: Browser): Promise<BrowserContext> {
 const email = process.env.NAV_E2E_EMAIL;
 const password = process.env.NAV_E2E_PASSWORD;
 if (!email || !password) throw new Error("NAV_E2E_EMAIL and NAV_E2E_PASSWORD are required");
 const context = await browser.newContext({ baseURL });
 const login = await context.request.post((process.env.NAV_E2E_ACCOUNT_URL || "http://192.168.5.7:3000")+"/api/v1/auth/login", {data:{email,password}});
 expect(login.ok()).toBeTruthy();
 const bootstrap = await context.request.get("/auth/login?return_to=/manage");
 expect(new URL(bootstrap.url()).origin).toBe(baseURL);
 return context;
}

test("public navigation renders on desktop and mobile", async ({page}) => {
 const errors: string[] = [];
 page.on("pageerror", e => errors.push(e.message));
 for (const width of [1440,390]) {
  await page.setViewportSize({width,height:900});
  const response = await page.goto("/");
  expect(response?.status()).toBe(200);
  await expect(page.locator("main")).toBeVisible();
  expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth + 1)).toBeTruthy();
 }
 expect(errors).toEqual([]);
});

test("authenticated admin refresh keeps server-rendered shell", async ({browser}) => {
 const context = await authenticated(browser);
 try {
  const page = await context.newPage();
  const errors:string[]=[];
  page.on("pageerror", e=>errors.push(e.message));
  for(let i=0;i<3;i++){
   const response = i ? await page.reload() : await page.goto("/manage");
   expect(await response!.text()).toContain("data-admin-shell");
   await expect(page.locator("[data-admin-shell]")).toBeVisible();
  }
  for(const path of ["/contribute","/manage/settings","/manage/members","/manage/checks"]){
   await page.goto(path); await expect(page.locator("main")).toBeVisible();
  }
  expect(errors).toEqual([]);
 } finally { await context.close(); }
});

test("real category/group/link lifecycle returns raw DTOs, 201 and empty 204", async ({browser}) => {
 const context = await authenticated(browser);
 let categoryID="",groupID="",linkID="";
 try {
  const category = await context.request.post("/api/v1/admin/nav/categories", {data:{title:"HTTP result test",description:"",icon:"i-tabler-folder",sortOrder:0}});
  categoryID=(await category.json()).category.id;
  expect(category.status()).toBe(201);
  const group = await context.request.post("/api/v1/admin/nav/groups", {data:{categoryId:categoryID,title:"HTTP group test",description:"",sortOrder:0}});
  groupID=(await group.json()).group.id;
  expect(group.status()).toBe(201);
  expect(group.headers().location, "BFF must preserve creation Location").toBe("/api/v1/nav/groups/"+groupID);
  const groupRead=await context.request.get(group.headers().location!);
  expect(groupRead.status()).toBe(200);
  const input={categoryId:categoryID,groupId:groupID,title:"HTTP link test",url:"https://example.com/",description:"Contract validation",kind:"tool",status:"draft",tags:[],keywords:[],sortOrder:0};
  const created=await context.request.post("/api/v1/admin/nav/links",{data:input});
  expect(created.status()).toBe(201); linkID=(await created.json()).link.id;
  const list=await context.request.get("/api/v1/admin/nav/links?q=HTTP%20link%20test");
  const body=await list.json();
  expect(body.items.some((item:{id:string})=>item.id===linkID)).toBeTruthy();
  expect(body).not.toHaveProperty("code");
  const invalid=await context.request.patch("/api/v1/admin/nav/links/"+linkID,{data:{...input,url:"file:///private"}});
  expect(invalid.status()).toBe(400);
  const failure=await invalid.json();
  expect(failure.code).toBe("nav.invalid_input");
  expect(failure.violations[0].pointer).toBe("/url");
  expect(failure).not.toHaveProperty("detail");
  const removed=await context.request.delete("/api/v1/admin/nav/links/"+linkID);
  expect(removed.status()).toBe(204); expect(await removed.text()).toBe(""); linkID="";
 } finally {
  if(linkID) await context.request.delete("/api/v1/admin/nav/links/"+linkID);
  if(groupID) await context.request.delete("/api/v1/admin/nav/groups/"+groupID);
  if(categoryID) await context.request.delete("/api/v1/admin/nav/categories/"+categoryID);
  await context.close();
 }
});

test("accepted check jobs expose a usable Location and top-level operation DTO", async ({browser}) => {
 const context=await authenticated(browser);
 try {
  const response=await context.request.post("/api/v1/admin/nav/checks/run", {data:{scope:"filtered",q:"__nav_http_result_no_matching_links__"}});
  expect(response.status()).toBe(202);
  const job=await response.json();
  expect(job.id).toEqual(expect.any(String));
  expect(["running","completed"]).toContain(job.status);
  expect(job).not.toHaveProperty("job");
  expect(response.headers().location).toBe("/api/v1/admin/nav/checks/jobs/"+job.id);
  const status=await context.request.get(response.headers().location!);
  expect(status.status()).toBe(200);
  expect((await status.json()).id).toBe(job.id);
 } finally {await context.close();}
});

test("settings keeps field violations, unknown summary and separate trace details", async ({browser}) => {
 const context=await authenticated(browser);
 try {
  const page=await context.newPage();
  await page.route("**/api/v1/admin/nav/settings",async route=>{
   if(route.request().method()!=="PUT") return route.continue();
   await route.fulfill({status:400,contentType:"application/problem+json",body:JSON.stringify({
    type:"https://errors.yueli.dev/problems/nav.invalid_input",status:400,code:"nav.invalid_input",params:{},traceId:"nav-e2e-trace",
    violations:[{pointer:"/profile/identity/name",code:"validation.required",params:{}},{pointer:"/unknown",code:"validation.invalid",params:{}}]
   })});
  });
  await page.goto("/manage/settings");
  const field=page.getByRole("textbox",{name:"导航名称"});
  await field.fill("HTTP Result E2E");
  await page.getByRole("button",{name:"保存",exact:true}).click();
  await expect(field).toHaveAttribute("aria-invalid","true");
  await expect(field).toHaveAttribute("aria-describedby",/.+/);
  await expect(page.getByText("此项内容不符合要求。",{exact:false})).toBeVisible();
  await page.getByText("技术详情",{exact:true}).click();
  await expect(page.locator("code").filter({hasText:"nav-e2e-trace"})).toBeVisible();
 } finally {await context.close();}
});
