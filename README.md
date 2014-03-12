This is a super simple TCP proxy.

It takes three arguments:
1. The local TCP port to listen on.
1. The remote hostname or IP address to connect to.
1. The remote TCP port to connect to.

The ports can be specified either numerically or by service name.

Once started, the proxy listens for inbound connections. When one arrives,
it connects to the remote machine and copies bytes back and forth between
the two connections until both ends hang-up.

Logging can be enabled by setting the LOG environment variable.

