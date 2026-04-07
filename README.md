# qwertify
A simple yet secure command-line password manager written in Go.

![screenshot](screenshot.png)
## Features
- All data is stored locally in `~/.qwertify/safe.json`.
- Uses AES-GCM for data encryption and Scrypt for key derivation.
- Protected by bcrypt-hashed master password.
- Automatically copies passwords to clipboard on retrieval.
- File lock to prevent race condition.
## Install
**Go 1.24.4 or later** is required!
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
5. qwfy del <url>      Delete an entry
6. qwfy all            List all entries
7. qwfy help           Show help
```
## Work In Progress
I work on qwertify when I have some free time and right mood so the process is pretty slow but I do have some plans for it so one day it becomes a convenient tool for managing your passwords :)
```
1. Tink AEAD encryption for better security and UX
2. TUI (tview or bubbletea)
3. Password generation
```
## WARNING
I really do not recommend using qwertify as your main password manager RIGHT NOW since the current state is unfinished and requires a lot of work to be done.\
But feel free to try it anyway.

---
## License
[MIT](https://github.com/temaelkin/qwertify/blob/main/LICENSE)
