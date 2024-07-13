create table request_logs (
    id bigserial primary key,
    endpoint text not null, -- endpoint path, ex: GET /ping
    latency real not null default 0.00, -- request latency time in (ms)
    user_agent text not null default 0,  -- useragent detected on request
    requested_at timestamp with time zone not null default now()
);
