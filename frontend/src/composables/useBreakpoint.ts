import { ref, onMounted, onUnmounted } from "vue";
import {
  breakpoints,
  type Breakpoint,
} from "@/design-system/tokens/breakpoints";

/**
 * Composable para detectar el breakpoint actual de forma reactiva
 */
export function useBreakpoint() {
  const currentBreakpoint = ref<Breakpoint>("sm");

  const updateBreakpoint = () => {
    if (typeof window === "undefined") {
      currentBreakpoint.value = "sm";
      return;
    }

    const width = window.innerWidth;
    if (width >= parseInt(breakpoints["2xl"])) {
      currentBreakpoint.value = "2xl";
    } else if (width >= parseInt(breakpoints.xl)) {
      currentBreakpoint.value = "xl";
    } else if (width >= parseInt(breakpoints.lg)) {
      currentBreakpoint.value = "lg";
    } else if (width >= parseInt(breakpoints.md)) {
      currentBreakpoint.value = "md";
    } else if (width >= parseInt(breakpoints.sm)) {
      currentBreakpoint.value = "sm";
    } else {
      currentBreakpoint.value = "sm";
    }
  };

  onMounted(() => {
    updateBreakpoint();
    window.addEventListener("resize", updateBreakpoint);
  });

  onUnmounted(() => {
    window.removeEventListener("resize", updateBreakpoint);
  });

  const isMobile = () => currentBreakpoint.value === "sm";
  const isTablet = () => currentBreakpoint.value === "md";
  const isDesktop = () => ["lg", "xl", "2xl"].includes(currentBreakpoint.value);

  return {
    currentBreakpoint,
    isMobile,
    isTablet,
    isDesktop,
  };
}
