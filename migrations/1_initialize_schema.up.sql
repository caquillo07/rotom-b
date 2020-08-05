-- server config
CREATE TABLE server_config (
  id SERIAL PRIMARY KEY,
  guild_name TEXT NOT NULL,
  guild_id TEXT NOT NULL,
  last_updated_by TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_server_config_guild_id on server_config (guild_id);
