# Area

Area is a tiny CLI tool that gives you **named, persistent container environments** on top of Podman, without making you think much about containers.

The core idea is **persistence**, not security isolation.
A area is a named place where state lives. You create one, keep using it, and remove it when you are done.

```sh
area use myarea
```

That creates or enters the same persistent environment every time.

## What is this?

Area is:

- a thin wrapper around Podman
- based on **convention**, not configuration
- focused on **persistent environments**, host integration, and usability

It is useful when you want to keep some toolchain, application, or “foreign junk” in a separate persisted environment without building Dockerfiles or writing long Podman command lines yourself.

Examples:

- keep a messy SDK or dependency stack out of your host root filesystem
- keep an application containerized without maintaining a lot of config
- keep different persistent environments for different kinds of work

## What it is not

Area is **not**:

- a container orchestrator
- a hard security sandbox

Area may mount your real home directory, expose desktop/session resources, and mount the host filesystem into the container. That is intentional: the goal is a practical, persistent environment that still behaves like a normal desktop app environment.

If you want strong isolation, use a tool designed primarily for that.

## Install

It is simple:

```sh
go install github.com/grasparv/area/v2/cmd/area@latest
```

## Usage

```text
$ area

    use.......Create or enter a area
    run.......Run /bin/area-guest in a area
     ls.......List all areas
     rm.......Remove a area
```

### Enter or create a area

```sh
area use <name>
```

Enter the shell for the area. It creates the container if it does not exist yet, or starts it if it already exists but is stopped.

### Use the default program for a area

```sh
area run <name>
```

This runs:

```text
/usr/local/bin/area-guest
```

inside the area.

You define what that does.

This is meant for having a standard program that you can just `run`, while keeping that behavior as part of the persistent environment itself.

### List and remove areas

```sh
area ls
area rm <name>
```

`ls` shows all managed areas.
`rm` deletes the persistent container and its state.

## Podman inside containers

Area uses **Podman outside Podman**. That means `area` behaves the same **inside and outside** a area: containers are created as siblings on the host through the host Podman socket.

This is implemented by mounting the host Podman socket into the container at:

```text
/run/podman/podman.sock
```

The `area` binary is also refreshed into the guest so the guest-side wrapper stays in sync with the host binary.

## Host integration

When available, Area forwards host resources such as:

- host network
- X11
- Wayland
- audio (PipeWire/PulseAudio)
- GPU DRI

Area can display GUI apps on your desktop and use your graphics/audio stack, but it does not try to authenticate into or participate in your desktop session.

### Filesystem layout

Inside the container:

```text
/home/<your-username>   your real home from the host
/mnt/host               full host filesystem
```

Again, this is convenient, but it is also why Area should not be described as a hard sandbox.

## Design philosophy

Area is intentionally **opinionated**.

It does not try to expose every possible container option, and it does not aim to be a general-purpose interface to Podman or Docker. Those tools already exist-and they are very good at what they do.

Area exists to provide something different: a **simple, persistent environment you can enter and reuse**, without having to think about how containers work.

Area avoids adding flags, config files, or per-area settings. These are not missing features-they are deliberately excluded. If you need that level of control, you should use Podman directly.

Adding options would make Area more flexible-but also more complex, less predictable, and closer to being “just another container wrapper.”

Area chooses the opposite trade-off:
- fewer features
- stronger defaults
- a smaller mental model

The goal is that you can understand and use the entire tool in seconds, and then stop thinking about it.

If Area does not fit your needs, that is expected. It is designed to be small, focused, and opinionated.

