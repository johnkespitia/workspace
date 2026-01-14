import { describe, it, expect } from "vitest";
import { mount } from "@vue/test-utils";
import { defineComponent, h } from "vue";
import { withLoading } from "@/hoc/withLoading";

// Componente de prueba para envolver
const TestComponent = defineComponent({
  name: "TestComponent",
  setup() {
    return () => h("div", { class: "test-content" }, "Test Content");
  },
});

describe("withLoading HOC", () => {
  it("renders wrapped component when not loading", () => {
    const WrappedComponent = withLoading(TestComponent);
    const wrapper = mount(WrappedComponent, {
      props: {
        loading: false,
      },
    });

    expect(wrapper.find(".test-content").exists()).toBe(true);
    expect(wrapper.text()).toBe("Test Content");
  });

  it("renders loading state when loading is true", () => {
    const WrappedComponent = withLoading(TestComponent);
    const wrapper = mount(WrappedComponent, {
      props: {
        loading: true,
        loadingText: "Cargando datos...",
      },
    });

    expect(wrapper.find(".test-content").exists()).toBe(false);
    expect(wrapper.find(".loading-container").exists()).toBe(true);
    expect(wrapper.find(".spinner").exists()).toBe(true);
    expect(wrapper.text()).toContain("Cargando datos...");
  });

  it("has correct ARIA attributes when loading", () => {
    const WrappedComponent = withLoading(TestComponent);
    const wrapper = mount(WrappedComponent, {
      props: {
        loading: true,
      },
    });

    const container = wrapper.find(".loading-container");
    expect(container.attributes("aria-live")).toBe("polite");
    expect(container.attributes("aria-busy")).toBe("true");
  });

  it("uses default loading text when not provided", () => {
    const WrappedComponent = withLoading(TestComponent);
    const wrapper = mount(WrappedComponent, {
      props: {
        loading: true,
      },
    });

    expect(wrapper.text()).toContain("Cargando...");
  });
});
