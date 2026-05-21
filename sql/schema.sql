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

CREATE TABLE birthdays (
  id     INTEGER PRIMARY KEY,
  discord_uuid TEXT NOT NULL UNIQUE,
  day INTEGER NOT NULL,
  month INTEGER NOT NULL,

  FOREIGN KEY(discord_uuid) REFERENCES players(discord_uuid) ON DELETE CASCADE
);
