create table if not exists racl_resources (
  -- Standard
  id           text        not null,
  created_at   timestamptz default CURRENT_TIMESTAMP,
  updated_at   timestamptz default CURRENT_TIMESTAMP,

  -- References 

  -- Data

  -- Constraints
  primary key(id),
  constraint resource_id_min_len check(char_length(id) > 0)
);
