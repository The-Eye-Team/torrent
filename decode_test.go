package torrentfile

import (
	"testing"
	"io/ioutil"
)

func TestUbuntuTorrentFile(t *testing.T) {
	var v Torrent

	dat, _ := ioutil.ReadFile("./ubuntu-17.10.1-desktop-amd64.iso.torrent")
	Unmarshal(dat, &v)

	if v.Announce != "http://torrent.ubuntu.com:6969/Announce" {
		println("Expected: \"http://torrent.ubuntu.com:6969/Announce\", Actual: ", v.Announce)
		t.Fail()
	}
	if v.Info.Length != 1502576640 {
		println("Expected: \"1502576640\", Actual: ", v.Info.Length)
		t.Fail()
	}
	if v.Info.PieceLength != 524288 {
		println("Expected: \"524288\", Actual: ", v.Info.PieceLength)
		t.Fail()
	}
	if v.Info.Name != "ubuntu-17.10.1-desktop-amd64.iso" {
		println("Expected: \"ubuntu-17.10.1-desktop-amd64.iso\", Actual: ", v.Info.Name)
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
