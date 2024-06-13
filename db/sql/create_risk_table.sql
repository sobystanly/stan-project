CREATE TABLE IF NOT EXISTS risks (
    risk_id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    state TEXT NOT NULL
)