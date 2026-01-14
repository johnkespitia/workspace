import type { Meta, StoryObj } from "@storybook/vue3";
import Table from "./Table.vue";
import type { TableHeader } from "./Table.vue";

const meta: Meta<typeof Table> = {
  title: "Design System/Table",
  component: Table,
  tags: ["autodocs"],
  argTypes: {
    sortable: {
      control: "boolean",
    },
    striped: {
      control: "boolean",
    },
    hoverable: {
      control: "boolean",
    },
    clickable: {
      control: "boolean",
    },
  },
};

export default meta;
type Story = StoryObj<typeof Table>;

const sampleHeaders: TableHeader[] = [
  { key: "ticker", label: "Ticker", sortable: true },
  { key: "companyName", label: "Compañía", sortable: true },
  { key: "rating", label: "Rating", sortable: true },
  { key: "target", label: "Target", sortable: true },
];

const sampleData = [
  {
    ticker: "AAPL",
    companyName: "Apple Inc.",
    rating: "Buy",
    target: "$180.00",
  },
  {
    ticker: "MSFT",
    companyName: "Microsoft Corporation",
    rating: "Strong Buy",
    target: "$420.00",
  },
  {
    ticker: "GOOGL",
    companyName: "Alphabet Inc.",
    rating: "Buy",
    target: "$150.00",
  },
];

export const Default: Story = {
  args: {
    headers: sampleHeaders,
    data: sampleData,
    sortable: true,
    striped: true,
    hoverable: true,
    clickable: false,
    ariaLabel: "Tabla de acciones",
  },
};

export const WithEmptyState: Story = {
  args: {
    headers: sampleHeaders,
    data: [],
    sortable: true,
    striped: true,
    hoverable: true,
    ariaLabel: "Tabla vacía",
  },
};

export const ClickableRows: Story = {
  args: {
    headers: sampleHeaders,
    data: sampleData,
    sortable: true,
    striped: true,
    hoverable: true,
    clickable: true,
    ariaLabel: "Tabla con filas clickeables",
  },
  play: async ({ args }) => {
    // Simular click en fila
    console.log("Row clicked", args);
  },
};

export const WithoutStripes: Story = {
  args: {
    headers: sampleHeaders,
    data: sampleData,
    sortable: true,
    striped: false,
    hoverable: true,
    ariaLabel: "Tabla sin rayas",
  },
};

export const WithCustomCells: Story = {
  args: {
    headers: sampleHeaders,
    data: sampleData,
    sortable: true,
    striped: true,
    hoverable: true,
    ariaLabel: "Tabla con celdas personalizadas",
  },
  render: (args) => ({
    components: { Table },
    setup() {
      return { args };
    },
    template: `
      <Table v-bind="args">
        <template #cell-rating="{ value }">
          <span class="px-2 py-1 rounded text-xs font-semibold bg-blue-100 text-blue-800">
            {{ value }}
          </span>
        </template>
      </Table>
    `,
  }),
};
