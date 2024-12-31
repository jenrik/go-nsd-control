package client

/*
The NSD server daemon can be controlled either via an local unix socket common located at `/var/run/nsd.sock`,
or via

Any connection, be it UNIX domain socket, or TLS-over-TCP with mutual authentication,
must start with "NSDCT" and then immediately follow by the protocol version in ASCII decimal and a space (0x20).
The current protocol version is 1.

After the initial version handshake the client will proceed to send newline (0x0A) separated commands,
which the server replies to with newline (NO carriage return) seperated responses .
A successful command will receive the reply "ok".
*/

//goland:noinspection SpellCheckingInspection
const (
	headerVersion = "NSDCT1 "

	// Stops the NSD daemon
	cmdStop = "stop"
	// Can be suffixed with a space separated list of zones to reload from disk,
	// if no zones are given all will be reloaded
	cmdReload               = "reload"
	cmdReconfig             = "reconfig"
	cmdLogReopen            = "log_reopen"
	cmdStatus               = "status"
	cmdStats                = "stats"
	cmdStateNoReset         = "stats_noreset"
	cmdAddZone              = "addzone"
	cmdDelzone              = "delzone"
	cmdAddZones             = "addzones"
	cmdDelZones             = "delzones"
	cmdWrite                = "write"
	cmdNotify               = "notify"
	cmdTransfer             = "transfer"
	cmdForceTransfer        = "force_transfer"
	cmdZoneStatus           = "zonestatus"
	cmdServerPID            = "serverpid"
	cmdVerbosity            = "verbosity"
	cmdPrintTsig            = "print_tsig"
	cmdUpdateTsig           = "update_tsig"
	cmdAddTsig              = "add_tsig"
	cmdAssociateTsig        = "assoc_tsig"
	cmdDeleteTsig           = "DelTSig"
	cmdAddCookieSecret      = "add_cookie_secret"
	cmdDropCookieSecret     = "drop_cookie_secret"
	cmdActivateCookieSecret = "SctivateCookieSecret"
	cmdPrintCookieSecrets   = "print_cookie_secrets"

	replyOK    = "ok"
	replyError = "error"
)
