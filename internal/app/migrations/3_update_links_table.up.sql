alter table links
add constraint links_url_key
unique (url);