// Cliente GraphQL para comunicación con el backend
import { Client, cacheExchange, fetchExchange } from '@urql/core';

// Configuración del cliente GraphQL
const GRAPHQL_ENDPOINT = import.meta.env.VITE_GRAPHQL_ENDPOINT || 'http://localhost:8080/query';

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
  ticker?: string;
  companyName?: string;
  rating?: string[];
  action?: string;
}

export interface StockSort {
  field: 'ticker' | 'companyName' | 'ratingTo' | 'targetTo' | 'createdAt';
  direction: 'asc' | 'desc';
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
