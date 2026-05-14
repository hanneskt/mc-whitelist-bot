# Minecraft Whitelist Discord Bot

A discord bot for whitelisting people after filling in a form on discord.

## Setup

Run once to generate the `config.json` file
It's a file with these fields:
|Field|Value|
|---|---|
|`"token"`|The secret bot token from the [discord developer portal](https://discord.com/developers/home)|
|`"guild_id"`|The snowflake ID of the guild this bot is in|
|`"welcome_channel_id"`|The snowflake ID of the channel the bot should post welcome messages|
|`"server_name"`|The name of your server, used in the application form|

## Running

```bash
go run .
```
