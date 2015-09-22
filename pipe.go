package ev

func NewPipe() (*PipeReader, *PipeWriter, error) {
	rd, wd, err := pipe()
	if err != 0 {
		return nil, nil, err
	}
	return &PipeReader{rfd(rd)}, &PipeWriter{wfd(wd)}, nil
}

type PipeReader struct {
	r *evfd
}

type PipeWriter struct {
	w *evfd
}

func (p *PipeWriter) Write(b []byte) (int, error) {
	return p.w.Write(b)
}

func (p *PipeReader) Read(b []byte) (int, error) {
	return p.r.Read(b)
}

func (p *PipeReader) Close() error {
	return p.r.Close()
}

func (p *PipeWriter) Close() error {
	return p.w.Close()
}
