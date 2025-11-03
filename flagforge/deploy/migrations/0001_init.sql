CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS orgs (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    owner_id UUID,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS projects (
    id UUID PRIMARY KEY,
    org_id UUID REFERENCES orgs(id),
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS environments (
    id UUID PRIMARY KEY,
    project_id UUID REFERENCES projects(id),
    key TEXT NOT NULL CHECK (key IN ('dev','stage','prod'))
);

CREATE UNIQUE INDEX IF NOT EXISTS environments_project_key ON environments(project_id, key);

CREATE TABLE IF NOT EXISTS flags (
    id UUID PRIMARY KEY,
    project_id UUID REFERENCES projects(id),
    key TEXT NOT NULL UNIQUE,
    type TEXT NOT NULL CHECK (type IN ('bool','number','string','json')),
    description TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS flag_versions (
    id UUID PRIMARY KEY,
    flag_id UUID REFERENCES flags(id),
    environment_id UUID REFERENCES environments(id),
    value_json JSONB NOT NULL,
    rollout INT DEFAULT 100,
    rules_json JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by UUID
);

CREATE INDEX IF NOT EXISTS flag_versions_flag_env ON flag_versions(flag_id, environment_id);

CREATE TABLE IF NOT EXISTS api_keys (
    id UUID PRIMARY KEY,
    environment_id UUID REFERENCES environments(id),
    key TEXT NOT NULL UNIQUE,
    role TEXT NOT NULL CHECK (role IN ('client','server')),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY,
    actor_id UUID,
    org_id UUID,
    entity_type TEXT NOT NULL,
    entity_id UUID NOT NULL,
    action TEXT NOT NULL,
    diff_json JSONB,
    ts TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS audit_logs_entity ON audit_logs(entity_id, ts DESC);

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email CITEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
