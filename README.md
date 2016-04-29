Crestron Definitions:

- `Attributes` are the same thing as `Signals`
- `Signals` are the individual "variables" we can see in Roomview
- `Processors` area source where you can get `Symbols`
- `Rooms` contain one or more `Processors`
- `Symbols` are containers of `Signals`
- `Assets` are contained by `Rooms` (don't worry about this for now)

High-Pass Overview:

1. Autodiscover a room (with correct symbols, etc.)
1. Delete that room
1. Delete its attendant processor (TODO: Figure out how to programmatically delete all processors from Fusion)
1. Use this repo to pull attributes from ELK `crestron-fusion-add-host -op S`
