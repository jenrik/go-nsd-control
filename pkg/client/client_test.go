package client

import (
	"reflect"
	"strings"
	"testing"
)

type StaticReply []string

func (s StaticReply) readReply() ([]string, error) {
	return s, nil
}

func NewStaticReply(lines []string) StaticReply {
	return lines
}

func Test_parseCookieSecretsReply(t *testing.T) {
	type args struct {
		c replyReader
	}
	tests := []struct {
		name    string
		args    args
		want    *CookieSecrets
		wantErr bool
	}{
		{
			name: "no staged cookie secret",
			args: args{
				c: NewStaticReply([]string{
					"source : \"/var/db/nsd/cookiesecrets.txt\"",
					"active : cd0636b6a5f8b9b1004b2450155ffca1",
				}),
			},
			want: &CookieSecrets{
				Source:  "/var/db/nsd/cookiesecrets.txt",
				Active:  "cd0636b6a5f8b9b1004b2450155ffca1",
				Staging: nil,
			},
			wantErr: false,
		}, {
			name: "staged cookie secret",
			args: args{
				c: NewStaticReply([]string{
					"source : \"/var/db/nsd/cookiesecrets.txt\"",
					"active : 4f08019819f6b945e03b6e91aafa0e8e",
					"staging: cd0636b6a5f8b9b1004b2450155ffca1",
				}),
			},
			want: &CookieSecrets{
				Source:  "/var/db/nsd/cookiesecrets.txt",
				Active:  "4f08019819f6b945e03b6e91aafa0e8e",
				Staging: &[]string{"cd0636b6a5f8b9b1004b2450155ffca1"}[0],
			},
			wantErr: false,
		}, {
			name: "no manually configured cookie secret",
			args: args{
				c: NewStaticReply([]string{
					"source : random generated",
					"active : 8234dff32ace962428c8da3d22da0d49",
				}),
			},
			want: &CookieSecrets{
				Source:  "random generated",
				Active:  "8234dff32ace962428c8da3d22da0d49",
				Staging: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseCookieSecretsReply(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseCookieSecretsReply() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseCookieSecretsReply() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_expectOk(t *testing.T) {
	type args struct {
		c replyReader
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				c: NewStaticReply([]string{
					"ok",
				}),
			},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				c: NewStaticReply([]string{
					"error",
				}),
			},
			wantErr: true,
		},
		{
			name: "malformed",
			args: args{
				c: NewStaticReply([]string{
					"asdf",
				}),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := expectOk(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("expectOk() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_parseZoneStatus(t *testing.T) {
	type args struct {
		c replyReader
	}
	tests := []struct {
		name    string
		args    args
		want    *ZoneStatus
		wantErr bool
	}{
		{
			name: "static slave zone refreshing",
			args: args{
				NewStaticReply(strings.Split(`zone:	example.org
	state: refreshing
	served-serial: none
	commit-serial: none
	wait: "99 sec between attempts"`, "\n")),
			},
			want: &ZoneStatus{
				Zone:  "example.org",
				State: "refreshing",
				Attributes: map[string]string{
					"served-serial": "none",
					"commit-serial": "none",
					"wait":          "99 sec between attempts",
				},
			},
			wantErr: false,
		},
		{
			name: "patterned primary zone",
			args: args{
				NewStaticReply(strings.Split(`zone:	example.dk.
	pattern: replica
	state: primary`, "\n")),
			},
			want: &ZoneStatus{
				Zone:  "example.dk.",
				State: "primary",
				Attributes: map[string]string{
					"pattern": "replica",
				},
			},
			wantErr: false,
		},
		{
			name: "status primary zone",
			args: args{
				NewStaticReply(strings.Split(`zone:	example.com
	state: primary`, "\n")),
			},
			want: &ZoneStatus{
				Zone:       "example.com",
				State:      "primary",
				Attributes: map[string]string{},
			},
			wantErr: false,
		},
		{
			name: "not configured zone",
			args: args{
				NewStaticReply([]string{}),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseZoneStatus(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseZoneStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseZoneStatus() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseAddCookieSecretReply(t *testing.T) {
	type args struct {
		c replyReader
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "wrong length",
			args: args{
				NewStaticReply(strings.Split(`invalid cookie secret: invalid argument length
please provide a 128bit hex encoded secret`, "\n")),
			},
			wantErr: true,
		},
		{
			name: "missing argument",
			args: args{
				NewStaticReply([]string{"error: missing argument (cookie_secret)"}),
			},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				NewStaticReply([]string{"ok"}),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := parseAddCookieSecretReply(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("parseAddCookieSecretReply() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
