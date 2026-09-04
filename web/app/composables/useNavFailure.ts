import { computed, shallowRef } from "vue";
import { navFailureFeedback } from "../utils/navFailureFeedback";

export function useNavFailure(fields?: Readonly<Record<string, string>>) {
 const feedback = shallowRef<ReturnType<typeof navFailureFeedback>>();
 const message = computed(() => feedback.value
   ? [feedback.value.message, feedback.value.recovery, ...feedback.value.summary].filter(Boolean).join(" ")
   : "");
 function capture(error: unknown, fallback: string) {
   feedback.value = navFailureFeedback(error, fallback, fields);
   return message.value;
 }
 function clear() { feedback.value = undefined; }
 return { feedback, message, capture, clear };
}
