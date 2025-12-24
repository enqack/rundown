# Profiling Rundown

This guide explains how to profile the Rundown application itself to identify performance bottlenecks.

## Enabling Profiling

Start Rundown with the `--profile` flag to enable the pprof HTTP server:

```bash
rundown --profile
```

By default, the pprof server listens on `localhost:6060`. You can change the port:

```bash
rundown --profile --profile-port 8080
```

## Capturing Profiles

### CPU Profile

Capture a 30-second CPU profile:

```bash
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
```

This will download the profile and open an interactive pprof session.

### Memory Profile

Capture a heap memory profile:

```bash
go tool pprof http://localhost:6060/debug/pprof/heap
```

### Goroutine Profile

View active goroutines:

```bash
go tool pprof http://localhost:6060/debug/pprof/goroutine
```

## Analyzing Profiles

Once in the pprof interactive session, useful commands include:

- `top` - Show top functions by CPU/memory usage
- `list <function>` - Show source code for a function
- `web` - Generate a visual graph (requires graphviz)
- `pdf` - Generate a PDF report
- `help` - Show all available commands

### Example Analysis

```bash
# Capture CPU profile
$ go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# In pprof session:
(pprof) top10
(pprof) list updateCpuViewport
(pprof) web
```

## Web Interface

You can also browse profiles via web browser:

```
http://localhost:6060/debug/pprof/
```

This provides links to all available profiles and a flame graph visualization.

## Performance Tips

Based on profiling, here are some optimizations already implemented:

1. **Active Tab Only Updates**: Viewports only update for the currently visible tab
2. **Configurable Update Interval**: Use `+/-` keys to adjust refresh rate
3. **Efficient String Building**: Minimized allocations in view rendering

## Common Issues

### High CPU Usage

If you see high CPU usage:

1. Check the update interval (default 1s) - increase with `+` key
2. Profile to identify hot functions
3. Consider which tab you're viewing (Process tab is most expensive)

### Memory Growth

If memory usage grows over time:

1. Capture heap profiles at different times
2. Compare with `go tool pprof -base=<old> <new>`
3. Look for goroutine leaks with goroutine profile
