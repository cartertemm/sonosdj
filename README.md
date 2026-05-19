# SonosDJ

Codex and Claude Code are your DJ now.

Tell it what you want to hear. It could be anything: a mood, a scene, a decade, a hard constraint, the room you want to hear it in, "something random", or any combination of these. It will then find something, and play it on your Sonos system automatically.

Say things like:

- `give me 90s Memphis rap, nothing too obvious`
- `put on some 70s rock that I probably won't know in the kitchen`
- `just one artist tonight: Sade`
- `look at my liked songs and play some I probably haven't heard but might like`
- `I'm cleaning and I want to hear some upbeat house, no vocals`
- `give me late 70s New York punk`
- `I'm working late again`

It is meant to feel more like talking to a friend with taste than operating a music app.
If you ever want something else, switch back to your terminal and **just ask**.

## What It Does

SonosDJ is a single launcher that configures your environment, checks to make sure you have Spotify set up within Sonos, then starts either Codex or Claude Code with a comprehensive built-in DJ prompt.

It knows how to:

- control Sonos completely
- search the entire Spotify catalog
- respect mood, genre, artist, era, scene, and negative constraints
- be safer or more adventurous depending on permissiveness
- build a session instead of just tossing on one random song

## Getting Started

SonosDJ runs on Windows, macOS, and Linux. You will need

- Go on your `PATH`
- either `claude` or `codex` on your `PATH`
- Sonos speakers on your local network
- Spotify linked in Sonos

then install the launcher:

```bash
go install github.com/cartertemm/sonosdj@latest
```

If the `sonos` CLI is missing, SonosDJ will offer to install it for you on first run. If you would rather install it yourself:

```bash
go install github.com/steipete/sonoscli/cmd/sonos@latest
```

The first time SonosDJ talks to Spotify through Sonos, it will run an interactive Spotify auth flow against your chosen room. Follow the prompts it prints, then SonosDJ will continue automatically.

## Usage

Run it:

```bash
sonosdj
```

If only one of `claude` and `codex` is on your `PATH`, SonosDJ uses that one. If both are available and you don't pick, SonosDJ asks. Force a specific agent:

```bash
sonosdj --claude
sonosdj --codex
```

Set a default room:

```bash
sonosdj -r "Living Room"
```

Set the default permissiveness:

```bash
sonosdj -p low
sonosdj -p medium
sonosdj -p high
```

Permissiveness controls how far the DJ is willing to stray from a literal reading of your request:

- `low` stays literal and predictable
- `medium` stays close, but can branch out a bit (default)
- `high` follows the mood and goes hunting, even if that means leaving the most obvious path. Ultimate discovery mode

Turn on verbose startup output:

```bash
sonosdj -V
```

You can combine them:

```bash
sonosdj --codex -p high -V
```

## Flags

| Flag | Meaning | Default |
| --- | --- | --- |
| `--claude` | Launch Claude directly | |
| `--codex` | Launch Codex directly | |
| `-r`, `--room` | Default Sonos room (skips the room question on startup) | none |
| `-p`, `--permissive` | Default permissiveness: `low`, `medium`, or `high` | `medium` |
| `-V`, `--verbose` | Print more startup detail | off |
| `-h`, `--help` | Print usage and exit | |

## How It Behaves

The DJ is built to:

- keep chatter short
- avoid obvious lazy picks when there is room for taste

If something is already playing, it tends to stay in that lane unless you tell it otherwise.

If nothing is playing, it will list the available Sonos rooms and ask where to play. Pass `-r` to skip that question.

When launched as Claude, `Bash(sonos:*)` and `Bash(wait:*)` are pre-approved, so the DJ can act without a permission prompt on every command.

## Startup

When you run `sonosdj`, it:

1. finds `claude` and/or `codex`
2. picks one, or asks if both are available
3. checks the `sonos` CLI is installed (and offers to install it if not)
4. checks that Sonos speakers are discoverable
5. checks Spotify through Sonos, and runs the auth flow if needed
6. opens the agent with the DJ prompt already loaded; if you did not pass `-r` and nothing is playing, the DJ will ask which room to use

## Troubleshooting

If startup fails, check these first:

- `claude` or `codex` is installed and on your `PATH`
- `go` is installed and on your `PATH`
- `sonos discover` can see your speakers
- Spotify is linked in Sonos

Useful commands:

- `sonos --version` — confirm the `sonos` CLI is installed and on your `PATH`
- `sonos discover` — list the speakers visible on your network
- `sonos smapi services` — show which music services (like Spotify) are linked in Sonos
