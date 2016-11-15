jamesd
======
A universal packetmanager for heterogeneous fleets of machines with server side version control.

# Concept
Version control in distributed and/or heterogeneous systems is a pain. Traditionally you install your software once, and hope that you will never have to change it.
If you need to change it, you put together a packet and distribute it as an update packet which needs to be installed on every system. You can multiply this process by the number of different architectures you maintain. At least on debian based systems you could also upload your packet to your packet-server and run `apt-get upgrade` on each system. Never the less, You need to interact with every device you want to upgrade.

Jamesd solves this.

It provides a simple client software which runs in the background on your devices. These software periodically asks the server which software should be installed and acts accordingly. You can then specify the desired state of your devices on the server and lean back while letting jamesd updating all your devices.

# Parts
* jamesd
  * this runs on the repository server
  * responsible for:
    * managing packets
    * managing device states
    * installing / uninstalling packets from devices
* jamesc
  * this is the client which runs on the devices
  * is a longrunning daemon
  * periodically sends its state (installed packets) to jamesd
  * accepts install/uninstall instructions from jamesd
* jamesd-ctl
  * commandline tool to interact with jamesd
  * provides interface to:
    * upload new packages
    * delete packages
    * get list of connected devices
    * get state information about devices
    * set/get the desired state of the devices

# Tags
To find matching packets for each device, jamesd uses tags for devices and packets.
Jamesc's state information contains a so-called system-tag-list. These tags describe the system. This list could look like this: `['armv7l', 'glibc-2.23', 'systemd']`
Each packet has a name and also a tag-list. The meta-information for a packet could look like this:
```js
{
  name: 'test-packet',
  tags: ['v1.0.0', 'armv7l', 'glibc-2.23']
}
```
If you instruct the system that `test-packet` with version `v1.0.0` should be installed on the system with the given system-tag-list, jamesd will look for a packet `test-packet` which fullfills the combined taglist of the specification and system-tag-list. A packet fullfills the requirements when all tags of the packet, are part of the combined taglist.

# Packet-Management
You can use jamesd-ctl to manage your packets. To create a new packet you must create a compressed tar archive containing all the files you need. The folder structure must represent the files location on the target device.
e.g:
```
openvpn-root
├── etc
│   ├── openvpn
│   │   └── openvpn.conf
│   └── systemd
│       └── system
│           ├── multi-user.target.wants
│           │   └── openvpn.service -> /etc/systemd/system/openvpn.service
│           └── openvpn.service
└── usr
    └── bin
       └── openvpn
```

Create an archive like this:
```bash
> tar cfvJ openvpn.tar.xz -C openvpn-root/ .
```
You can add pre/post install/uninstall scripts:
```bash
> touch preinst.sh postinst.sh prerm.sh postrm.sh
```
Now you are ready to upload your packet to jamesd:
```bash
> jamesd-ctl \
    --cmd add-packet \
    --name openvpn \
    --data openvpn.tar.xz \
    --tags v1.0.0,armv7l,systemd,glibc-2.23 \
    --preinst preinst.sh \
    --prerm prerm.sh \
    --postinst postinst.sh \
    --postrm postrm.sh
```

# Device Management
Lets assume there is a device with the id `device-1`. To get the state of it do this:
```
> jamesd-ctl --cmd get-state --id device-1
```
This could output something like this:
```yaml
id: device-1
systemtags:
- armv7l
- glibc-2.23
- systemd
apps:
- name: openvpn
  tags:
  - v0.9.9
```
Now create a file `state.yaml` with the state as content, but adjust the version to this:
```yaml
id: device-1
systemtags:
- armv7l
- glibc-2.23
- systemd
apps:
- name: openvpn
  tags:
  - v1.0.0
```
Now set the desired state to your adjusted configuration:
```
> jamesd-ctl --cmd set-desired-state --file state.yaml
```
The next time when `device-1` sends it state, jamesd will recognize the version mismatch, and will start with uninstalling the old version (executing prerm, deleting all packet files, executing postrm) and continue with installing the new version (executing preinst, unpacking all packet files, executing postinst)

# Setup Jamesd
Installing jamesd is simple.
The only requierement is a running mongodb server.
```
go get github.com/trusch/jamesd/jamesd
go get github.com/trusch/jamesd/jamesd-ctl
sudo -E cp $GOPATH/bin/jamesd $GOPATH/bin/jamesd-ctl /usr/local/bin
```
To manage jamesd as a systemd service you can use [jamesd.service](./jamesd.service)
```
> sudo -E cp $GOPATH/src/github.com/trusch/jamesd/jamesd.service /etc/systemd/systemd/
> sudo adduser jamesd
> sudo systemctl enable jamesd.service
> sudo systemctl start jamesd.service
```

# Setup Jamesc
If you are on a 'normal' architecture (eg. amd64, x86) just install via go get
```
go get github.com/trusch/jamesd/jamesc
sudo -E cp $GOPATH/bin/jamesc /usr/local/bin
```
To manage it via systemd you can use [jamesc.service](./jamesc.service) as template.

## Statically crosscompile for arm
first 'go get' it like described above. Then execute the following:
```
> export CC=arm-linux-gnueabihf-gcc
> export GOOS=linux
> export GOARCH=arm
> export CGO_ENABLED=1
> go install -v -ldflags '-linkmode external -extldflags -static' \
    github.com/trusch/jamesd/jamesc
```
Now you can copy the binary (`$GOPATH/bin/linux_arm/jamesc`) to any (you know what I mean) arm device.

# Key management
You need a proper public key infrastructure to work with jamesd. You can use the openssl client binary to maintain a CA yourself, or you use a simple wrapper like [pkitool](https://github.com/trusch/vpntool/tree/master/pki/pkitool)
## Install pkitool
```
go get github.com/trusch/vpntool/pki/pkitool
```
## Setup pki
```
pkitool --init
pkitool --add-server jamesd
pkitool --add-client device-1
pkitool --add-client device-2
pkitool --add-client device-3
```
Keep the created directory safe! As you can see in the unit files, jamesd and jamesc need their own keys/certificates and the CA certificate. The specified names are the device id's reported by jamesd.
