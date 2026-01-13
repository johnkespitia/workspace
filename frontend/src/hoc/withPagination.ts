import { defineComponent, h, Component, computed } from 'vue';

interface PaginationProps {
  currentPage: number;
  totalPages: number;
  pageSize: number;
  totalItems: number;
}

export function withPagination<T extends Component>(WrappedComponent: T) {
  return defineComponent({
    name: `withPagination(${WrappedComponent.name || 'Component'})`,
    props: {
      currentPage: {
        type: Number,
        default: 1,
      },
      totalPages: {
        type: Number,
        required: true,
      },
      pageSize: {
        type: Number,
        default: 10,
      },
      totalItems: {
        type: Number,
        required: true,
      },
    },
    emits: ['page-change'],
    setup(props, { slots, attrs, emit }) {
      const paginationInfo = computed(() => ({
        currentPage: props.currentPage,
        totalPages: props.totalPages,
        pageSize: props.pageSize,
        totalItems: props.totalItems,
        startItem: (props.currentPage - 1) * props.pageSize + 1,
        endItem: Math.min(props.currentPage * props.pageSize, props.totalItems),
      }));

      const goToPage = (page: number) => {
        if (page >= 1 && page <= props.totalPages) {
          emit('page-change', page);
        }
      };

      return () => {
        return h('div', { class: 'pagination-wrapper' }, [
          h(WrappedComponent, {
            ...attrs,
            pagination: paginationInfo.value,
          }, slots),
          h('div', {
            class: 'pagination-controls flex items-center justify-between mt-4 px-4 py-3 bg-white border border-gray-200 rounded-lg',
          }, [
            h('div', { class: 'text-sm text-gray-600' }, [
              `Mostrando ${paginationInfo.value.startItem}-${paginationInfo.value.endItem} de ${props.totalItems}`,
            ]),
            h('div', { class: 'flex items-center gap-2' }, [
              h('button', {
                class: 'px-3 py-1.5 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2',
                disabled: props.currentPage === 1,
                onClick: () => goToPage(props.currentPage - 1),
                'aria-label': 'Página anterior',
              }, '← Anterior'),
              h('span', {
                class: 'px-3 py-1.5 text-sm text-gray-700',
              }, `Página ${props.currentPage} de ${props.totalPages}`),
              h('button', {
                class: 'px-3 py-1.5 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2',
                disabled: props.currentPage === props.totalPages,
                onClick: () => goToPage(props.currentPage + 1),
                'aria-label': 'Página siguiente',
              }, 'Siguiente →'),
            ]),
          ]),
        ]);
      };
    },
  });
}
