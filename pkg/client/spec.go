package client

/*
The NSD server daemon can be controlled either via an local unix socket common located at `/var/run/nsd.sock`,
or via

Any connection, be it UNIX domain socket, or TLS-over-TCP with mutual authentication,
must start with "NSDCT" and then immediately follow by the protocol version in ASCII decimal and a space (0x20).
The current protocol version is 1.

After the initial version handshake the client will proceed to send newline (0x0A) separated commands,
which the server replies to with newline (NO carriage return) seperated responses .
A successful command will receive the reply "ok" or a relevant message.
Failed and invalid commands will receive a single-line response starting with "error", with some exceptions.
All replies will end with a blank line signaling end of message.

Commands can be found in the NSD source code: https://github.com/NLnetLabs/nsd/blob/NSD_4_11_0_REL/remote.c#L2602
*/

//goland:noinspection SpellCheckingInspection
const (
	headerVersion = "NSDCT1 "

	cmdActivateCookieSecret = "activate_cookie_secret"
	cmdAddCookieSecret      = "add_cookie_secret"
	cmdAddTsig              = "add_tsig"
	cmdAddZone              = "addzone"
	cmdAddZones             = "addzones"
	cmdAssociateTsig        = "assoc_tsig"
	cmdChangeZone           = "changezone"
	cmdDeleteTsig           = "del_tsig"
	cmdDelZone              = "delzone"
	cmdDelZones             = "delzones"
	cmdDropCookieSecret     = "drop_cookie_secret"
	cmdForceTransfer        = "force_transfer"
	cmdLogReopen            = "log_reopen"
	cmdNotify               = "notify"
	cmdPrintCookieSecrets   = "print_cookie_secrets"
	cmdPrintTsig            = "print_tsig"
	// Same as repattern
	cmdReconfig = "reconfig"
	// Can be suffixed with a space separated list of zones to reload from disk,
	// if no zones are given all will be reloaded
	cmdReload = "reload"
	// Reload the configuration file if possible and apply keys and pattern anew
	cmdRepattern    = "repattern"
	cmdServerPID    = "serverpid"
	cmdStats        = "stats"
	cmdStatsNoReset = "stats_noreset"
	cmdStatus       = "status"
	// Stops the NSD daemon
	cmdStop       = "stop"
	cmdTransfer   = "transfer"
	cmdUpdateTsig = "update_tsig"
	cmdVerbosity  = "verbosity"
	cmdWrite      = "write"
	cmdZoneStatus = "zonestatus"

	replyOK = "ok"
	// Error messages always start with this sequence
	replyError = "error"
)

type CookieSecrets struct {
	Source  string
	Active  string
	Staging *string
}

type ZoneStatus struct {
	Zone       string
	State      string
	Attributes map[string]string
}
