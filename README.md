# Qwertify
A simple yet secure command-line password manager written in Go.

![screenshot](screenshot.png)
## Features
- **Offline Storage:** All data is stored locally in `~/.qwertify/safe.json`.
- **Strong Encryption:** Uses AES-GCM for data encryption and Scrypt for key derivation.
- **Master Password:** Protected by bcrypt-hashed master password.
- **Clipboard Support:** Automatically copies passwords to clipboard on retrieval.
- **Fast & Lightweight:** Built with Go for performance and simplicity.
- **Safe writing:** File lock to prevent race condition.
## Install
**Go** is required!
```
go install github.com/temaelkin/qwertify/cmd/qwfy@latest
```
## Usage
```
# Make a good master password!
# You will need it more than once

1. qwfy init           Initialize a new safe
2. qwfy add <url>      Add a new entry
3. qwfy get <url>      Retrieve an entry
4. qwfy edit <url>     Edit an existing entry
5. qwfy all            List all entries
6. qwfy help           Show help
```
---
## License
[MIT](https://github.com/temaelkin/qwertify/blob/main/LICENSE)
