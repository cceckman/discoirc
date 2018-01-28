package testhelper

import (
	"bytes"
	"io/ioutil"
	"fmt"
)


var configFile string

func init() {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		panic(fmt.Sprintf("could not init oragno.yaml: %s", err))
	}
	configFile = f.Name()
	_, err = config.WriteTo(f)
	if err != nil {
		panic(fmt.Sprintf("could not write oragno.yaml: %s", err))
	}
}


var config = bytes.NewBufferString(`
# oragono IRCd config

# network configuration
network:
    # name of the network
    name: OragonoTest

# server configuration
server:
    # server name
    name: oragono.test

    # addresses to listen on; local, SSL only
    listen:
        - "[::1]:6697"
        - "127.0.0.1:6697"

    # tls listeners
    tls-listeners:
        # listener on ":6697"
        ":6697":
            key: tls.key
            cert: tls.crt

    # strict transport security, to get clients to automagically use TLS
    sts:
        # whether to advertise STS
        #
        # to stop advertising STS, leave this enabled and set 'duration' below to "0". this will
        # advertise to connecting users that the STS policy they have saved is no longer valid
        enabled: false

        # how long clients should be forced to use TLS for.
        # setting this to a too-long time will mean bad things if you later remove your TLS.
        # the default duration below is 1 month, 2 days and 5 minutes.
        duration: 1mo2d5m

        # tls port - you should be listening on this port above
        port: 6697

        # should clients include this STS policy when they ship their inbuilt preload lists?
        preload: false

    # use ident protocol to get usernames
    check-ident: true

    # password to login to the server
    # generated using  "oragono genpasswd"
    #password: ""

    # motd filename
    # if you change the motd, you should move it to ircd.motd
    motd: oragono.motd

    # motd formatting codes
    # if this is true, the motd is escaped using formatting codes like $c, $b, and $i
    #motd-formatting: true

    # addresses/hostnames the PROXY command can be used from
    # this should be restricted to 127.0.0.1 and localhost at most
    # you should also add these addresses to the connection limits and throttling exemption lists
    proxy-allowed-from:
        # - localhost
        # - "127.0.0.1"

    # controls the use of the WEBIRC command (by IRC<->web interfaces, bouncers and similar)
    # webirc:
        # one webirc block -- should correspond to one set of gateways
        -
            # tls fingerprint the gateway must connect with to use this webirc block
            # fingerprint: 938dd33f4b76dcaf7ce5eb25c852369cb4b8fb47ba22fc235aa29c6623a5f182

            # password the gateway uses to connect, made with  oragono genpasswd
            # password: JDJhJDA0JG9rTTVERlNRa0hpOEZpNkhjZE95SU9Da1BseFdlcWtOTEQxNEFERVlqbEZNTkdhOVlYUkMu

            # hosts that can use this webirc command
            # hosts:
                # - localhost
                # - "127.0.0.1"
                # - "0::1"

    # maximum length of clients' sendQ in bytes
    # this should be big enough to hold /LIST and HELP replies
    max-sendq: 16k

    # maximum number of connections per subnet
    connection-limits:
        # whether to enforce connection limits or not
        enabled: true

        # how wide the cidr should be for IPv4
        cidr-len-ipv4: 32

        # how wide the cidr should be for IPv6
        cidr-len-ipv6: 64

        # maximum number of IPs per subnet (defined above by the cird length)
        ips-per-subnet: 16

        # IPs/networks which are exempted from connection limits
        exempted:
            - "127.0.0.1"
            - "127.0.0.1/8"
            - "::1/128"

    # automated connection throttling
    connection-throttling:
        # whether to throttle connections or not
        enabled: true

        # how wide the cidr should be for IPv4
        cidr-len-ipv4: 32

        # how wide the cidr should be for IPv6
        cidr-len-ipv6: 64

        # how long to keep track of connections for
        duration: 10m

        # maximum number of connections, per subnet, within the given duration
        max-connections: 32

        # how long to ban offenders for, and the message to use
        # after banning them, the number of connections is reset (which lets you use UNDLINE to unban people)
        ban-duration: 10m
        ban-message: You have attempted to connect too many times within a short duration. Wait a while, and you will be able to connect.

        # IPs/networks which are exempted from connection limits
        exempted:
            - "127.0.0.1"
            - "127.0.0.1/8"
            - "::1/128"

# account options
accounts:
    # account registration
    registration:
        # can users register new accounts?
        enabled: true

        # length of time a user has to verify their account before it can be re-registered
        # default is 120 hours, or 5 days
        verify-timeout: "120h"

        # callbacks to allow
        enabled-callbacks:
            - none # no verification needed, will instantly register successfully

        # allow multiple account registrations per connection
        # this is for testing purposes and shouldn't be allowed on real networks
        allow-multiple-per-connection: false

    # is account authentication enabled?
    authentication-enabled: true

# channel options
channels:
    # modes that are set when new channels are created
    # +n is no-external-messages and +t is op-only-topic
    # see  /QUOTE HELP cmodes  for more channel modes
    default-modes: +nt

    # channel registration - requires an account
    registration:
        # can users register new channels?
        enabled: true

# operator classes
oper-classes:
    # local operator
    "local-oper":
        # title shown in WHOIS
        title: Local Operator

        # capability names
        capabilities:
            - "oper:local_kill"
            - "oper:local_ban"
            - "oper:local_unban"

    # network operator
    "network-oper":
        # title shown in WHOIS
        title: Network Operator

        # oper class this extends from
        extends: "local-oper"

        # capability names
        capabilities:
            - "oper:remote_kill"
            - "oper:remote_ban"
            - "oper:remote_unban"

    # server admin
    "server-admin":
        # title shown in WHOIS
        title: Server Admin

        # oper class this extends from
        extends: "local-oper"

        # capability names
        capabilities:
            - "oper:rehash"
            - "oper:die"
            - "samode"

# ircd operators
opers:
    # operator named 'dan'
    dan:
        # which capabilities this oper has access to
        class: "server-admin"

        # custom whois line
        whois-line: is a cool dude

        # custom hostname
        vhost: "n"

        # modes are the modes to auto-set upon opering-up
        modes: +is acjknoqtux

        # password to login with /OPER command
        # generated using  "oragono genpasswd"
        password: JDJhJDA0JE1vZmwxZC9YTXBhZ3RWT2xBbkNwZnV3R2N6VFUwQUI0RUJRVXRBRHliZVVoa0VYMnlIaGsu

# logging, takes inspiration from Insp
logging:
    -   # how to log these messages
        #
        #   file    log to given target filename
        #   stdout  log to stdout
        #   stderr  log to stderr
        method: stderr
        # filename to log to, if file method is selected
        filename: ircd.log
        # type(s) of logs to keep here. you can use - to exclude those types
        #
        # exclusions take precedent over inclusions, so if you exclude a type it will NEVER
        # be logged, even if you explicitly include it
        #
        # useful types include:
        #   *               everything (usually used with exclusing some types below)
        #   accounts        account registration and authentication
        #   channels        channel creation and operations
        #   commands        command calling and operations
        #   opers           oper actions, authentication, etc
        #   password        password hashing and comparing
        #   userinput       raw lines sent by users
        #   useroutput      raw lines sent to users
        type: "* "

        # one of: debug info warn error
        level: debug

# debug options
debug:
    # when enabled, oragono will attempt to recover from certain kinds of
    # client-triggered runtime errors that would normally crash the server.
    # this makes the server more resilient to DoS, but could result in incorrect
    # behavior. deployments that would prefer to "start from scratch", e.g., by
    # letting the process crash and auto-restarting it with systemd, can set
    # this to false.
    recover-from-errors: false

    # enabling StackImpact profiling
    stackimpact:
        # whether to use StackImpact
        enabled: false

        # the AgentKey to use
        agent-key: examplekeyhere

        # the app name to report
        app-name: Oragono

# datastore configuration
datastore:
    # path to the datastore
    path: ircd.db

# limits - these need to be the same across the network
limits:
    # nicklen is the max nick length allowed
    nicklen: 32

    # channellen is the max channel length allowed
    channellen: 64

    # awaylen is the maximum length of an away message
    awaylen: 500

    # kicklen is the maximum length of a kick message
    kicklen: 1000

    # topiclen is the maximum length of a channel topic
    topiclen: 1000

    # maximum number of monitor entries a client can have
    monitor-entries: 100

    # whowas entries to store
    whowas-entries: 100

    # maximum length of channel lists (beI modes)
    chan-list-modes: 60

    # maximum length of IRC lines
    # this should generally be 1024-2048, and will only apply when negotiated by clients
    linelen:
        # tags section
        tags: 2048

        # rest of the message
        rest: 2048
`)


