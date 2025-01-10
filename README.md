# 1 App 5 Stacks ported to Go+Templ+Datastar

The [original code](https://github.com/t3dotgg/1app5stacks ) that goes with [this video](https://www.youtube.com/watch?v=O-EWIlZW0mM) I think did Go dirty.

# Run

## With Docker

If you have Docker installed, it takes care of the setup for you in a Dockerized environment, which allows you to get it up and running quickly & easily.

The only requirement is that you have [Docker](https://www.docker.com/products/docker-desktop) installed (you do not need `golang`, `tmpl`, `sqlc` or any other project dependencies installed locally).

In terminal, `cd` to the project directory, and then type:

```
make dev
```

The first time you run this command, it may take a bit of time to build the Docker image, and download all of the appropriate packages, and cache them locally.

Then just navigate to `http://localhost:4321` in your browser, and the project will be up and running.

## Without Docker

If you have [Task](https://taskfile.dev/#/) and [Go 1.23.3](https://golang.org/) installed you can run the following:

```bash
task tools
task -w
```

# What's different?

1. 321LOC of Go across 6 files (including HTTPServer) and ~125LOC of [Templ](https://templ.guide/) across 3 files for UI.  Could be one of each but tried to match the spirit of the original.
   1. Templ could be shorter but expanded lines for readability.
   2. Pretty sure it's the smallest codebase of any of the apps in the original.
2. Metrics
   1. I'm pretty sure it's the fastest (queries take < 1ms on my machine)
   2. Smallest memory footprint (19MB sustained) of any of the apps in the original.
3. It handles speculation rules asynchronously for precaching the pokemon images.
4. Use pure Go SQLite so N+1 queries are not a problem.
5. It's a single binary with no dependencies.
6. Datastar
   1. One time cost of 12KiB of JS.
   2. Handles all UI interactions and backend communication.
   3. The version is [ALL the plugins](https://data-star.dev/bundler) we are only using a handful of them for this demo
   4. No websockets, just normal HTTP requests that can work on HTTP 2/3.

I hope this clear up why it's not
> [from a glance like Live View](https://x.com/theo/status/1858032204355612770)
>
> ~ [Theo](https://x.com/theo)