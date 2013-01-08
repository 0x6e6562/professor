# About

Professor is a Go library for the Cassandra CQL3 binary protocol.

# Limitations

Currently this is a synchronous implementation of CQL2 - i.e. it mimics the Thrift RPC semantics. Hopefully this limitation can be relaxed in due course.

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

# License

Copyright (c) 2013, Ben Hood
All rights reserved.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
* Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.