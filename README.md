JamesD
======

A universal packetmanager for heterogeneous fleets of machines with server side version control.

JamesD is a centralistic packet management system for heterogenous device fleets.
You can use it to manage multiple versions for each packet you maintain, and to automatically distribute updates.

## Problem
Imagine you have to maintain multiple fleets of devices, with many different architectures, operating systems, libc versions and so on. You may have ended fucked up like this because the devices your company is selling have a long support duration, or you're retrofitting many different old devices etc...

How do you organize and distribute software for all those devices?

You can not simply use common packet repositories, because they are designed to support one OS with at least support for diffent architectures. So how to deal with diffent OS or same OS and different base library versions? Another problem is that classical packet management solutions assume that the repository holds the packets, and the clients decide which packets they want to install. This is another no-go since our devices are smart-devices whicht are rather stupid when it comes to intelligent decisions like which packets should be installed. There is no operator who could decide this, sometimes there is not even physical access to the devices.

## Solution
JamesD is here to rescue you! It is completly agnostic about anything in your software packets, so you can supply packets for every possible scenario. The following core features support you:

* upload arbitary packets
* get packets based on device labels
* specify which packets should be installed
* powerfull command line interface to manage specs and packets
* client side daemon which polls for changes and installs/uninstalls accordingly

## Concept
### Parts

JamesD consist of three main parts:
* repository server (`jamesd`)
  * hosts packets
  * hosts specs
  * provides restfull api
* client-side daemon (`jamesc`)
  * asks server what should be installed
  * installs packets
  * uninstalls deprecated packets
* commandline tool (`jamesd-ctl`)
  * manage packets
  * manage specs
  * manually install / uninstall packets

So the repository server hosts your packets and specification about what should be installed on which machines and the client-daemon acts according to the specs.
To create or upload packets and specs you can use the commandline tool, or speak directly to the HTTP api.

### Labels
The matching which packets should be installed is done via labels. A `label` is a simple pair of two strings like (version, 1.0.0), (arch, armv7l) or (fleet, temp-sensors). You can specify labels as you like, as long as it makes sense for your usecase.

Labels are used to identify packets and devices.

To identify a single packet, you would need the `name` of the packet and a `labelset`.
A `labelset` is a collection of zero or more labels, where the key part of the has to be unique (so an invalid `labelset` would be { (arch, amd64), (arch, armhf) } ).
There can't be two packets with the same `name` and the same `labelset`.

If a device asks for packets, it tells its device `labelset`. This typically contains two types of labels:
* labels to identify the devices usecase
* labels about the devices requierements (architecture, libc version etc.)

### Specs
A spec is simply a definition of what packets should be installed on which devices.
Therefore a spec definition looks quiet simple:
```yaml
id: logger-spec
target:
  fleet: alpha
apps:
  - name: logger
    labels:
      version: 1.0.0
```
It contains:
  * an ID to identify the spec
  * a target `labelset` which says "this spec should be applied to all devices with (fleet, alpha) in its labelset"
  * a list of apps to be installed
    * this is not a specific packet!
    * it says: "give each device the best matching logger packet which contains (version, 1.0.0) in its labelset"

### Packet Matching
Assume we have the following server config in our repository:
```yaml
packets:
- name: logger
  labels:
    version: 1.0.0
    arch: amd64

- name: logger
  labels:
    version: 1.0.0
    arch: armv6l

- name: logger
  labels:
    version: 1.0.0
    arch: armv7l

specs:
- id: logger-spec
  target:
    fleet: alpha
  apps:
    - name: logger
      labels:
        version: 1.0.0

```

If a client asks for packets, he sends its device `labelset`, for example:
```yaml
fleet: alpha
arch: armv7l
```

At first the server searches for matching specs, and will find the logger spec. This spec tells him, that the logger packet should be installed in version 1.0.0.

To find the best packet, the server now merges the device labels with the logger labels from the spec, resulting in this new `labelset`:
```yaml
version: 1.0.0
fleet: alpha
arch: armv7l
```
Now the database is queried for a packet whichs `labelset` is a subset of this merged `labelset`. As a result the correct logger packet (with armv7l and 1.0.0) will be returned and the ID of it will be reported to the client.
