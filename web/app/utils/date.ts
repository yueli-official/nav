const formatter = new Intl.RelativeTimeFormat("zh-CN", { numeric: "auto" });

export function rel(iso?: string): string {
  if (!iso) return "";
  const timestamp = Date.parse(iso);
  if (!Number.isFinite(timestamp)) return "";
  const seconds = Math.round((timestamp - Date.now()) / 1_000);
  const ranges = [
    ["year", 31_536_000],
    ["month", 2_592_000],
    ["day", 86_400],
    ["hour", 3_600],
    ["minute", 60],
  ] as const;
  for (const [unit, size] of ranges) {
    if (Math.abs(seconds) >= size)
      return formatter.format(Math.round(seconds / size), unit);
  }
  return formatter.format(seconds, "second");
}
