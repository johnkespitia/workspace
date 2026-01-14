// Cliente GraphQL para comunicación con el backend
import { Client, cacheExchange, fetchExchange } from "@urql/core";

// Configuración del cliente GraphQL
// Desde el contenedor, usar el nombre del servicio; desde el host, usar localhost
const getGraphQLEndpoint = () => {
  const envEndpoint = import.meta.env.VITE_GRAPHQL_ENDPOINT;
  if (envEndpoint) {
    return envEndpoint;
  }
  // Si estamos en el navegador (cliente), usar localhost
  // Si estamos en el servidor (SSR), usar el nombre del servicio
  return typeof window !== "undefined"
    ? "http://localhost:8080/query"
    : "http://api:8080/query";
};

const GRAPHQL_ENDPOINT = getGraphQLEndpoint();

export const graphqlClient = new Client({
  url: GRAPHQL_ENDPOINT,
  exchanges: [cacheExchange, fetchExchange],
});

// Tipos para las queries GraphQL
export interface Stock {
  id: string;
  ticker: string;
  companyName: string;
  brokerage?: string;
  action?: string;
  ratingFrom?: string;
  ratingTo?: string;
  targetFrom?: number;
  targetTo?: number;
  createdAt: string;
  updatedAt: string;
}

export interface StockFilter {
  ticker?: string; // Búsqueda exacta por ticker
  companyName?: string; // Búsqueda parcial por nombre de empresa
  ratings?: string[]; // ✅ CORRECTO: Array de ratings (PLURAL)
  action?: string; // Búsqueda exacta por acción
}

export type StockSortField =
  | "TICKER"
  | "COMPANY_NAME"
  | "RATING_TO"
  | "TARGET_TO"
  | "CREATED_AT";
export type SortDirection = "ASC" | "DESC";

export interface StockSort {
  field?: StockSortField; // Campo por el cual ordenar (opcional)
  direction?: SortDirection; // Dirección del ordenamiento (opcional)
}

export interface StockConnection {
  stocks: Stock[];
  totalCount: number;
  pageInfo: {
    hasNextPage: boolean;
    hasPreviousPage: boolean;
  };
}

export interface Recommendation {
  stock: Stock;
  score: number;
  priceChange: number;
  ratingScore: number;
  actionScore: number;
}

// Queries GraphQL
export const GET_STOCKS_QUERY = `
  query GetStocks($filter: StockFilter, $sort: StockSort, $limit: Int, $offset: Int) {
    stocks(filter: $filter, sort: $sort, limit: $limit, offset: $offset) {
      stocks {
        id
        ticker
        companyName
        brokerage
        action
        ratingFrom
        ratingTo
        targetFrom
        targetTo
        createdAt
        updatedAt
      }
      totalCount
      pageInfo {
        hasNextPage
        hasPreviousPage
      }
    }
  }
`;

export const GET_STOCK_QUERY = `
  query GetStock($ticker: String!) {
    stock(ticker: $ticker) {
      id
      ticker
      companyName
      brokerage
      action
      ratingFrom
      ratingTo
      targetFrom
      targetTo
      createdAt
      updatedAt
    }
  }
`;

export const GET_RECOMMENDATIONS_QUERY = `
  query GetRecommendations($limit: Int) {
    recommendations(limit: $limit) {
      stock {
        id
        ticker
        companyName
        brokerage
        action
        ratingFrom
        ratingTo
        targetFrom
        targetTo
        createdAt
        updatedAt
      }
      score
      priceChange
      ratingScore
      actionScore
    }
  }
`;

export const SYNC_STOCKS_MUTATION = `
  mutation SyncStocks {
    syncStocks {
      success
      message
      stocksSynced
    }
  }
`;
