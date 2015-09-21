package ev

import (
	"testing"
)

func TestPipe(t *testing.T) {
	p, err := NewPipe()
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()

	_, err = p.Write([]byte("hello, world!"))
	if err != nil {
		t.Fatal(err)
	}
	out := make([]byte, 13)
	_, err = p.Read(out)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != "hello, world!" {
		t.Errorf("expected %q; got %q", "hello, world!", out)
	}
}