package common

import (
	"encoding/json"
	"errors"
)

type Peer struct {
	IP string;
	Port int;
};

func (self *Peer) FromDict(data map[string]interface{}) error {
	var ok bool;
	self.IP, ok = data["ip"].(string)
	if (!ok) {
		return errors.New("No IP")
	}
	self.Port, ok = data["port"].(int)
	if (!ok) {
		return errors.New("No port")
	}
	return nil
}

func (self Peer) MarshalJSON() ([]byte, error) {
	dict := map[string]interface{} {
		"ip": self.IP,
		"port": self.Port,
	}
	return json.Marshal(dict)
}

type FileDownloadData struct {
	FileID string;
	FileName string;
	Peers []Peer;
};

func (self FileDownloadData) MarshalJSON() ([]byte, error) {
	dict := map[string]interface{} {
		"file_id": self.FileID,
		"file_name": self.FileName,
		"peers": self.Peers,
	}
	return json.Marshal(dict)
}

func (self *FileDownloadData) UnmarshalJSON(data []byte) error {
	dict := make(map[string]interface{})
	err := json.Unmarshal(data, &dict)
	if (err != nil) {
		return err
	}
	var ok bool
	self.FileID, ok = dict["file_id"].(string)
	if (!ok) {
		return errors.New("No file_id")
	}
	self.FileName, ok = dict["file_name"].(string)
	if (!ok) {
		return errors.New("No file_name")
	}
	rawPeerList, ok := dict["peers"].([]map[string]interface{})
	if (!ok) {
		return errors.New("No peerlist")
	}
	self.Peers = make([]Peer, len(rawPeerList))
	for i := 0; i < len(rawPeerList); i++ {
		err = self.Peers[i].FromDict(rawPeerList[i])
	}
	return nil
}

const TRACKER_TCP_PORT int = 3141;
const TRACKER_UDP_PORT int = 1618;
const PEER_QUIC_PORT int = 2718;
