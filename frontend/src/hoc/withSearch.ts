import { defineComponent, h, Component, ref, watch } from 'vue';
import { useDebounce } from '@/composables/useDebounce';

export function withSearch<T extends Component>(WrappedComponent: T) {
  return defineComponent({
    name: `withSearch(${WrappedComponent.name || 'Component'})`,
    props: {
      searchPlaceholder: {
        type: String,
        default: 'Buscar...',
      },
      debounceMs: {
        type: Number,
        default: 300,
      },
    },
    emits: ['search'],
    setup(props, { slots, attrs, emit }) {
      const searchQuery = ref('');
      const debouncedQuery = useDebounce(searchQuery, props.debounceMs);

      watch(debouncedQuery, (newQuery) => {
        emit('search', newQuery);
      });

      return () => {
        return h('div', { class: 'search-wrapper' }, [
          h('div', { class: 'search-input-container mb-4' }, [
            h('div', { class: 'relative' }, [
              h('input', {
                type: 'text',
                class: 'w-full px-4 py-2 pl-10 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200',
                placeholder: props.searchPlaceholder,
                value: searchQuery.value,
                onInput: (e: Event) => {
                  searchQuery.value = (e.target as HTMLInputElement).value;
                },
                'aria-label': 'Campo de búsqueda',
              }),
              h('span', {
                class: 'absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400',
                'aria-hidden': 'true',
              }, '🔍'),
            ]),
          ]),
          h(WrappedComponent, {
            ...attrs,
            searchQuery: debouncedQuery.value,
          }, slots),
        ]);
      };
    },
  });
}
