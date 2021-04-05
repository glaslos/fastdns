package fastdns

import (
	"encoding/hex"
	"net"
	"testing"
)

var mockHandlerMessage = &Message{
	nil,
	[]byte("hk.phus.lu"),
	Header{
		ID:      0x0002,
		QR:      0x00,
		Opcode:  0x0000,
		AA:      0x00,
		TC:      0x00,
		RD:      0x01,
		RA:      0x00,
		Z:       0x00,
		RCODE:   0x00,
		QDCount: 0x01,
		ANCount: 0x00,
		NSCount: 0x00,
		ARCount: 0x00,
	},
	Question{
		Name:  []byte("\x02hk\x04phus\x02lu\x00"),
		Type:  TypeA,
		Class: ClassINET,
	},
}

func TestHandlerError(t *testing.T) {
	var cases = []struct {
		Hex   string
		Rcode Rcode
	}{
		{
			"000281030000000000000000",
			RcodeNameError,
		},
	}

	rw := &MemoryResponseWriter{}
	if rw.RemoteAddr() != nil {
		t.Errorf("MemoryResponseWriter shall return empty addr")
	}
	if rw.LocalAddr() != nil {
		t.Errorf("MemoryResponseWriter shall return empty addr")
	}
	for _, c := range cases {
		Error(rw, mockHandlerMessage, c.Rcode)
		if got, want := hex.EncodeToString(rw.Data), c.Hex; got != want {
			t.Errorf("Error(%v) error got=%#v want=%#v", c.Rcode, got, want)
		}
	}
}

func TestHandlerHost(t *testing.T) {
	var cases = []struct {
		Hex string
		IP  net.IP
		TTL uint32
	}{
		{
			"00028100000100010000000002686b0470687573026c750000010001c00c000100010000012c000401020408",
			net.IP{1, 2, 4, 8},
			300,
		},
	}

	rw := &MemoryResponseWriter{}
	for _, c := range cases {
		HOST(rw, mockHandlerMessage, c.TTL, []net.IP{c.IP})
		if got, want := hex.EncodeToString(rw.Data), c.Hex; got != want {
			t.Errorf("HOST(%v) error got=%#v want=%#v", c.IP, got, want)
		}
	}
}

func TestHandlerCNAME(t *testing.T) {
	var cases = []struct {
		Hex   string
		CNAME string
		TTL   uint32
	}{
		{
			"00028100000100010000000002686b0470687573026c750000010001c00c000500010000012c001305636e616d65076578616d706c6503636f6d00",
			"cname.example.com",
			300,
		},
	}

	rw := &MemoryResponseWriter{}
	for _, c := range cases {
		CNAME(rw, mockHandlerMessage, c.TTL, []string{c.CNAME}, nil)
		if got, want := hex.EncodeToString(rw.Data), c.Hex; got != want {
			t.Errorf("CNAME(%v) error got=%#v want=%#v", c.CNAME, got, want)
		}
	}
}

func TestHandlerNS(t *testing.T) {
	var cases = []struct {
		Hex        string
		Nameserver string
		TTL        uint32
	}{
		{
			"00028100000100010000000002686b0470687573026c750000010001c00c000200010000012c0010026e73076578616d706c6503636f6d00",
			"ns.example.com",
			300,
		},
	}

	rw := &MemoryResponseWriter{}
	for _, c := range cases {
		NS(rw, mockHandlerMessage, c.TTL, []string{c.Nameserver})
		if got, want := hex.EncodeToString(rw.Data), c.Hex; got != want {
			t.Errorf("NS(%v) error got=%#v want=%#v", c.Nameserver, got, want)
		}
	}
}

func TestHandlerSOA(t *testing.T) {
	var cases = []struct {
		Hex     string
		TTL     uint32
		MName   string
		RName   string
		Serial  uint32
		Refresh uint32
		Retry   uint32
		Expire  uint32
		Minimum uint32
	}{
		{
			"00028100000100010000000002686b0470687573026c750000010001c00c000600010000012c003a036e733106676f6f676c6503636f6d0009646e732d61646d696e06676f6f676c6503636f6d00400000000000038400000384000007080000003c",
			300,
			"ns1.google.com",
			"dns-admin.google.com",
			1073741824,
			900,
			900,
			1800,
			60,
		},
	}

	rw := &MemoryResponseWriter{}
	for _, c := range cases {
		SOA(rw, mockHandlerMessage, c.TTL, c.MName, c.RName, c.Serial, c.Refresh, c.Retry, c.Expire, c.Minimum)
		if got, want := hex.EncodeToString(rw.Data), c.Hex; got != want {
			t.Errorf("SOA(%v) error got=%#v want=%#v", c.MName, got, want)
		}
	}
}

func TestHandlerSRV(t *testing.T) {
	var cases = []struct {
		Hex       string
		SRV       string
		Priovrity uint16
		Weight    uint16
		Port      uint16
		TTL       uint32
	}{
		{
			"00028100000100010000000002686b0470687573026c750000010001c00c002100010000012c001c03e803e81f41087365727669636531076578616d706c6503636f6d00",
			"service1.example.com",
			1000,
			1000,
			8001,
			300,
		},
	}

	rw := &MemoryResponseWriter{}
	for _, c := range cases {
		SRV(rw, mockHandlerMessage, c.TTL, c.SRV, c.Priovrity, c.Weight, c.Port)
		if got, want := hex.EncodeToString(rw.Data), c.Hex; got != want {
			t.Errorf("SRV(%v) error got=%#v want=%#v", c.SRV, got, want)
		}
	}
}

func TestHandlerMX(t *testing.T) {
	var cases = []struct {
		Hex string
		MX  string
		TTL uint32
	}{
		{
			"00028100000100010000000002686b0470687573026c750000010001c00c000f00010000012c0013000a03707472076578616d706c65036f726700",
			"ptr.example.org",
			300,
		},
	}

	rw := &MemoryResponseWriter{}
	for _, c := range cases {
		MX(rw, mockHandlerMessage, c.TTL, []MXRecord{{10, c.MX}})
		if got, want := hex.EncodeToString(rw.Data), c.Hex; got != want {
			t.Errorf("MX(%v) error got=%#v want=%#v", c.MX, got, want)
		}
	}
}

func TestHandlerPTR(t *testing.T) {
	var cases = []struct {
		Hex string
		PTR string
		TTL uint32
	}{
		{
			"00028100000100010000000002686b0470687573026c750000010001c00c000c00010000012c001103707472076578616d706c65036f726700",
			"ptr.example.org",
			300,
		},
	}

	rw := &MemoryResponseWriter{}
	for _, c := range cases {
		PTR(rw, mockHandlerMessage, c.TTL, c.PTR)
		if got, want := hex.EncodeToString(rw.Data), c.Hex; got != want {
			t.Errorf("PTR(%v) error got=%#v want=%#v", c.PTR, got, want)
		}
	}
}

func TestHandlerTXT(t *testing.T) {
	var cases = []struct {
		Hex string
		TXT string
		TTL uint32
	}{
		{
			"00028100000100010000000002686b0470687573026c750000010001c00c001000010000012c000e0d69616d617478747265636f7264",
			"iamatxtrecord",
			300,
		},
	}

	rw := &MemoryResponseWriter{}
	for _, c := range cases {
		TXT(rw, mockHandlerMessage, c.TTL, c.TXT)
		if got, want := hex.EncodeToString(rw.Data), c.Hex; got != want {
			t.Errorf("TXT(%v) error got=%#v want=%#v", c.TXT, got, want)
		}
	}
}

type nilResponseWriter struct{}

func (rw *nilResponseWriter) RemoteAddr() net.Addr { return nil }

func (rw *nilResponseWriter) LocalAddr() net.Addr { return nil }

func (rw *nilResponseWriter) Write(p []byte) (n int, err error) { return len(p), nil }

func BenchmarkHOST(b *testing.B) {
	ips := []net.IP{net.ParseIP("8.8.8.8")}
	for i := 0; i < b.N; i++ {
		HOST(&nilResponseWriter{}, mockHandlerMessage, 3000, ips)
	}
}

func BenchmarkCNAME(b *testing.B) {
	cnames := []string{"cname.example.org"}
	for i := 0; i < b.N; i++ {
		CNAME(&nilResponseWriter{}, mockHandlerMessage, 3000, cnames, nil)
	}
}

func BenchmarkNS(b *testing.B) {
	nameservers := []string{"ns1.google.com", "ns2.google.com"}
	for i := 0; i < b.N; i++ {
		NS(&nilResponseWriter{}, mockHandlerMessage, 3000, nameservers)
	}
}

func BenchmarkSOA(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SOA(&nilResponseWriter{}, mockHandlerMessage, 3000, "ns1.google.com", "dns-admin.google.com", 42, 900, 900, 1800, 60)
	}
}

func BenchmarkSRV(b *testing.B) {
	srv := "service1.example.org"
	for i := 0; i < b.N; i++ {
		SRV(&nilResponseWriter{}, mockHandlerMessage, 3000, srv, 100, 100, 443)
	}
}

func BenchmarkPTR(b *testing.B) {
	ptr := "ptr.example.org"
	for i := 0; i < b.N; i++ {
		PTR(&nilResponseWriter{}, mockHandlerMessage, 3000, ptr)
	}
}

func BenchmarkMX(b *testing.B) {
	mx := []MXRecord{{100, "mail.google.com"}}
	for i := 0; i < b.N; i++ {
		MX(&nilResponseWriter{}, mockHandlerMessage, 3000, mx)
	}
}

func BenchmarkTXT(b *testing.B) {
	txt := "iamatxtrecord"
	for i := 0; i < b.N; i++ {
		TXT(&nilResponseWriter{}, mockHandlerMessage, 3000, txt)
	}
}
