-- Migration: Create stocks table
-- Created: 2024

-- Crear tabla stocks
CREATE TABLE IF NOT EXISTS stocks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticker VARCHAR(10) NOT NULL UNIQUE,
    company_name VARCHAR(255) NOT NULL,
    brokerage VARCHAR(255),
    action VARCHAR(50),
    rating_from VARCHAR(50) NOT NULL,
    rating_to VARCHAR(50) NOT NULL,
    target_from DECIMAL(10,2) NOT NULL,
    target_to DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

-- Índices para optimización de consultas
CREATE INDEX IF NOT EXISTS idx_stocks_ticker ON stocks(ticker);
CREATE INDEX IF NOT EXISTS idx_stocks_rating_to ON stocks(rating_to);
CREATE INDEX IF NOT EXISTS idx_stocks_target_to ON stocks(target_to);
CREATE INDEX IF NOT EXISTS idx_stocks_company_name ON stocks(company_name);

-- Índice compuesto para búsquedas frecuentes
CREATE INDEX IF NOT EXISTS idx_stocks_rating_target ON stocks(rating_to, target_to);

-- Trigger para actualizar updated_at automáticamente
-- Nota: CockroachDB v23.1 puede tener limitaciones con triggers
-- Si el trigger falla, se puede actualizar updated_at manualmente en el código
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Crear el trigger
-- Nota: En algunas versiones de CockroachDB, los triggers pueden no estar completamente soportados
-- Si falla, puedes comentar esta sección y actualizar updated_at manualmente en el código Go
CREATE TRIGGER update_stocks_updated_at 
    BEFORE UPDATE ON stocks 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
