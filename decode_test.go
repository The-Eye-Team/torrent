package torrentfile

import (
	"testing"
	"io/ioutil"
)

func TestUbuntuTorrentFile(t *testing.T) {
	var v Torrent

	dat, _ := ioutil.ReadFile("./ubuntu-17.10.1-desktop-amd64.iso.torrent")
	Unmarshal(dat, &v)

	if v.announce != "http://torrent.ubuntu.com:6969/announce" {
		println("Expected: \"http://torrent.ubuntu.com:6969/announce\", Actual: ", v.announce)
		t.Fail()
	}
	if v.info.length != 1502576640 {
		println("Expected: \"1502576640\", Actual: ", v.info.length)
		t.Fail()
	}
	if v.info.pieceLength != 524288 {
		println("Expected: \"524288\", Actual: ", v.info.pieceLength)
		t.Fail()
	}
	if v.info.name != "ubuntu-17.10.1-desktop-amd64.iso" {
		println("Expected: \"ubuntu-17.10.1-desktop-amd64.iso\", Actual: ", v.info.name)
		t.Fail()
	}
}

// from fib_test.go
func BenchmarkTestUbuntuTorrentFile(b *testing.B) {
	dat, _ := ioutil.ReadFile("./ubuntu-17.10.1-desktop-amd64.iso.torrent")

	for n := 0; n < b.N; n++ {
		var v Torrent
		Unmarshal(dat, &v)
	}
}
