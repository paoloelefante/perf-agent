# perf-agent

I'm learning eBPF. This is the project I'm using to do it.

The idea is simple: build a Kubernetes agent that uses eBPF to observe what happens at the kernel level inside a cluster. Syscalls, network activity, process execution: the kind of things you can't see from application metrics alone.

I'm not an expert on this. I work on distributed systems and I've been curious about eBPF for a while. At some point curiosity has to become code, so here we are.

⚠️ This is experimental. Don't run it anywhere that matters.

---

## What's here now

A working skeleton. The agent starts, exposes `/healthz`, `/readyz` and `/metrics` over HTTP, and shuts down cleanly on SIGTERM. It builds into a small container image and installs on Kubernetes via Helm.

The agent now runs as a DaemonSet across all cluster nodes with eBPF capabilities
active (CAP_BPF, CAP_PERFMON, CAP_SYS_PTRACE), hostPID enabled, and /sys and
/sys/kernel/debug mounted. No eBPF programs yet, but the deployment model is ready.

No eBPF yet. But the scaffolding is solid and the deployment model works.
I've tested it on a real cluster. That was the point of this phase.

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

## Why the DaemonSet needs all those privileges

eBPF agents require specific kernel-level access that normal pods don't have. Here's what each setting does and why it's necessary.

### hostPID: true

Linux has PID namespaces that isolate process visibility. Normally a container only sees its own processes. With `hostPID: true`, the agent shares the host's PID namespace and sees all processes on the node.

This matters because eBPF probes report PIDs from the kernel's perspective (host PIDs). Without `hostPID`, the agent couldn't correlate events to processes - a kprobe returning `pid=12345` would be meaningless if that PID doesn't exist in the container's `/proc`.

Note: Kubernetes namespaces (like `default` or `kube-system`) are just logical groupings. They don't create kernel-level isolation. With `hostPID`, the agent sees processes from all K8s namespaces on that node.

### Volume mounts

**/sys** (read-only) - Contains `/sys/fs/bpf`, where the kernel exposes pinned BPF programs and maps. Needed to inspect loaded eBPF objects.

**/sys/kernel/debug** (read-write) - Contains the `tracing/` interface for kprobes, tracepoints, and ftrace. Some probe attachment and debugging operations require writing here.

### Security context

Instead of running as `privileged: true` (which grants full access), the agent requests only the specific Linux capabilities needed for eBPF:

| Capability | Purpose |
|------------|---------|
| CAP_BPF | Load eBPF programs into the kernel via `bpf()` syscall |
| CAP_PERFMON | Access perf events and hardware counters |
| CAP_SYS_PTRACE | Read memory/state of other processes (needed for uprobes and stack traces) |

Other settings:
- `drop: ALL` - Remove all default capabilities first (least privilege)
- `runAsUser: 0` - Root is required to load eBPF even with the capabilities
- `readOnlyRootFilesystem: true` - Container can't write to its own filesystem
- `allowPrivilegeEscalation: false` - Process can't gain additional privileges

### Tolerations

```yaml
tolerations:
  - operator: Exists
```

This means "schedule on any node regardless of taints". Nodes can have taints (like `node-role.kubernetes.io/control-plane`) that repel pods. A DaemonSet observability agent needs to run everywhere - including master nodes and nodes with issues - so it tolerates all taints.

### ClusterRole (pods/nodes watch)

eBPF gives you kernel-level data: PIDs, cgroup IDs, syscalls. But it doesn't know Kubernetes concepts. To correlate "PID 12345 made syscall X" with "that's the Cassandra pod in namespace production", the agent needs to watch pods via the K8s API and build a mapping table.

Without this, you get raw data without K8s context.

---

## License

Apache 2.0
