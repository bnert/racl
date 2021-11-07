create table if not exists racl_acls (
  -- Standard
  id           uuid        default uuid_generate_v4(),
  created_at   timestamptz default CURRENT_TIMESTAMP,
  updated_at   timestamptz default CURRENT_TIMESTAMP,

  -- References
  resource_id  text        references racl_resources(id) ON DELETE CASCADE not null,

  -- Data
  entity       text        not null,
  capabilities char[]      not null,

  -- Constraints
  primary key (id),
  unique(resource_id, entity),
  constraint entity_length check(char_length(entity) > 0),
  constraint valid_capabilities check(capabilities <@ '{"c", "r", "u", "d", "a"}')
);
