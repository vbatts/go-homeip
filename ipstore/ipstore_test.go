package ipstore

import (
	"os"
	"testing"
)

func TestInitClose(t *testing.T) {
	filename := "/tmp/foo.sqlite"
	err := InitFilename(filename)
	if err != nil {
		t.Fatal("failed to open", filename, "with err:", err)
	}

	err = Close()
	if err != nil {
		t.Fatal("failed to close", filename, "with err:", err)
	}

	os.Remove(filename)
}

func TestTransactions(t *testing.T) {
	filename := "/tmp/foo.sqlite"
	err := InitFilename(filename)
	if err != nil {
		t.Fatal("failed to open", filename, "with err:", err)
	}
	defer Close()

	exists, err := HostExists("foobar")
	if err != nil {
    t.Fatal("Could not HostExists due to:", err)
	}
	if exists != false {
		t.Error("This initial check should be false")
	}

	err = SetHostIp("foobar", "0.0.0.0")
	if err != nil {
		t.Fatal("Could not SetHostIp due to:", err)
	}

	ip, err := GetHostIp("foobar")
	if err != nil {
		t.Fatal("Could not GetHostIp due to:", err)
	}
	if ip != "0.0.0.0" {
		t.Error("the ip addresses should be the same")
	}

	os.Remove(filename)
}
