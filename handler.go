package fastdns

import (
	"net"
)

// Handler is implemented by any value that implements ServeDNS.
type Handler interface {
	ServeDNS(rw ResponseWriter, req *Message)
}

// Error replies to the request with the specified Rcode.
func Error(rw ResponseWriter, req *Message, code Rcode) {
	req.Raw = AppendHeaderQuestion(req.Raw[:0], req, code, 0, 0, 0, 0)
	_, _ = rw.Write(req.Raw)
}

// HOST replies to the request with the specified Host records.
func HOST(rw ResponseWriter, req *Message, ttl uint32, ips []net.IP) {
	req.Raw = AppendHeaderQuestion(req.Raw[:0], req, RcodeSuccess, 1, uint16(len(ips)), 0, 0)
	req.Raw = AppendHostRecord(req.Raw, req, ttl, ips)
	_, _ = rw.Write(req.Raw)
}

// CNAME replies to the request with the specified CName and Host records.
func CNAME(rw ResponseWriter, req *Message, ttl uint32, cnames []string, ips []net.IP) {
	req.Raw = AppendHeaderQuestion(req.Raw[:0], req, RcodeSuccess, 1, uint16(len(cnames)+len(ips)), 0, 0)
	req.Raw = AppendCNAMERecord(req.Raw, req, ttl, cnames, ips)
	_, _ = rw.Write(req.Raw)
}

// NS replies to the request with the specified CName and Host records.
func NS(rw ResponseWriter, req *Message, ttl uint32, nameservers []string) {
	req.Raw = AppendHeaderQuestion(req.Raw[:0], req, RcodeSuccess, 1, uint16(len(nameservers)), 0, 0)
	req.Raw = AppendNSRecord(req.Raw, req, ttl, nameservers)
	_, _ = rw.Write(req.Raw)
}

// SOA replies to the request with the specified SOA records.
func SOA(rw ResponseWriter, req *Message, ttl uint32, mname, rname string, serial, refresh, retry, expire, minimum uint32) {
	req.Raw = AppendHeaderQuestion(req.Raw[:0], req, RcodeSuccess, 1, 1, 0, 0)
	req.Raw = AppendSOARecord(req.Raw, req, ttl, mname, rname, serial, refresh, retry, expire, minimum)
	_, _ = rw.Write(req.Raw)
}

// SRV replies to the request with the specified SRV records.
func SRV(rw ResponseWriter, req *Message, ttl uint32, srv string, priovrity, weight, port uint16) {
	req.Raw = AppendHeaderQuestion(req.Raw[:0], req, RcodeSuccess, 1, 1, 0, 0)
	req.Raw = AppendSRVRecord(req.Raw, req, ttl, srv, priovrity, weight, port)
	_, _ = rw.Write(req.Raw)
}

// MX replies to the request with the specified MX records.
func MX(rw ResponseWriter, req *Message, ttl uint32, mx []MXRecord) {
	req.Raw = AppendHeaderQuestion(req.Raw[:0], req, RcodeSuccess, 1, uint16(len(mx)), 0, 0)
	req.Raw = AppendMXRecord(req.Raw, req, ttl, mx)
	_, _ = rw.Write(req.Raw)
}

// PTR replies to the request with the specified PTR records.
func PTR(rw ResponseWriter, req *Message, ttl uint32, ptr string) {
	req.Raw = AppendHeaderQuestion(req.Raw[:0], req, RcodeSuccess, 1, 1, 0, 0)
	req.Raw = AppendPTRRecord(req.Raw, req, ttl, ptr)
	_, _ = rw.Write(req.Raw)
}

// TXT replies to the request with the specified TXT records.
func TXT(rw ResponseWriter, req *Message, ttl uint32, txt string) {
	req.Raw = AppendHeaderQuestion(req.Raw[:0], req, RcodeSuccess, 1, 1, 0, 0)
	req.Raw = AppendTXTRecord(req.Raw, req, ttl, txt)
	_, _ = rw.Write(req.Raw)
}
