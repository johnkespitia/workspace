import type { Meta, StoryObj } from '@storybook/vue3';
import Card from './Card.vue';

const meta: Meta<typeof Card> = {
  title: 'Design System/Card',
  component: Card,
  tags: ['autodocs'],
  argTypes: {
    variant: {
      control: 'select',
      options: ['default', 'elevated', 'outlined'],
    },
    padding: {
      control: 'select',
      options: ['none', 'sm', 'md', 'lg'],
    },
  },
};

export default meta;
type Story = StoryObj<typeof Card>;

export const Default: Story = {
  args: {
    title: 'Título de la Tarjeta',
    children: 'Contenido de la tarjeta',
  },
  render: (args) => ({
    components: { Card },
    setup() {
      return { args };
    },
    template: '<Card v-bind="args">{{ args.children }}</Card>',
  }),
};

export const Elevated: Story = {
  args: {
    title: 'Tarjeta Elevada',
    variant: 'elevated',
    children: 'Esta tarjeta tiene sombra',
  },
  render: (args) => ({
    components: { Card },
    setup() {
      return { args };
    },
    template: '<Card v-bind="args">{{ args.children }}</Card>',
  }),
};
