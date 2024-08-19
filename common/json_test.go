package common;

import (
	"testing";
	"fmt";
);

func TestPeerMarshal(t *testing.T) {
	p := Peer {
		IP: "192.168.1.1",
		Port: 1234,
	}
	res, _ := p.MarshalJSON()
	fmt.Println(string(res[:]))
}

func TestFDDMarshal(t *testing.T) {
	fdd := FileDownloadData {
		FileID: "file id haha",
		FileName: "file name",
		Peers: []Peer {
			{"192.168.1.1", 1234},
			{"122.134.145.167", 789},
		},
	}
	res, _ := fdd.MarshalJSON()
	fmt.Println(string(res[:]))
	fdd.UnmarshalJSON(res)
	res, _ = fdd.MarshalJSON()
	fmt.Println(string(res[:]))
}
