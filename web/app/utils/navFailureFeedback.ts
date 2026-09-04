import { resolveFailureFeedback, type ProblemParams } from "@yueli/http-runtime";
import { navFailurePresentation, type NavFailureCode } from "../generated/navFailure";

const messages: Record<NavFailureCode, string> = {
 "nav.not_found": "请求的内容不存在或已被删除。",
 "nav.forbidden": "你没有执行此操作的权限。",
 "nav.conflict": "当前内容存在关联或已被其他操作修改，请刷新后检查。",
 "nav.not_initialized": "站点配置尚未就绪，请联系管理员。",
 "nav.authorization_unavailable": "权限服务暂时不可用，请稍后重试。",
 "nav.membership_unavailable": "成员服务暂时不可用，请稍后重试。",
 "nav.membership_suspended": "你的导航成员资格已暂停，请联系管理员。",
 "nav.site_profile_revision_conflict": "设置已更新，请刷新后重新修改。",
 "nav.site_profile_precondition_required": "请先加载最新设置再保存。",
 "nav.invalid_input": "提交内容不符合要求，请检查标出的字段。",
};
const violations: Record<string, string> = {
 required: "请填写此项。", invalid: "此项内容不符合要求。",
 absolute_http_url: "请输入完整的 HTTP 或 HTTPS 网址。",
 one_of: "请选择有效选项。", in: "请选择有效选项。",
 different: "新值不能与原值相同。", not_found: "选择的内容不存在。",
 category_mismatch: "主题不属于所选分类，请重新选择。",
 unknown: "请选择有效选项。", regex: "格式不正确。",
};
function resolveText(code: string, params: ProblemParams) {
 if (Object.hasOwn(navFailurePresentation, code)) return { message: messages[code as NavFailureCode] };
 if (code === "common.validation_failed") return { message: messages["nav.invalid_input"] };
 if (code === "common.internal") return { message: "服务暂时无法完成操作，请稍后重试。" };
 if (code === "common.rate_limited") return { message: "操作过于频繁，请稍后重试。" };
 if (code.startsWith("validation.")) {
  const rule = code.slice(11);
  if (rule === "minimum" && typeof params.min === "number") return { message: `不能小于 ${params.min}。` };
  if ((rule === "maximum" || rule === "length") && typeof params.max === "number") return { message: `不能超过 ${params.max}。` };
  return { message: violations[rule] ?? "此项内容不符合要求。" };
 }
 return undefined;
}
export function navFailureFeedback(error: unknown, fallback: string, fields?: Readonly<Record<string, string>>) {
 return resolveFailureFeedback(error, { fallback, fields, resolveText });
}
export function navFailureMessage(error: unknown, fallback: string) {
 const feedback = navFailureFeedback(error, fallback);
 return [feedback.message, feedback.recovery, ...feedback.summary].filter(Boolean).join(" ");
}
