package ev

func NewPipe() (*Pipe, error) {
	rd, wd, err := pipe()
	if err != 0 {
		return nil, err
	}
	return &Pipe{
		w: wfd(wd),
		r: rfd(rd),
	}, nil
}

type Pipe struct {
	w *evfd
	r *evfd
}

func (p *Pipe) Write(b []byte) (int, error) {
	return p.w.Write(b)
}

func (p *Pipe) Read(b []byte) (int, error) {
	return p.r.Read(b)
}

func (p *Pipe) Close() error {
	we, re := p.w.Close(), p.r.Close()
	if we == nil {
		we = re
	}
	return we
}