import { defineComponent, h, Component, PropType } from 'vue';

export function withError<T extends Component>(WrappedComponent: T) {
  return defineComponent({
    name: `withError(${WrappedComponent.name || 'Component'})`,
    props: {
      error: {
        type: [String, Error, null] as PropType<string | Error | null>,
        default: null,
      },
      errorMessage: {
        type: String,
        default: 'Ha ocurrido un error',
      },
    },
    emits: ['retry'],
    setup(props, { slots, attrs, emit }) {
      const getErrorMessage = () => {
        if (!props.error) return null;
        if (typeof props.error === 'string') return props.error;
        return props.error.message || props.errorMessage;
      };

      return () => {
        const errorMsg = getErrorMessage();

        if (errorMsg) {
          return h('div', {
            class: 'error-container p-6 bg-red-50 border border-red-200 rounded-lg',
            role: 'alert',
            'aria-live': 'assertive',
          }, [
            h('div', { class: 'flex items-start gap-4' }, [
              h('div', { class: 'text-2xl', 'aria-hidden': 'true' }, '⚠️'),
              h('div', { class: 'flex-1' }, [
                h('h3', {
                  class: 'text-lg font-semibold text-red-800 mb-2',
                }, 'Error'),
                h('p', {
                  class: 'text-red-700 mb-4',
                }, errorMsg),
                h('button', {
                  class: 'px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2',
                  onClick: () => emit('retry'),
                  'aria-label': 'Reintentar',
                }, 'Reintentar'),
              ]),
            ]),
          ]);
        }

        return h(WrappedComponent, attrs, slots);
      };
    },
  });
}
