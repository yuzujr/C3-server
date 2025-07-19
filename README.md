# C3 Server
C3 Server is a Go-based backend for the C3 framework, offering better scalability, easier deployment, and higher reliability compared to previous versions.

## Under Development

The server has been refactored and now fully replicates the core features of the original Node.js version. The project is moving into active development.

### Outlook

- More new features upcoming.

## Features

- **Web frontend**: A pretty web interface that provides:
  - Multi-user management
  - Sending real-time commands
  - Viewing screenshots (supports multiple monitors)
  - A reverse shell based on xterm.js (frontend), WebSocket (communication), and pty (backend), offering millisecond-level responsiveness and a native experience

- **WebSocket Communication**: Enables real-time interaction with clients.
- **Database Integration**: Uses PostgreSQL for persistent storage.
- **Environment Configuration**: Uses `.env` files for easy configuration management.

## Quick Start

1. **Get Server**: Download from [Releases](https://github.com/yuzujr/C3-server/releases) or [Build from Source](docs/build.md)
2. **Setup**: Follow the [Setup Guide](docs/setup.md)

## Legal

**⚠️ DISCLAIMER: Educational and Research Use Only**

This software is provided for educational, research, and authorized system administration purposes only. Users are solely responsible for ensuring compliance with all applicable laws and regulations in their jurisdiction.

**Prohibited Uses:**
- Unauthorized access to computer systems
- Any illegal or malicious activities
- Violation of computer crime laws

**User Responsibilities:**
- Obtain proper authorization before deployment
- Comply with all applicable laws and regulations
- Use only on systems you own or have explicit permission to access
- Respect privacy and data protection laws

The developers assume no liability for misuse of this software.
