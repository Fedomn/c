# Commander Operator

[![Build Status](https://github.com/Fedomn/c/workflows/Go/badge.svg)](https://github.com/Fedomn/c/actions)
[![Build Status](https://github.com/Fedomn/c/workflows/Release/badge.svg)](https://github.com/Fedomn/c/actions)

Make command operation easier to use.

![Commander Operator](assets/normal-mode.png)
![Commander Operator](assets/search-mode.png)

# Features

* YAML format make configuration easier
* Terminal UI make operation faster
* Fuzzy Search make searching more convenient
* Including flexible normal mode and search mode

# Usage

configuration demo:

```yaml
-
 name: jump server
 cmd: ssh jump@server
-
 name: show date
 cmd: date
```

Terminal UI shortcuts in normal mode:

| key | operation in Normal Mode list |
| :--- | :--- |
| `j` / `<Down>` | Scroll Down |
| `k` / `<Up>` | Scroll Up |
| `<C-d>` | Scroll Half Page Down |
| `<C-u>` | Scroll Half Page Up |
| `<C-f>` | Scroll Page Down |
| `<C-b>` | Scroll Page Up |
| `q` / `<C-c>` / `<Escape>` | Close App |
| `/` | Into Search Mode |
| `Enter` | Select a command |
| `Backspace` | Delete the last letter of search string |


Terminal UI shortcuts in search mode:

| key | operation in Search Mode list |
| :--- | :--- |
| `<C-j>` / `<Down>` | Scroll Down |
| `<C-k>` / `<Up>` | Scroll Up |
| `<C-u>` | Erase search string |
| `<C-c>` / `<Escape>` | Back to Normal Mode |
| `Backspace` | Delete the last letter of search string |