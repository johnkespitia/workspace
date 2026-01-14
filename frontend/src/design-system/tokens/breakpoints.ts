// Design Tokens - Breakpoints (Responsive Design)
export const breakpoints = {
  sm: "640px", // Small devices (landscape phones)
  md: "768px", // Medium devices (tablets)
  lg: "1024px", // Large devices (desktops)
  xl: "1280px", // Extra large devices (large desktops)
  "2xl": "1536px", // 2X Extra large devices
} as const;

export type Breakpoint = keyof typeof breakpoints;

// Media queries helpers para uso en JavaScript/TypeScript
export const mediaQueries = {
  sm: `(min-width: ${breakpoints.sm})`,
  md: `(min-width: ${breakpoints.md})`,
  lg: `(min-width: ${breakpoints.lg})`,
  xl: `(min-width: ${breakpoints.xl})`,
  "2xl": `(min-width: ${breakpoints["2xl"]})`,
} as const;

// Media queries para max-width (mobile-first)
export const maxMediaQueries = {
  sm: `(max-width: ${parseInt(breakpoints.sm) - 1}px)`,
  md: `(max-width: ${parseInt(breakpoints.md) - 1}px)`,
  lg: `(max-width: ${parseInt(breakpoints.lg) - 1}px)`,
  xl: `(max-width: ${parseInt(breakpoints.xl) - 1}px)`,
  "2xl": `(max-width: ${parseInt(breakpoints["2xl"]) - 1}px)`,
} as const;

// Helper function para usar en composables
export function useBreakpoint() {
  if (typeof window === "undefined") {
    return { current: "sm" as Breakpoint };
  }

  const getCurrentBreakpoint = (): Breakpoint => {
    const width = window.innerWidth;
    if (width >= parseInt(breakpoints["2xl"])) return "2xl";
    if (width >= parseInt(breakpoints.xl)) return "xl";
    if (width >= parseInt(breakpoints.lg)) return "lg";
    if (width >= parseInt(breakpoints.md)) return "md";
    if (width >= parseInt(breakpoints.sm)) return "sm";
    return "sm";
  };

  return { current: getCurrentBreakpoint() };
}
