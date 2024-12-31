package client

import (
	"bufio"
	"fmt"
	"io"
)

// Transport agnostic Client for the NSD server's control socket.
// This is *not* thread-safe, it's the consumers responsibility to protect the Client from concurrent use.
type Client struct {
	// Server-side command parsing logic: https://github.com/NLnetLabs/nsd/blob/149049ca0a8e5536d2cfe60461b9f74d4f8ccc02/remote.c#L2606

	socket  io.ReadWriteCloser
	scanner *bufio.Scanner
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

func (c *Client) readReply() (string, error) {
	if !c.scanner.Scan() {
		return "", c.scanner.Err()
	}
	return c.scanner.Text(), nil
}

func (c *Client) Close() error {
	return c.socket.Close()
}

// Stop request that NSD daemon stops
func (c *Client) Stop() error {
	if err := c.sendCmd(cmdStop); err != nil {
		return err
	}

	if reply, err := c.readReply(); err != nil {
		return err
	} else if reply != replyOK {
		return fmt.Errorf("unexpected reply: %s", reply)
	}

	return nil
}

// Reload causes NSD to reload modified zone files from disk
func (c *Client) Reload(zones []string) {
	//TODO implement me
	panic("implement me")
}

// Repattern reloads the config file.
// Alias of reconfig, https://github.com/NLnetLabs/nsd/blob/149049ca0a8e5536d2cfe60461b9f74d4f8ccc02/remote.c#L2640-L2643
// proto: repattern
func (c *Client) Repattern() {
	//TODO implement me
	panic("implement me")
}

// Reopen logfile (for log rotate
func (c *Client) LogReopen() {
	//TODO implement me
	panic("implement me")
}

// Status of server
func (c *Client) Status() {
	//TODO implement me
	panic("implement me")
}

func (c *Client) Stats() {
	//TODO implement me
	panic("implement me")
}

func (c *Client) StatsNoReset() {
	//TODO implement me
	panic("implement me")
}

func (c *Client) AddZone(name string, pattern string) {
	//TODO implement me
	panic("implement me")
}

func (c *Client) DelZone(name string) {
	//TODO implement me
	panic("implement me")
}

func (c *Client) ChangeZone(name string, pattern string) {
	//TODO implement me
	panic("implement me")
}

func (c *Client) Write(zones []string) {
	//TODO implement me
	panic("implement me")
}

func (c *Client) Notify(zones []string) {
	//TODO implement me
	panic("implement me")
}

func (c *Client) Transfer(zones []string) {
	//TODO implement me
	panic("implement me")
}

func (c *Client) ForceTransfer(zones []string) {
	//TODO implement me
	panic("implement me")
}

func (c *Client) ZoneStatus(zones []string) {
	//TODO implement me
	panic("implement me")
}

func (c *Client) ServerPID() {
	//TODO implement me
	panic("implement me")
}

func (c *Client) Verbosity(verbosity int) {
	//TODO implement me
	panic("implement me")
}

func (c *Client) GetTSig(keyName []string) {
	//TODO implement me
	panic("implement me")
}

func (c *Client) UpdateTSig(name string, secret string) {
	//TODO implement me
	panic("implement me")
}

func (c *Client) AddTSig(name string, secret string, algo *string) {
	//TODO implement me
	panic("implement me")
}

func (c *Client) AssocTSig(zone string, keyName string) {
	//TODO implement me
	panic("implement me")
}

func (c *Client) DelTSig(keyName string) {
	//TODO implement me
	panic("implement me")
}

func (c *Client) AddCookieSecret(secret string) {
	//TODO implement me
	panic("implement me")
}

func (c *Client) DropCookieSecret() {
	//TODO implement me
	panic("implement me")
}

func (c *Client) SctivateCookieSecret() {
	//TODO implement me
	panic("implement me")
}

func (c *Client) GetCookieSecrets() {
	//TODO implement me
	panic("implement me")
}
