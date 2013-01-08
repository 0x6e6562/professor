# About

Professor is a Go library for the Cassandra CQL3 binary protocol.

# Running 

In order to test this out you will need to have a Cassandra instance running on the localhost.
Since the protocol was only introduced in version 1.2, you cannot use versions below this.
As of the current version (1.2), you need to tweak the default cassandra.yaml configuration by
uncommenting the following lines:

	start_native_transport: true
	native_transport_port: 9042

Once you have set the GOPATH, run

	$ go test professor

# Status

Currently this is just a proof of concept library and as such, has not been tested properly.
