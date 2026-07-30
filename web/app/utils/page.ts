export const PAGE_WIDTHS = {
  narrow: "max-w-5xl",
  wide: "max-w-7xl",
  full: "max-w-[1440px]",
} as const;

export type PageWidth = keyof typeof PAGE_WIDTHS;
