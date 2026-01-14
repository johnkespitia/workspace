import { afterEach } from "vitest";
import { cleanup } from "@testing-library/vue";

// Limpiar después de cada test
afterEach(() => {
  cleanup();
});
