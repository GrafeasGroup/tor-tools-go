# ToR Tools

A toolchain for TranscribersOfReddit maintenance concerns

## Build

```shell
~/workspace/tor-tools-go $ make clean all
```

## Usage

```bash
#!/usr/bin/env bash

set -euo pipefail

SLACK_CHANNEL='red-alerts'
SLACK_TOKEN='<generate-your-own>'
BOT_USERNAMES='bot1 bot2 bot3'  # <- these are reddit.com usernames

export SLACK_CHANNEL SLACK_TOKEN BOT_USERNAMES
./output/tor-tools-darwin-amd64 # <- for MacOS, otherwise choose your architecture
```
