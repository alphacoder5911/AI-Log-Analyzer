CREATE TABLE incidents (
    id SERIAL PRIMARY KEY,

    service_name TEXT NOT NULL,

    level TEXT NOT NULL,

    message TEXT NOT NULL,

    root_cause TEXT,

    impact TEXT,

    suggested_fix TEXT,

    severity TEXT,

    confidence INTEGER,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);