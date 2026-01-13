import type { Meta, StoryObj } from '@storybook/vue3';
import Input from './Input.vue';

const meta: Meta<typeof Input> = {
  title: 'Design System/Input',
  component: Input,
  tags: ['autodocs'],
  argTypes: {
    type: {
      control: 'select',
      options: ['text', 'email', 'password', 'number', 'search'],
    },
  },
};

export default meta;
type Story = StoryObj<typeof Input>;

export const Default: Story = {
  args: {
    modelValue: '',
    label: 'Nombre',
    placeholder: 'Ingrese su nombre',
  },
};

export const WithError: Story = {
  args: {
    modelValue: '',
    label: 'Email',
    placeholder: 'Ingrese su email',
    error: 'El email no es válido',
  },
};

export const WithHint: Story = {
  args: {
    modelValue: '',
    label: 'Contraseña',
    placeholder: 'Ingrese su contraseña',
    hint: 'Mínimo 8 caracteres',
    type: 'password',
  },
};
