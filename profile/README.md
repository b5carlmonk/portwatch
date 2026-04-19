# profile

The `profile` package manages named scan profiles for portwatch.

A **profile** groups together a set of hosts, ports, and protocol so you can
switch between common scanning configurations without retyping options.

## Usage

```go
store := profile.New()

_ = store.Add(profile.Profile{
    Name:     "web",
    Hosts:    []string{"192.168.1.1"},
    Ports:    []int{80, 443, 8080},
    Protocol: "tcp",
})

p, ok := store.Get("web")

_ = store.Save("/etc/portwatch/profiles.json")

loaded, _ := profile.Load("/etc/portwatch/profiles.json")
```

## Persistence

Profiles are stored as JSON. `Load` returns an empty store when the file does
not yet exist, making first-run safe.
