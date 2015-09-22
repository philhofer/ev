package ev

import (
	"bytes"
	"crypto/rand"
	"io"
	"os"
	"testing"
)

func copy(dst *PipeWriter, src io.Reader, t *testing.T) {
	_, err := io.Copy(dst, src)
	if err != nil {
		t.Fatal(err)
	}
	if err := dst.Close(); err != nil {
		t.Errorf("(*PipeWriter).Close(): %s", err)
	}
}

func TestPipe(t *testing.T) {
	r, w, err := NewPipe()
	if err != nil {
		t.Fatal(err)
	}

	var outbuf bytes.Buffer
	var inbuf bytes.Buffer

	tr := io.TeeReader(&io.LimitedReader{
		R: rand.Reader,
		N: 1<<16 + 500,
	}, &inbuf)

	go copy(w, tr, t)

	_, err = io.Copy(&outbuf, r)
	if err != nil {
		t.Fatal(err)
	}
	err = r.Close()
	if err != nil {
		t.Errorf("(*PipeReader).Close(): %s", err)
	}

	if !bytes.Equal(outbuf.Bytes(), inbuf.Bytes()) {
		t.Error("out != in")
	}
}

type nowhere struct{}

func (n *nowhere) Write(p []byte) (int, error) { return len(p), nil }

func benchmarkPipe(b *testing.B, size int64) {
	data := make([]byte, size)
	out := make([]byte, size)
	rand.Read(data)
	r, w, err := NewPipe()
	if err != nil {
		b.Fatal(err)
	}

	b.SetBytes(size)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		go w.Write(data)
		_, err = io.ReadFull(r, out)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPipe4MB(b *testing.B) {
	benchmarkPipe(b, 1<<22)
}

func BenchmarkOSPipe4MB(b *testing.B) {
	data := make([]byte, 1<<22)
	out := make([]byte, 1<<22)
	rand.Read(data)
	r, w, err := os.Pipe()
	if err != nil {
		b.Fatal(err)
	}

	b.SetBytes(1 << 22)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		go w.Write(data)
		_, err = io.ReadFull(r, out)
		if err != nil {
			b.Fatal(err)
		}
	}
}
