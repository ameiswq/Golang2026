create table if not exists idempotency_keys (
    key text primary key, 
    status text not null,
    response_code integer,
    response_body text,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);