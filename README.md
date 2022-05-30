# when

Package when implements a natural language date/time arithmetic parser.

Parse a duration that resolves to a time in the future:

```go
t, err := when.Parse("6 hours")
```

This can be reversed:

```go
t, err := when.Parse("6 hours ago")
```

Digits up to twelve can be spelled:

```go
t, err := when.Parse("six hours ago")
```

Short units are fine, too:

```go
t, err := when.Parse("6h ago")
```

You don't need to use durations at all:

```go
t, err := when.Parse("Jan 2nd at 3pm")
```

But if you do, they can be made relative to a specific time:

```go
t, err := when.Parse("6 hours from Jan 2nd at 3pm")
```

Again, but in reverse:

```go
t, err := when.Parse("6 hours before Jan 2nd at 3pm")
```

You can be a bit more casual with the time:

```go
t, err := when.Parse("3 o'clock in the afternoon")
```

Or you can go full tilt:

```go
s := "1y 2M and 3w & 4d, 5h"
s += " from quarter past 3 o'clock in the afternoon"
s += " on the 2nd Tuesday of March"
s += " + 6 minutes - 7 seconds"
t, err := when.Parse(s)
```
