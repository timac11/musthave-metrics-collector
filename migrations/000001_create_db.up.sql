CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE metrics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    mtype VARCHAR(10) NOT NULL CHECK (mtype IN ('gauge', 'counter')),
    delta BIGINT NULL,
    value DOUBLE PRECISION NULL,
    hash VARCHAR(255) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT check_type_value_consistency CHECK (
        (mtype = 'counter' AND delta IS NOT NULL AND value IS NULL) OR
        (mtype = 'gauge' AND delta IS NULL AND value IS NOT NULL)
    ),
    
    -- Ensure each metric name has only one type
    CONSTRAINT unique_metric_name_type UNIQUE (name, mtype)
);