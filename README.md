# Aqua

Easily run saved commands per project

https://github.com/user-attachments/assets/0e382d3f-4ed0-4c04-ae82-abf0e9f3606c

# How to get started

Create **lake.json** in project root with this schema

```json
[
  {
    "cmd": "go run .",
    "title": "Start application"
  }
]
```

# How project root is determined?
By git files. When aqua is started it will search up to 3 directories up for the **.git** files

# Warning
* Do not use this application to run TUI apps, or any other app that require user input during its lifetime<br>
It's best to use aqua for commands that will only spew out logs, so for example web backends

