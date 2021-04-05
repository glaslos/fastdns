package fastdns

import (
	"bytes"
	"encoding/hex"
	"reflect"
	"testing"
)

func TestParseMessageOK(t *testing.T) {
	var cases = []struct {
		Message *Message
	}{
		{
			/*
				Domain Name System (query)
				    Transaction ID: 0x0001
				    Flags: 0x0100 Standard query
				        0... .... .... .... = Response: Message is a query
				        .000 0... .... .... = Opcode: Standard query (0)
				        .... ..0. .... .... = Truncated: Message is not truncated
				        .... ...1 .... .... = Recursion desired: Do query recursively
				        .... .... .0.. .... = Z: reserved (0)
				        .... .... ...0 .... = Non-authenticated data: Unacceptable
				    Questions: 1
				    Answer RRs: 0
				    Authority RRs: 0
				    Additional RRs: 0
				    Queries
				        1.50.168.192.in-addr.arpa: type PTR, class IN
				            Name: 1.50.168.192.in-addr.arpa
				            [Name Length: 25]
				            [Label Count: 6]
				            Type: PTR (domain name PoinTeR) (12)
				            Class: IN (0x0001)
			*/
			&Message{
				[]byte("\x00\x01\x01\x00\x00\x01\x00\x00\x00\x00\x00\x00\x01\x31\x02\x35\x30\x03\x31\x36\x38\x03\x31\x39\x32\x07\x69\x6e\x2d\x61\x64\x64\x72\x04\x61\x72\x70\x61\x00\x00\x0c\x00\x01"),
				[]byte("1.50.168.192.in-addr.arpa"),
				Header{
					ID:      0x0001,
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
					Name:  []byte("\x011\x0250\x03168\x03192\x07in-addr\x04arpa\x00"),
					Type:  TypePTR,
					Class: ClassINET,
				},
			},
		},
		{
			/*
				Domain Name System (query)
				    Transaction ID: 0x0002
				    Flags: 0x0100 Standard query
				        0... .... .... .... = Response: Message is a query
				        .000 0... .... .... = Opcode: Standard query (0)
				        .... ..0. .... .... = Truncated: Message is not truncated
				        .... ...1 .... .... = Recursion desired: Do query recursively
				        .... .... .0.. .... = Z: reserved (0)
				        .... .... ...0 .... = Non-authenticated data: Unacceptable
				    Questions: 1
				    Answer RRs: 0
				    Authority RRs: 0
				    Additional RRs: 0
				    Queries
				        hk.phus.lu: type A, class IN
				            Name: hk.phus.lu
				            [Name Length: 10]
				            [Label Count: 3]
				            Type: A (Host Address) (1)
				            Class: IN (0x0001)
			*/
			&Message{
				[]byte("\x00\x02\x01\x00\x00\x01\x00\x00\x00\x00\x00\x00\x02\x68\x6b\x04\x70\x68\x75\x73\x02\x6c\x75\x00\x00\x01\x00\x01"),
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
			},
		},
	}

	for _, c := range cases {
		msg := AcquireMessage()
		err := ParseMessage(msg, c.Message.Raw, true)
		if err != nil {
			t.Errorf("ParseMessage(%x) error: %+v", c.Message.Raw, err)
		}
		if got, want := msg, c.Message; !reflect.DeepEqual(got, want) {
			t.Errorf("ParseMessage(%x) error got=%#v want=%#v", c.Message.Raw, got, want)
		}
		ReleaseMessage(msg)
	}
}

func TestParseMessageError(t *testing.T) {
	var cases = []struct {
		Hex   string
		Error error
	}{
		{
			"0001010000010000000000",
			ErrInvalidHeader,
		},
		{
			"00020100000000000000000002686b0470687573026c7500000100",
			ErrInvalidHeader,
		},
		{
			"00020100000100000000000002686b0470687573026c7500000100",
			ErrInvalidQuestion,
		},
	}

	for _, c := range cases {
		payload, err := hex.DecodeString(c.Hex)
		if err != nil {
			t.Errorf("hex.DecodeString(%v) error: %+v", c.Hex, err)
		}
		var msg Message
		err = ParseMessage(&msg, payload, true)
		if err != c.Error {
			t.Errorf("ParseMessage(%x) should error: %+v", payload, c.Error)
		}
	}
}

func TestAppendMessage(t *testing.T) {
	var cases = []struct {
		Message *Message
	}{
		{
			/*
				Domain Name System (query)
				    Transaction ID: 0x0001
				    Flags: 0x0100 Standard query
				        0... .... .... .... = Response: Message is a query
				        .000 0... .... .... = Opcode: Standard query (0)
				        .... ..0. .... .... = Truncated: Message is not truncated
				        .... ...1 .... .... = Recursion desired: Do query recursively
				        .... .... .0.. .... = Z: reserved (0)
				        .... .... ...0 .... = Non-authenticated data: Unacceptable
				    Questions: 1
				    Answer RRs: 0
				    Authority RRs: 0
				    Additional RRs: 0
				    Queries
				        1.50.168.192.in-addr.arpa: type PTR, class IN
				            Name: 1.50.168.192.in-addr.arpa
				            [Name Length: 25]
				            [Label Count: 6]
				            Type: PTR (domain name PoinTeR) (12)
				            Class: IN (0x0001)
			*/
			&Message{
				[]byte("\x00\x01\x01\x00\x00\x01\x00\x00\x00\x00\x00\x00\x01\x31\x02\x35\x30\x03\x31\x36\x38\x03\x31\x39\x32\x07\x69\x6e\x2d\x61\x64\x64\x72\x04\x61\x72\x70\x61\x00\x00\x0c\x00\x01"),
				[]byte("1.50.168.192.in-addr.arpa"),
				Header{
					ID:      0x0001,
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
					Name:  []byte("\x011\x0250\x03168\x03192\x07in-addr\x04arpa\x00"),
					Type:  TypePTR,
					Class: ClassINET,
				},
			},
		},
		{
			/*
				Domain Name System (query)
				    Transaction ID: 0x0002
				    Flags: 0x0100 Standard query
				        0... .... .... .... = Response: Message is a query
				        .000 0... .... .... = Opcode: Standard query (0)
				        .... ..0. .... .... = Truncated: Message is not truncated
				        .... ...1 .... .... = Recursion desired: Do query recursively
				        .... .... .0.. .... = Z: reserved (0)
				        .... .... ...0 .... = Non-authenticated data: Unacceptable
				    Questions: 1
				    Answer RRs: 0
				    Authority RRs: 0
				    Additional RRs: 0
				    Queries
				        hk.phus.lu: type A, class IN
				            Name: hk.phus.lu
				            [Name Length: 10]
				            [Label Count: 3]
				            Type: A (Host Address) (1)
				            Class: IN (0x0001)
			*/
			&Message{
				[]byte("\x00\x02\x01\x00\x00\x01\x00\x00\x00\x00\x00\x00\x02\x68\x6b\x04\x70\x68\x75\x73\x02\x6c\x75\x00\x00\x01\x00\x01"),
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
			},
		},
	}

	for _, c := range cases {
		if got, want := AppendMessage(nil, c.Message), c.Message.Raw; !bytes.Equal(got, want) {
			t.Errorf("AppendMessage(%v) error got=%#v want=%#v", c.Message, got, want)
		}
	}
}

func BenchmarkParseMessage(b *testing.B) {
	payload, _ := hex.DecodeString("00020100000100000000000002686b0470687573026c750000010001")
	var msg Message

	for i := 0; i < b.N; i++ {
		if err := ParseMessage(&msg, payload, false); err != nil {
			b.Errorf("ParseMessage(%+v) error: %+v", payload, err)
		}
	}
}