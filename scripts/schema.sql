CREATE TABLE IF NOT EXISTS user_profiles (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    gender TEXT NOT NULL,
    gender_probability NUMERIC(3, 2),
    sample_size BIGINT,
    age BIGINT,
    age_group TEXT NOT NULL,
    country_id TEXT NOT NULL,
    country_probability NUMERIC(3, 2),
    created_at TIMESTAMPTZ DEFAULT NOW()
);
