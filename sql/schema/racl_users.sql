create table if not exists racl_users (
  -- Standard
  id           text        not null,
  created_at   timestamptz default CURRENT_TIMESTAMP,
  updated_at   timestamptz default CURRENT_TIMESTAMP,

  -- References
  name    text not null,
  api_key text not null,
  api_secret text not null,

  -- Data
  primary key(id),
  constraint min_key_len check(char_length(api_key) >= 24),
  constraint max_key_len check(char_length(api_key) <= 32),
  constraint min_sec_len check(char_length(api_secret) >= 64),
  constraint max_sec_len check(char_length(api_secret) <= 72)
);
