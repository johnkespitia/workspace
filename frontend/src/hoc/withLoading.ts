import { defineComponent, h, Component, PropType } from 'vue';

export function withLoading<T extends Component>(WrappedComponent: T) {
  return defineComponent({
    name: `withLoading(${WrappedComponent.name || 'Component'})`,
    props: {
      loading: {
        type: Boolean,
        default: false,
      },
      loadingText: {
        type: String,
        default: 'Cargando...',
      },
    },
    setup(props, { slots, attrs }) {
      return () => {
        if (props.loading) {
          return h('div', {
            class: 'loading-container flex items-center justify-center min-h-[200px]',
            'aria-live': 'polite',
            'aria-busy': 'true',
          }, [
            h('div', { class: 'flex flex-col items-center gap-4' }, [
              h('div', {
                class: 'spinner w-12 h-12 border-4 border-indigo-200 border-t-indigo-600 rounded-full animate-spin',
              }),
              h('p', { class: 'text-gray-600' }, props.loadingText),
            ]),
          ]);
        }
        return h(WrappedComponent, attrs, slots);
      };
    },
  });
}
