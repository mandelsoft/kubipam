# Kubernetes based IP Address Management

This project provides an IP address management controller
based on Kubernetes resources.

There are two resources, an `IPAMRange` resource describing
managed IP ranges and an `IPAMRequest` resource for requesting
an IP CIDR or a single IP address.

### Ranges

The `IPAMRange` resource provides a pool of managed IP addresses.
It consists of a list of CIDRs, arbitrary IP ranges of the form `<ip>-<ip>` or
single IP addresses.

Optionally the range may provide a default `chunkSize`, the width of
network bits allocated for a request, if it does not request a dedicated
size.
 
A size lower than the IP size (32 for IPv4 and 128 for IPv6)
only makes sense for a list of CIDRs with a smaller network part.

```yaml
  apiVersion: ipam.mandelsoft.org/v1alpha1
  kind: IPAMRange
  metadata:
    name: mynetworkpool
    namespace: default
  spec:
    chunkSize: 24
    ranges:
      - 192.168.0.0/16
```

A range cannot be deleted as long as there are requests refering
to this range.

### Requests

The `IPAMRequest` resource is used to request the allocation
of a single IP or CIDR from a dedicated pool specified by an
`IPAMRange` object.

The allocated CIDR is reported in the status field of the
object.

```yaml
  apiVersion: ipam.mandelsoft.org/v1alpha1
  kind: IPAMRequest
  metadata:
    name: mynet
    namespace: default
  spec:
    ipam:
      name: mynetworkpool
  status:
    cidr: 192.168.1.0/24
    state: Ready
```

A request may specify a dedicated size (for example 32 for a dedicated IPv4
Address). If no size is given (either by the referenced pool or by the request
itself) a single IP address is allocated (/32 or /128).

By specifying a request CIDR or IP in the field `request` in the request object
it is possible to request the allocation of a dedicated range or IP.
If it could be granted the status is set accordinly as for an anonymous
request.

The allocation is released again, when the request object is deleted.


### Constraints

Once created the specification of a range or request MUST never
be modified.

So far there is no validating webhook yet, that prevents such operations.