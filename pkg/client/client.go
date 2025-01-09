package client

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// Transport agnostic Client for the NSD server's control socket.
// This is *not* thread-safe, it's the consumers responsibility to protect the Client from concurrent use.
type Client struct {
	// Server-side command parsing logic: https://github.com/NLnetLabs/nsd/blob/149049ca0a8e5536d2cfe60461b9f74d4f8ccc02/remote.c#L2606

	socket  io.ReadWriteCloser
	scanner *bufio.Scanner
}

type replyReader interface {
	readReply() ([]string, error)
}

func (c *Client) init() (err error) {
	c.scanner = bufio.NewScanner(c.socket)
	c.scanner.Split(bufio.ScanLines)
	_, err = c.socket.Write([]byte(headerVersion))
	return
}

func (c *Client) sendCmd(line string) error {
	_, err := c.socket.Write([]byte(line + "\n"))
	return err
}

func (c *Client) readLine() (string, error) {
	if !c.scanner.Scan() {
		return "", c.scanner.Err()
	}
	return c.scanner.Text(), nil
}

func (c *Client) readReply() ([]string, error) {
	var lines []string
	for {
		line, err := c.readLine()
		if err != nil {
			return nil, err
		}
		if line == "" {
			break
		}
		lines = append(lines, line)
	}
	return lines, nil
}

func expectOk(c replyReader) error {
	if reply, err := c.readReply(); err != nil {
		return err
	} else if len(reply) != 1 {
		return fmt.Errorf("unexpected reply: %s", reply)
	} else if len(reply) == 1 {
		if reply[0] == replyOK {
			return nil
		} else {
			return fmt.Errorf("server send error: %s", reply[0])
		}
	} else {
		return fmt.Errorf("unexpected reply: %+v", reply)
	}
}

func (c *Client) Close() error {
	return c.socket.Close()
}

// Stop request that NSD daemon stops
func (c *Client) Stop() error {
	// NSD handler: https://github.com/NLnetLabs/nsd/blob/NSD_4_11_0_REL/remote.c#L881
	if err := c.sendCmd(cmdStop); err != nil {
		return err
	}

	return expectOk(c)
}

// Reload causes NSD to reload modified zone files from disk
func (c *Client) Reload(zones []string) {
	// NSD handler: https://github.com/NLnetLabs/nsd/blob/NSD_4_11_0_REL/remote.c#L902
	//TODO implement me
	panic("implement me")
}

// Repattern reloads the config file.
// Alias of reconfig, https://github.com/NLnetLabs/nsd/blob/149049ca0a8e5536d2cfe60461b9f74d4f8ccc02/remote.c#L2640-L2643
// proto: repattern
func (c *Client) Repattern() {
	// NSD handler: https://github.com/NLnetLabs/nsd/blob/NSD_4_11_0_REL/remote.c#L2047
	//TODO implement me
	panic("implement me")
}

// Reopen logfile (for log rotate)
func (c *Client) LogReopen() error {
	// NSD handler: https://github.com/NLnetLabs/nsd/blob/NSD_4_11_0_REL/remote.c#L894
	if err := c.sendCmd(cmdLogReopen); err != nil {
		return err
	}

	return expectOk(c)
}

// Status of server
func (c *Client) Status() ([]string, error) {
	// NSD handler: https://github.com/NLnetLabs/nsd/blob/NSD_4_11_0_REL/remote.c#L1249
	if err := c.sendCmd(cmdStatus); err != nil {
		return nil, err
	}

	// Issue: There is no end of message
	// https://github.com/NLnetLabs/nsd/blob/NSD_4_11_0_REL/remote.c#L1249
	lines := make([]string, 0, 2)
	for {
		if reply, err := c.readLine(); err != nil {
			return nil, err
		} else {
			if reply == "" {
				return lines, nil
			}
			lines = append(lines, reply)
		}
	}
}

func (c *Client) Stats() ([]string, error) {
	// NSD handler: https://github.com/NLnetLabs/nsd/blob/NSD_4_11_0_REL/remote.c#L1266
	if err := c.sendCmd(cmdStats); err != nil {
		return nil, err
	}

	// Issue: There is no end of message
	// https://github.com/NLnetLabs/nsd/blob/NSD_4_11_0_REL/remote.c#L1249
	lines := make([]string, 0, 2)
	for {
		if reply, err := c.readLine(); err != nil {
			return nil, err
		} else {
			if reply == "" {
				return lines, nil
			}
			lines = append(lines, reply)
		}
	}
}

func (c *Client) StatsNoReset() ([]string, error) {
	// NSD handler: https://github.com/NLnetLabs/nsd/blob/NSD_4_11_0_REL/remote.c#L1266
	if err := c.sendCmd(cmdStatsNoReset); err != nil {
		return nil, err
	}

	// Issue: There is no end of message
	// https://github.com/NLnetLabs/nsd/blob/NSD_4_11_0_REL/remote.c#L1249
	lines := make([]string, 0, 2)
	for {
		if reply, err := c.readLine(); err != nil {
			return nil, err
		} else {
			if reply == "" {
				return lines, nil
			}
			lines = append(lines, reply)
		}
	}
}

func (c *Client) AddZone(domain string, pattern string) error {
	// NSD handler: https://github.com/NLnetLabs/nsd/blob/NSD_4_11_0_REL/remote.c#L1532
	cmd := fmt.Sprintf("%s %s %s", cmdAddZone, domain, pattern)
	if err := c.sendCmd(cmd); err != nil {
		return err
	}

	return expectOk(c)
}

func (c *Client) DelZone(domain string) error {
	// NSD handler: https://github.com/NLnetLabs/nsd/blob/NSD_4_11_0_REL/remote.c#L1541
	cmd := fmt.Sprintf("%s %s", cmdDelZone, domain)
	if err := c.sendCmd(cmd); err != nil {
		return err
	}

	return expectOk(c)
}

func (c *Client) ChangeZone(domain string, pattern string) error {
	// NSD handler: https://github.com/NLnetLabs/nsd/blob/NSD_4_11_0_REL/remote.c#L1550
	cmd := fmt.Sprintf("%s %s %s", cmdChangeZone, domain, pattern)
	if err := c.sendCmd(cmd); err != nil {
		return err
	}

	return expectOk(c)
}

func (c *Client) Write(zones []string) {
	// NSD handler: https://github.com/NLnetLabs/nsd/blob/NSD_4_11_0_REL/remote.c#L915
	//TODO implement me
	panic("implement me")
}

func (c *Client) Notify(zones []string) {
	// NSD handler: https://github.com/NLnetLabs/nsd/blob/NSD_4_11_0_REL/remote.c#L928
	//TODO implement me
	panic("implement me")
}

func (c *Client) Transfer(zones []string) {
	// NSD handler: https://github.com/NLnetLabs/nsd/blob/NSD_4_11_0_REL/remote.c#L953
	//TODO implement me
	panic("implement me")
}

func (c *Client) ForceTransfer(zones []string) {
	// NSD handler: https://github.com/NLnetLabs/nsd/blob/NSD_4_11_0_REL/remote.c#L994
	//TODO implement me
	panic("implement me")
}

func (c *Client) ZoneStatus(zone string) (*ZoneStatus, error) {
	// NSD handler: https://github.com/NLnetLabs/nsd/blob/NSD_4_11_0_REL/remote.c#L1033
	if err := c.sendCmd(cmdServerPID); err != nil {
		return nil, err
	}

	return parseZoneStatus(c)
}

func parseZoneStatus(c replyReader) (*ZoneStatus, error) {
	// NSD handler:
	status := &ZoneStatus{
		Attributes: make(map[string]string),
	}
	hasZone := false
	hasState := false

	reply, err := c.readReply()
	if err != nil {
		return nil, err
	}
	for _, line := range reply {
		if strings.HasPrefix(line, replyError) {
			return nil, fmt.Errorf("server send error: %s", reply)
		}

		match := commonKeyValueRegex.FindStringSubmatch(line)
		if match == nil {
			return nil, fmt.Errorf("unexpected reply: %s", line)
		}

		key := match[commonKeyValueRegex.SubexpIndex("key")]
		value := match[commonKeyValueRegex.SubexpIndex("value")]
		switch key {
		case "zone":
			status.Zone = value
			hasZone = true
		case "state":
			status.State = value
			hasState = true
		default:
			status.Attributes[key] = value
		}
	}
	if !hasZone || !hasState {
		return nil, fmt.Errorf("malformed reply, expected reply to contain zone name and state")
	}
	return status, nil
}

func (c *Client) ServerPID() (int, error) {
	// NSD handler: https://github.com/NLnetLabs/nsd/blob/NSD_4_11_0_REL/remote.c#L2130
	if err := c.sendCmd(cmdServerPID); err != nil {
		return -1, err
	}

	if reply, err := c.readLine(); err != nil {
		return -1, err
	} else {
		if v, err := strconv.Atoi(reply); err != nil {
			return -1, err
		} else {
			return v, nil
		}
	}
}

func (c *Client) Verbosity(verbosity int) error {
	// NSD handler: https://github.com/NLnetLabs/nsd/blob/NSD_4_11_0_REL/remote.c#L1187
	cmd := fmt.Sprintf("%s %d", cmdVerbosity, verbosity)
	if err := c.sendCmd(cmd); err != nil {
		return err
	}

	return expectOk(c)
}

func (c *Client) GetTSig(keyName []string) {
	// NSD handler: https://github.com/NLnetLabs/nsd/blob/NSD_4_11_0_REL/remote.c#L2137
	//TODO implement me
	panic("implement me")
}

func (c *Client) UpdateTSig(name string, secret string) {
	// NSD handler: https://github.com/NLnetLabs/nsd/blob/NSD_4_11_0_REL/remote.c#L2159
	//TODO implement me
	panic("implement me")
}

func (c *Client) AddTSig(name string, secret string, algo *string) {
	// NSD handler: https://github.com/NLnetLabs/nsd/blob/NSD_4_11_0_REL/remote.c#L2210
	//TODO implement me
	panic("implement me")
}

func (c *Client) AssocTSig(zone string, keyName string) {
	// NSD handler: https://github.com/NLnetLabs/nsd/blob/NSD_4_11_0_REL/remote.c#L2289
	//TODO implement me
	panic("implement me")
}

func (c *Client) DelTSig(keyName string) {
	// NSD handler: https://github.com/NLnetLabs/nsd/blob/NSD_4_11_0_REL/remote.c#L2348
	//TODO implement me
	panic("implement me")
}

func (c *Client) AddCookieSecret(secret string) error {
	// NSD handler: https://github.com/NLnetLabs/nsd/blob/NSD_4_11_0_REL/remote.c#L2500
	cmd := fmt.Sprintf("%s %s", cmdAddCookieSecret, secret)
	if err := c.sendCmd(cmd); err != nil {
		return err
	}

	return parseAddCookieSecretReply(c)
}

func parseAddCookieSecretReply(c replyReader) error {
	reply, err := c.readReply()
	if err != nil {
		return err
	}
	if len(reply) == 2 && reply[0] == "invalid cookie secret: invalid argument length" {
		return fmt.Errorf("invalid argument length")
	} else if len(reply) == 1 {
		if reply[0] == replyOK {
			return nil
		} else if strings.HasPrefix(reply[0], replyError) {
			return fmt.Errorf("server send error: %s", reply)
		} else {
			return fmt.Errorf("unexpected reply: %s", reply)
		}
	} else {
		return fmt.Errorf("unexpected reply: %s", reply)
	}
}

func (c *Client) DropCookieSecret() error {
	// NSD handler: https://github.com/NLnetLabs/nsd/blob/NSD_4_11_0_REL/remote.c#L2474
	if err := c.sendCmd(cmdDropCookieSecret); err != nil {
		return err
	}

	return expectOk(c)
}

func (c *Client) ActivateCookieSecret() error {
	// NSD handler: https://github.com/NLnetLabs/nsd/blob/NSD_4_11_0_REL/remote.c#L2448
	if err := c.sendCmd(cmdActivateCookieSecret); err != nil {
		return err
	}

	return expectOk(c)
}

var commonKeyValueRegex = regexp.MustCompile(`^\s*(?P<key>[^:\s]+)\s*:\s+"?(?P<value>[^"].*?)"?$`)

func (c *Client) GetCookieSecrets() (*CookieSecrets, error) {
	// NSD handler: https://github.com/NLnetLabs/nsd/blob/NSD_4_11_0_REL/remote.c#L2549
	if err := c.sendCmd(cmdPrintCookieSecrets); err != nil {
		return nil, err
	}

	return parseCookieSecretsReply(c)
}

func parseCookieSecretsReply(c replyReader) (*CookieSecrets, error) {
	cookieSecrets := &CookieSecrets{}
	reply, err := c.readReply()
	if err != nil {
		return nil, err
	}
	for _, line := range reply {
		if strings.HasPrefix(line, replyError) {
			return nil, fmt.Errorf("server send error: %s", reply)
		}

		match := commonKeyValueRegex.FindStringSubmatch(line)
		if match == nil {
			break
		}

		key := match[commonKeyValueRegex.SubexpIndex("key")]
		value := match[commonKeyValueRegex.SubexpIndex("value")]
		switch key {
		case "source":
			cookieSecrets.Source = value
		case "active":
			cookieSecrets.Active = value
		case "staging":
			cookieSecrets.Staging = &value
		}
	}
	return cookieSecrets, nil
}
