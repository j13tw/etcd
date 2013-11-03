package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/coreos/etcd/tests"
	"github.com/stretchr/testify/assert"
)

// Ensure that we can start a v2 node from the log of a v1 node.
func TestV1Migration(t *testing.T) {
	path, _ := ioutil.TempDir("", "etcd-")
	os.RemoveAll(path)
	defer os.RemoveAll(path)

	// Copy over fixture files.
	if err := exec.Command("cp", "-r", "../fixtures/v1/complete", path).Run(); err != nil {
		panic("Fixture initialization error")
	}

	procAttr := new(os.ProcAttr)
	procAttr.Files = []*os.File{nil, os.Stdout, os.Stderr}
	args := []string{"etcd", fmt.Sprintf("-d=%s", path)}

	process, err := os.StartProcess(EtcdBinPath, args, procAttr)
	if err != nil {
		t.Fatal("start process failed:" + err.Error())
		return
	}
	defer process.Kill()
	time.Sleep(time.Second)


	// Ensure deleted message is removed.
	resp, err := tests.Get("http://localhost:4001/v2/keys/message")
	tests.ReadBody(resp)
	assert.Nil(t, err, "")
	assert.Equal(t, resp.StatusCode, 404, "")

	// Ensure TTL'd message is removed.
	resp, err = tests.Get("http://localhost:4001/v2/keys/foo")
	tests.ReadBody(resp)
	assert.Nil(t, err, "")
	assert.Equal(t, resp.StatusCode, 404, "")
}

