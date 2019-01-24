package go2node

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestExecNode_Reader(t *testing.T) {
	cmd := exec.Command("node", "node_test.js", "reader")
	channel, err := ExecNode(cmd)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		cmd.Process.Kill()
	}()

	msg := <-channel.Reader
	const expectedContent = `{"black":"heart"}`
	if strings.Compare(string(msg.Message), expectedContent) != 0 {
		t.Fatal("Message not matched: ", string(msg.Message))
	}
}

func TestExecNode_Writer(t *testing.T) {
	cmd := exec.Command("node", "node_test.js", "writer")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	channel, err := ExecNode(cmd)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		cmd.Process.Kill()
	}()

	sp, _ := Socketpair()
	msg := &NodeMessage{
		Message: `65535`,
		Handle:  sp[0],
	}
	channel.Writer <- msg

	msg = <-channel.Reader
	if string(msg.Message) != `{"value":"6553588"}` {
		t.Fatal("Message not matched: ", msg.Message)
	}
	if msg.Handle.Fd() == 0 {
		t.Fatal("Handle is empty")
	}
}
