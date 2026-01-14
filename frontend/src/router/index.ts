import { createRouter, createWebHistory } from "vue-router";

// Code splitting: Lazy loading de rutas para mejorar performance inicial
const StockList = () => import("@/views/StockList.vue");
const StockDetail = () => import("@/views/StockDetail.vue");
const Recommendations = () => import("@/views/Recommendations.vue");

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/",
      redirect: "/stocks",
    },
    {
      path: "/stocks",
      name: "StockList",
      component: StockList,
      meta: {
        title: "Lista de Acciones",
      },
    },
    {
      path: "/stocks/:ticker",
      name: "StockDetail",
      component: StockDetail,
      props: true,
      meta: {
        title: "Detalle de Acción",
      },
    },
    {
      path: "/recommendations",
      name: "Recommendations",
      component: Recommendations,
      meta: {
        title: "Recomendaciones",
      },
    },
  ],
});

// Actualizar título de la página según la ruta
router.beforeEach((to, from, next) => {
  const title = to.meta.title as string;
  if (title) {
    document.title = `${title} - Stock Info`;
  }
  next();
});

export default router;
