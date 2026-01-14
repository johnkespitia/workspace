import type { Meta, StoryObj } from "@storybook/vue3";
import ThemeToggle from "./ThemeToggle.vue";
import { createPinia, setActivePinia } from "pinia";

const meta: Meta<typeof ThemeToggle> = {
  title: "Design System/ThemeToggle",
  component: ThemeToggle,
  tags: ["autodocs"],
  decorators: [
    (story) => {
      // Configurar Pinia para los tests
      const pinia = createPinia();
      setActivePinia(pinia);
      return {
        components: { story },
        template: '<div class="p-4"><story /></div>',
      };
    },
  ],
};

export default meta;
type Story = StoryObj<typeof ThemeToggle>;

export const Default: Story = {
  args: {},
};

export const LightTheme: Story = {
  args: {},
  parameters: {
    backgrounds: {
      default: "light",
    },
  },
};

export const DarkTheme: Story = {
  args: {},
  parameters: {
    backgrounds: {
      default: "dark",
    },
  },
};
