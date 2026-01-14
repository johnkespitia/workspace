// Utilidades de accesibilidad

/**
 * Verifica si el contraste entre dos colores cumple con WCAG AA
 * @param color1 Color de fondo (hex)
 * @param color2 Color de texto (hex)
 * @returns true si cumple con WCAG AA (ratio >= 4.5 para texto normal, >= 3 para texto grande)
 */
export function checkContrast(color1: string, color2: string): boolean {
  const getLuminance = (hex: string): number => {
    const rgb = hexToRgb(hex);
    if (!rgb) return 0;

    const [r, g, b] = [rgb.r, rgb.g, rgb.b].map((val) => {
      val = val / 255;
      return val <= 0.03928
        ? val / 12.92
        : Math.pow((val + 0.055) / 1.055, 2.4);
    });

    return 0.2126 * r + 0.7152 * g + 0.0722 * b;
  };

  const hexToRgb = (
    hex: string
  ): { r: number; g: number; b: number } | null => {
    const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
    return result
      ? {
          r: parseInt(result[1], 16),
          g: parseInt(result[2], 16),
          b: parseInt(result[3], 16),
        }
      : null;
  };

  const l1 = getLuminance(color1);
  const l2 = getLuminance(color2);

  const lighter = Math.max(l1, l2);
  const darker = Math.min(l1, l2);

  const ratio = (lighter + 0.05) / (darker + 0.05);

  // WCAG AA requiere ratio >= 4.5 para texto normal
  return ratio >= 4.5;
}

/**
 * Maneja navegación por teclado en listas
 * @param event Evento de teclado
 * @param items Array de items navegables
 * @param currentIndex Índice actual
 * @param onSelect Callback cuando se selecciona un item
 * @returns Nuevo índice o null
 */
export function handleKeyboardNavigation(
  event: KeyboardEvent,
  items: any[],
  currentIndex: number,
  onSelect?: (index: number) => void
): number | null {
  let newIndex: number | null = null;

  switch (event.key) {
    case "ArrowDown":
      event.preventDefault();
      newIndex = currentIndex < items.length - 1 ? currentIndex + 1 : 0;
      break;
    case "ArrowUp":
      event.preventDefault();
      newIndex = currentIndex > 0 ? currentIndex - 1 : items.length - 1;
      break;
    case "Home":
      event.preventDefault();
      newIndex = 0;
      break;
    case "End":
      event.preventDefault();
      newIndex = items.length - 1;
      break;
    case "Enter":
    case " ":
      event.preventDefault();
      if (onSelect) {
        onSelect(currentIndex);
      }
      return currentIndex;
  }

  if (newIndex !== null && onSelect) {
    onSelect(newIndex);
  }

  return newIndex;
}

/**
 * Focus trap para modales y diálogos
 */
export function createFocusTrap(container: HTMLElement, onEscape?: () => void) {
  const focusableElements = container.querySelectorAll(
    'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
  );

  const firstElement = focusableElements[0] as HTMLElement;
  const lastElement = focusableElements[
    focusableElements.length - 1
  ] as HTMLElement;

  const handleTab = (e: KeyboardEvent) => {
    if (e.key !== "Tab") return;

    if (e.shiftKey) {
      if (document.activeElement === firstElement) {
        e.preventDefault();
        lastElement?.focus();
      }
    } else {
      if (document.activeElement === lastElement) {
        e.preventDefault();
        firstElement?.focus();
      }
    }
  };

  const handleEscape = (e: KeyboardEvent) => {
    if (e.key === "Escape" && onEscape) {
      onEscape();
    }
  };

  container.addEventListener("keydown", handleTab);
  container.addEventListener("keydown", handleEscape);

  firstElement?.focus();

  return () => {
    container.removeEventListener("keydown", handleTab);
    container.removeEventListener("keydown", handleEscape);
  };
}

/**
 * Anuncia cambios a screen readers
 */
export function announceToScreenReader(
  message: string,
  priority: "polite" | "assertive" = "polite"
) {
  const announcement = document.createElement("div");
  announcement.setAttribute("role", "status");
  announcement.setAttribute("aria-live", priority);
  announcement.setAttribute("aria-atomic", "true");
  announcement.className = "sr-only";
  announcement.textContent = message;

  document.body.appendChild(announcement);

  setTimeout(() => {
    document.body.removeChild(announcement);
  }, 1000);
}
