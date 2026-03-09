# perf-agent

I'm learning eBPF. This is the project I'm using to do it.

The idea is simple: build a Kubernetes agent that uses eBPF to observe what happens at the kernel level inside a cluster. Syscalls, network activity, process execution: the kind of things you can't see from application metrics alone.

I'm not an expert on this. I work on distributed systems and I've been curious about eBPF for a while. At some point curiosity has to become code, so here we are.

⚠️ This is experimental. Don't run it anywhere that matters.

---

## What's here now

A working skeleton. The agent starts, exposes `/healthz`, `/readyz` and `/metrics` over HTTP, and shuts down cleanly on SIGTERM. It builds into a small container image and installs on Kubernetes via Helm.

No eBPF yet. But the scaffolding is solid and the deployment model works. I've tested it on a real cluster. That was the point of this phase.

---

## Where this is going

The next step is to integrate [cilium/ebpf](https://github.com/cilium/ebpf) and load the first probe. I want to start with something simple (syscall counters, maybe basic network activity) before touching anything more complex.

After that, the direction depends on what I learn. I have some ideas around DaemonSet deployment with proper Linux capabilities, Prometheus exposition, maybe CRD-based configuration, but I won't commit to a roadmap I can't honestly estimate. I'll update this as things become clearer.

Whether this project ever gets there is an open question.

---

## Running it

```bash
# build
make all

# container image
make image

# deploy (image must be pushed first)
helm install perf-agent ./charts/perf-agent \
  --namespace <your-namespace>
```

---

## License

Apache 2.0
