package utils

import (
	"testing"
)

func TestValidateDomain(t *testing.T) {
	tests := []struct {
		name    string
		domain  interface{}
		wantErr bool
	}{
		{"valid domain", "example.com", false},
		{"valid subdomain", "www.example.com", false},
		{"valid multi-level subdomain", "blog.www.example.com", false},
		{"valid domain with hyphen", "my-site.example.com", false},
		{"valid domain with numbers", "site123.example.com", false},
		{"valid short tld", "example.io", false},
		{"valid long tld", "example.photography", false},
		{"invalid - no tld", "example", true},
		{"invalid - starts with hyphen", "-example.com", true},
		{"invalid - ends with hyphen", "example-.com", true},
		{"invalid - double dots", "example..com", true},
		{"invalid - spaces", "example .com", true},
		{"invalid - single char tld", "example.c", true},
		{"invalid type", 123, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDomain(tt.domain)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDomain() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateSiteID(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{"valid short name", "abc", false},
		{"valid medium name", "mysite123", false},
		{"valid 16 char name", "abcdefghijklmnop", false},
		{"valid all numbers", "123456", false},
		{"invalid - too short", "ab", true},
		{"invalid - too long", "abcdefghijklmnopq", true},
		{"invalid - has hyphen", "my-site", true},
		{"invalid - has underscore", "my_site", true},
		{"invalid - has space", "my site", true},
		{"invalid - has special char", "mysite!", true},
		{"invalid type", 123, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSiteID(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSiteID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   interface{}
		wantErr bool
	}{
		{"valid email", "user@example.com", false},
		{"valid with plus", "user+tag@example.com", false},
		{"valid with dots", "user.name@example.com", false},
		{"valid with numbers", "user123@example.com", false},
		{"valid subdomain", "user@mail.example.com", false},
		{"invalid - no @", "userexample.com", true},
		{"invalid - no domain", "user@", true},
		{"invalid - no user", "@example.com", true},
		{"invalid - spaces", "user @example.com", true},
		{"invalid - double @", "user@@example.com", true},
		{"invalid type", 123, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePasswordStrength(t *testing.T) {
	tests := []struct {
		name     string
		password interface{}
		wantErr  bool
	}{
		{"valid password", "MyP@ssword123", false},
		{"valid with special chars", "Abc123!@#$%^", false},
		{"valid longer password", "MyVerySecureP@ss123", false},
		{"invalid - too short", "MyP@ss1", true},
		{"invalid - no uppercase", "myp@ssword123", true},
		{"invalid - no lowercase", "MYP@SSWORD123", true},
		{"invalid - no number", "MyP@sswordabc", true},
		{"invalid - no special", "MyPassword123", true},
		{"invalid - only lowercase", "mypassword!!", true},
		{"invalid type", 123, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePasswordStrength(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePasswordStrength() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateIP(t *testing.T) {
	tests := []struct {
		name    string
		ip      interface{}
		wantErr bool
	}{
		{"valid IPv4", "192.168.1.1", false},
		{"valid IPv4 zeros", "0.0.0.0", false},
		{"valid IPv4 max", "255.255.255.255", false},
		{"valid IPv6", "2001:0db8:85a3:0000:0000:8a2e:0370:7334", false},
		{"valid IPv6 short", "::1", false},
		{"invalid - too many octets", "192.168.1.1.1", true},
		{"invalid - too few octets", "192.168.1", true},
		{"invalid - out of range", "256.1.1.1", true},
		{"invalid - letters", "abc.def.ghi.jkl", true},
		{"invalid - domain name", "example.com", true},
		{"invalid type", 123, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateIP(tt.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateIP() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePort(t *testing.T) {
	tests := []struct {
		name    string
		port    interface{}
		wantErr bool
	}{
		{"valid port 22", 22, false},
		{"valid port 80", 80, false},
		{"valid port 443", 443, false},
		{"valid port 1", 1, false},
		{"valid port max", 65535, false},
		{"invalid - zero", 0, true},
		{"invalid - negative", -1, true},
		{"invalid - too large", 65536, true},
		{"invalid type", "22", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePort(tt.port)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePort() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
