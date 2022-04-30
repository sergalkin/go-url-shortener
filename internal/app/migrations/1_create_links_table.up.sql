create table if not exists links(
    id bigserial primary key,
    url_hash varchar(10) not null,
    url text not null,
    uid uuid not null,
    created_at timestamptz default NOW()
)