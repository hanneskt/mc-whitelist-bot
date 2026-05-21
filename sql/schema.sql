CREATE TABLE players (
  id     INTEGER PRIMARY KEY,
  mc_uuid      text NOT NULL UNIQUE,
  mc_username  text NOT NULL,
  discord_uuid text NOT NULL UNIQUE,

  country text NOT NULL,
  invited_by text NOT NULL,

  whitelisted  bool NOT NULL DEFAULT true,
  left_server bool NOT NULL DEFAULT false
);
