# get-ip

A beautiful terminal UI that shows every network interface on your Mac at a glance: device, IPv4, MAC, gateway, and whether the address came from DHCP or was set manually. Built in Go with [Bubble Tea](https://github.com/charmbracelet/bubbletea), [Bubbles](https://github.com/charmbracelet/bubbles), and [Lip Gloss](https://github.com/charmbracelet/lipgloss).

```
🌐 get-ip

╭───────────────────────────────────────────────────────────────────────────────────────╮
│  Device                  IPv4              MAC                  Gateway        Mode   │
│ ─────────────────────────────────────────────────────────────────────────────────────│
│  en0 (Wi-Fi)             192.168.8.207     0e:be:34:8c:d0:81    192.168.8.1    DHCP   │
│  utun7                   100.74.106.74     —                    192.168.8.1    Manual │
╰───────────────────────────────────────────────────────────────────────────────────────╯

  r refresh • w wifi • q quit
```

> Replace the ASCII above with a real screenshot once captured. Save as `docs/screenshot.png` and reference it here.

---

## Table of contents

- [Features](#features)
- [Requirements](#requirements)
- [Install](#install)
  - [Option 1: Homebrew (recommended)](#option-1-homebrew-recommended)
  - [Option 2: Makefile (manual)](#option-2-makefile-manual)
  - [Option 3: Run from source (no install)](#option-3-run-from-source-no-install)
- [The `ip` alias](#the-ip-alias)
- [Usage](#usage)
- [How it works](#how-it-works)
- [Troubleshooting](#troubleshooting)
- [Uninstall](#uninstall)
- [Development](#development)
- [Releasing a new version](#releasing-a-new-version)
- [License](#license)

---

## Features

- **One-glance network table** — device, IPv4, MAC, gateway, mode (DHCP / Manual) for every active interface.
- **Async refresh** (`r`) — re-runs detection without blocking the UI.
- **Current Wi-Fi panel** (`w`) — toggles a panel showing the connected SSID, channel, and signal (RSSI).
- **Friendly device labels** — `en0` is shown as `en0 (Wi-Fi)` using `networksetup -listnetworkserviceorder`.
- **Zero non-stdlib dependencies** beyond the four Charm libraries.
- **Boot-fast** — single static binary, no runtime config.

## Requirements

- **macOS** (primary target). The detection layer shells out to `networksetup`, `ipconfig`, `route`, and `system_profiler` — all stock macOS tools.
- **Go 1.22+** — only required if you build from source. Homebrew installs Go automatically as a build dep.
- **A real terminal** — Bubble Tea needs a TTY; piping output won't work.

Linux and Windows builds compile (`other.go` provides a stub) but `Collect()` returns `unsupported platform`. PRs welcome.

---

## Install

### Option 1: Homebrew (recommended)

Once a tagged release exists at `github.com/vianch/get-ip`, install via the in-repo formula:

```bash
brew tap vianch/get-ip https://github.com/vianch/get-ip
brew install vianch/get-ip/get-ip
```

To upgrade later:

```bash
brew update
brew upgrade get-ip
```

To uninstall:

```bash
brew uninstall get-ip
brew untap vianch/get-ip
```

> **Note:** Until v0.1.0 is tagged and `Formula/get-ip.rb` is updated with the real `sha256`, Homebrew install will fail. See [Releasing a new version](#releasing-a-new-version) for the publish workflow. To install from a local checkout right now:
>
> ```bash
> brew install --build-from-source ./Formula/get-ip.rb
> ```

### Option 2: Makefile (manual)

Best when you've cloned the repo and want a global install without Homebrew.

```bash
git clone https://github.com/vianch/get-ip.git
cd get-ip
make build      # produces ./bin/get-ip
make install    # copies to /usr/local/bin/get-ip
```

Override the install prefix if `/usr/local` requires sudo on your machine:

```bash
make install PREFIX=$HOME/.local
# then ensure $HOME/.local/bin is on your $PATH
```

All Makefile targets:

| Target                       | What it does                                                          |
| ---------------------------- | --------------------------------------------------------------------- |
| `make help`                  | Print this list (default target).                                     |
| `make build`                 | Compile to `bin/get-ip` with version info baked in.                   |
| `make install`               | `make build` + copy to `$(PREFIX)/bin/get-ip` (default `/usr/local`). |
| `make uninstall`             | Remove binary and the alias block from `~/.zshrc`.                    |
| `make alias`                 | Append a guarded `alias ip='get-ip'` block to `~/.zshrc` (idempotent).|
| `make unalias`               | Strip the alias block from `~/.zshrc` (keeps a `.bak`).               |
| `make run`                   | `go run ./` — handy during development.                               |
| `make test`                  | `go test ./...`.                                                      |
| `make fmt`                   | `gofmt -w . && go vet ./...`.                                         |
| `make clean`                 | Remove `bin/`.                                                        |
| `make release V=vX.Y.Z`      | Tag a release and print the formula update snippet.                   |
| `make formula-sha V=vX.Y.Z`  | Print SHA256 of an existing release tarball.                          |

### Option 3: Run from source (no install)

```bash
git clone https://github.com/vianch/get-ip.git
cd get-ip
go run ./
```

Or via Make:

```bash
make run
```

---

## The `ip` alias

The user-facing binary is `get-ip` — short, unambiguous, and doesn't collide with anything on a stock macOS install. If you'd like to invoke it as `ip`, opt in:

```bash
make alias        # appends a guarded block to ~/.zshrc
source ~/.zshrc   # reload your shell
ip                # launches get-ip
```

The block is idempotent (re-running `make alias` is a no-op):

```sh
# >>> get-ip alias >>>
alias ip='get-ip'
# <<< get-ip alias <<<
```

> **Heads up:** Linux's [iproute2](https://wiki.linuxfoundation.org/networking/iproute2) ships an `ip` command. If you ever install it on macOS (e.g., via `brew install iproute2mac`), this alias will shadow it. Remove with `make unalias`.

To remove the alias without uninstalling the binary:

```bash
make unalias
source ~/.zshrc
```

---

## Usage

Just run it:

```bash
get-ip
# or, if you ran `make alias`:
ip
```

Flags:

| Flag             | Effect                          |
| ---------------- | ------------------------------- |
| `-h`, `--help`   | Print usage and exit.           |
| `-v`, `--version`| Print build version and exit.   |

Keybindings inside the TUI:

| Key            | Action                                      |
| -------------- | ------------------------------------------- |
| `r`            | Refresh interface data.                     |
| `w`            | Toggle the current-Wi-Fi panel.             |
| `↑` / `↓`      | Move row selection.                         |
| `q`, `Ctrl-C`  | Quit.                                       |

---

## How it works

| Field          | macOS source                                                         |
| -------------- | -------------------------------------------------------------------- |
| Device + IPv4 + MAC | Go stdlib `net.Interfaces()` filtered to up + non-loopback + has-IPv4. |
| Friendly name  | `networksetup -listnetworkserviceorder` parsed for `(Hardware Port: X, Device: enN)` lines. |
| Gateway        | `route -n get -ifscope <iface> default`, falling back to the system default route. |
| **Mode (DHCP/Manual)** | `ipconfig getpacket <iface>` — non-empty stdout with a `yiaddr` line ⇒ `DHCP`, otherwise `Manual`. |
| Wi-Fi SSID/signal | `system_profiler SPAirPortDataType -json` — extracts `spairport_current_network_information`. |

Detection runs in a Bubble Tea `tea.Cmd` so the UI stays interactive while shellouts execute. The Wi-Fi panel lazy-loads on first `w` press and is cached until the next refresh.

---

## Troubleshooting

**The Wi-Fi panel says `(SSID hidden — grant Location Services to Terminal)`.**
Starting in macOS 14, Apple gates SSID visibility behind Location Services. Open `System Settings → Privacy & Security → Location Services` and enable it for your terminal app (Terminal, iTerm2, Ghostty, etc.). Then refresh with `r`.

**A row's MAC or Gateway is `—`.**
That field couldn't be detected for that interface. Common case: VPN tunnels (`utun*`) often have no MAC, and not every interface has its own gateway entry — `route -n get default` is consulted as a fallback.

**`get-ip: open /dev/tty: device not configured`.**
You're running in a non-interactive context (CI, pipe, IDE terminal stub). Bubble Tea requires a real TTY.

**`make install` fails with `Permission denied` writing to `/usr/local/bin/`.**
Either run with `sudo make install`, or pick a user-writable prefix: `make install PREFIX=$HOME/.local` and add `$HOME/.local/bin` to your `$PATH`.

**The header underline looks broken or columns are misaligned.**
You're on an old build. Rebuild: `make install`. The fix is in `internal/tui/styles.go` — header padding must match `bubbles/table`'s default cell padding.

---

## Uninstall

Homebrew install:

```bash
brew uninstall get-ip
brew untap vianch/get-ip
make unalias && source ~/.zshrc   # if you opted into the alias
```

Makefile install:

```bash
make uninstall    # removes binary and alias block in one go
source ~/.zshrc
```

---

## Development

```bash
git clone https://github.com/vianch/get-ip.git
cd get-ip
go mod download
make run          # iterate
make test         # parser tests (no shellouts)
make fmt          # gofmt + go vet
```

Project layout:

```
get-ip/
├── main.go                      # entry point, --help / --version
├── Formula/get-ip.rb            # Homebrew formula
├── Makefile                     # build, install, alias, release helpers
├── scripts/release.sh           # tag + SHA256 helper
└── internal/
    ├── network/                 # detection layer (darwin + cross-platform stub)
    └── tui/                     # Bubble Tea model, styles, keybindings
```

The `internal/network` package is split with build tags. Add new platforms by adding a sibling file with the appropriate `//go:build <os>` tag implementing the same exported funcs.

Tests live next to the code they test (`darwin_test.go`) and parse fixture strings — no live shellouts in `go test`, so it's safe to run in CI.

---

## Releasing a new version

You need push access to `github.com/vianch/get-ip`. The flow:

```bash
# 1. Make sure you're on main, clean, and origin points at vianch/get-ip.
git checkout main
git pull
git status

# 2. Tag, push, and get the formula snippet.
make release V=v0.2.0

# 3. The script prints the new url + sha256. Paste them into Formula/get-ip.rb.
$EDITOR Formula/get-ip.rb

# 4. Commit and push the formula bump.
git add Formula/get-ip.rb
git commit -m "chore: bump formula to v0.2.0"
git push
```

Users then upgrade via:

```bash
brew update
brew upgrade get-ip
```

If you ever need to recompute the SHA for an existing tag:

```bash
make formula-sha V=v0.2.0
```

---

## License

MIT — see [LICENSE](./LICENSE).
