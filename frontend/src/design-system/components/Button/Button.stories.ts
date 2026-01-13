import type { Meta, StoryObj } from '@storybook/vue3';
import Button from './Button.vue';

const meta: Meta<typeof Button> = {
  title: 'Design System/Button',
  component: Button,
  tags: ['autodocs'],
  argTypes: {
    variant: {
      control: 'select',
      options: ['primary', 'secondary', 'outline', 'ghost', 'danger'],
    },
    size: {
      control: 'select',
      options: ['sm', 'md', 'lg'],
    },
    disabled: {
      control: 'boolean',
    },
    loading: {
      control: 'boolean',
    },
  },
};

export default meta;
type Story = StoryObj<typeof Button>;

export const Primary: Story = {
  args: {
    variant: 'primary',
    children: 'Botón Primario',
  },
  render: (args) => ({
    components: { Button },
    setup() {
      return { args };
    },
    template: '<Button v-bind="args">{{ args.children }}</Button>',
  }),
};

export const Secondary: Story = {
  args: {
    variant: 'secondary',
    children: 'Botón Secundario',
  },
  render: (args) => ({
    components: { Button },
    setup() {
      return { args };
    },
    template: '<Button v-bind="args">{{ args.children }}</Button>',
  }),
};

export const Loading: Story = {
  args: {
    variant: 'primary',
    loading: true,
    loadingText: 'Cargando...',
  },
  render: (args) => ({
    components: { Button },
    setup() {
      return { args };
    },
    template: '<Button v-bind="args">Texto del botón</Button>',
  }),
};
