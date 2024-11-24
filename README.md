# 1 App 5 Stacks ported to Go+Templ+Datastar

The [original code](https://github.com/t3dotgg/1app5stacks ) that goes with [this video](https://www.youtube.com/watch?v=O-EWIlZW0mM) I think did Go dirty.

# Run

Haven't made a docker container yet but if you have [Task](https://taskfile.dev/#/) and [Go 1.23.3](https://golang.org/) installed you can run the following:

```bash
task tools
task -w
```

# What's different?

1. About 350 lines of Go code and 150 of [Templ](https://templ.guide/) for UI.  Templ could be shorter but expanded lines for readability.
2. I'm pretty sure it its the faster (queries take < 5ms on my machine) and the smaller memory footprint (24MB) of any of the apps in the original.
3. It's only a one time cost of 12KiB of JS
4. It handles speculation rules asynchronously for precaching the pokemon images.
5. Use pure Go SQLite so N+1 queries are not a problem.
6. It's a single binary with no dependencies.
7. Does lazy loading of results and works well on 3G testing
