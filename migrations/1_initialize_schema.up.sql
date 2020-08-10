-- Make sure the timezone is consistent DB wise
SET TIME ZONE 'UTC';

-- server config
CREATE TABLE guild_settings (
  id SERIAL PRIMARY KEY,
  guild_name TEXT NOT NULL,
  guild_id TEXT NOT NULL,
  last_updated_by TEXT NOT NULL,
  bot_prefix TEXT NOT NULL,
  bot_admin_roles TEXT ARRAY,
  bot_admin_users TEXT ARRAY,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_guild_settings_guild_id on guild_settings (guild_id);
