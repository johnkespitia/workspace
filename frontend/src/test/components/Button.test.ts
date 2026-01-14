import { describe, it, expect } from "vitest";
import { mount } from "@vue/test-utils";
import Button from "@/design-system/components/Button/Button.vue";

describe("Button Component", () => {
  it("renders correctly with default props", () => {
    const wrapper = mount(Button, {
      slots: {
        default: "Click me",
      },
    });

    expect(wrapper.text()).toBe("Click me");
    expect(wrapper.classes()).toContain("bg-indigo-600");
  });

  it("applies primary variant correctly", () => {
    const wrapper = mount(Button, {
      props: {
        variant: "primary",
      },
      slots: {
        default: "Primary Button",
      },
    });

    expect(wrapper.classes()).toContain("bg-indigo-600");
  });

  it("applies secondary variant correctly", () => {
    const wrapper = mount(Button, {
      props: {
        variant: "secondary",
      },
      slots: {
        default: "Secondary Button",
      },
    });

    expect(wrapper.classes()).toContain("bg-purple-600");
  });

  it("shows loading state correctly", () => {
    const wrapper = mount(Button, {
      props: {
        loading: true,
        loadingText: "Cargando...",
      },
      slots: {
        default: "Button Text",
      },
    });

    expect(wrapper.find(".spinner").exists()).toBe(true);
    expect(wrapper.text()).toContain("Cargando...");
  });

  it("is disabled when disabled prop is true", () => {
    const wrapper = mount(Button, {
      props: {
        disabled: true,
      },
      slots: {
        default: "Disabled Button",
      },
    });

    expect(wrapper.attributes("disabled")).toBeDefined();
    expect(wrapper.classes()).toContain("disabled:opacity-50");
  });

  it("emits click event when clicked", async () => {
    const wrapper = mount(Button, {
      slots: {
        default: "Click me",
      },
    });

    await wrapper.trigger("click");
    expect(wrapper.emitted("click")).toBeTruthy();
    expect(wrapper.emitted("click")?.length).toBe(1);
  });

  it("does not emit click when disabled", async () => {
    const wrapper = mount(Button, {
      props: {
        disabled: true,
      },
      slots: {
        default: "Disabled Button",
      },
    });

    await wrapper.trigger("click");
    expect(wrapper.emitted("click")).toBeFalsy();
  });

  it("does not emit click when loading", async () => {
    const wrapper = mount(Button, {
      props: {
        loading: true,
      },
      slots: {
        default: "Loading Button",
      },
    });

    await wrapper.trigger("click");
    expect(wrapper.emitted("click")).toBeFalsy();
  });

  it("applies correct size classes", () => {
    const sizes = ["sm", "md", "lg"] as const;

    sizes.forEach((size) => {
      const wrapper = mount(Button, {
        props: { size },
        slots: { default: "Button" },
      });

      if (size === "sm") {
        expect(wrapper.classes()).toContain("px-3");
      } else if (size === "md") {
        expect(wrapper.classes()).toContain("px-4");
      } else if (size === "lg") {
        expect(wrapper.classes()).toContain("px-6");
      }
    });
  });

  it("applies fullWidth class when fullWidth is true", () => {
    const wrapper = mount(Button, {
      props: {
        fullWidth: true,
      },
      slots: {
        default: "Full Width Button",
      },
    });

    expect(wrapper.classes()).toContain("w-full");
  });
});
