import { expect, test } from "vitest";
import { navFailureFeedback } from "../app/utils/navFailureFeedback";

test("maps known fields, retains unknown summaries and isolates technical details", () => {
 const feedback = navFailureFeedback({
  kind: "remote", status: 400, code: "nav.invalid_input", params: {},
  violations: [
   { pointer: "/url", code: "validation.absolute_http_url", params: {} },
   { pointer: "/newField", code: "validation.required", params: {} },
  ], traceId: "nav-test", reauth: "not-attempted",
 }, "保存失败", { "/url": "url" });
 expect(feedback.fieldErrors.url).toEqual(["请输入完整的 HTTP 或 HTTPS 网址。"]);
 expect(feedback.summary).toEqual(["请填写此项。"]);
 expect(feedback.technical).toEqual({code:"nav.invalid_input",traceId:"nav-test"});
 expect(feedback.message).not.toContain("nav-test");
 expect(navFailureFeedback(new Error("secret SQL"), "保存失败").message).toBe("保存失败");
});
