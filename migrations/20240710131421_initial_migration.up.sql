create table ponds (
    id bigserial primary key,
    farm_id bigint not null,
    name varchar(50) not null,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now(),
    deleted_at timestamp with time zone
);

create table farms (
    id bigserial primary key,
    name varchar(50) not null,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now(),
    deleted_at timestamp with time zone    
);
