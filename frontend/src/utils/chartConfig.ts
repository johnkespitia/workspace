// Configuración centralizada de Chart.js
// Esto asegura que los componentes estén registrados antes de usar los componentes de vue-chartjs
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  LineController,
  Title,
  Tooltip,
  Legend,
  Filler,
} from "chart.js";

// Registrar todos los componentes necesarios una sola vez
ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  LineController,
  Title,
  Tooltip,
  Legend,
  Filler
);

export default ChartJS;
