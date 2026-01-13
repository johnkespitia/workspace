import { createRouter, createWebHistory } from 'vue-router';
import StockList from '@/views/StockList.vue';
import StockDetail from '@/views/StockDetail.vue';
import Recommendations from '@/views/Recommendations.vue';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: '/stocks',
    },
    {
      path: '/stocks',
      name: 'StockList',
      component: StockList,
    },
    {
      path: '/stocks/:ticker',
      name: 'StockDetail',
      component: StockDetail,
      props: true,
    },
    {
      path: '/recommendations',
      name: 'Recommendations',
      component: Recommendations,
    },
  ],
});

export default router;
