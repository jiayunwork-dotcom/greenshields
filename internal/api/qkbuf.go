package api

// qkRecorder wraps the single-point q(k) response so handleQK can
// encode and then release the handle. Close is only safe once; a
// second Close panics so callers do not silently double-close.
type qkRecorder struct {
	closed bool
}

func newQKRecorder() *qkRecorder {
	return &qkRecorder{}
}

func (r *qkRecorder) Close() {
	if r.closed {
		panic("close of closed qk recorder")
	}
	r.closed = true
}

func (r *qkRecorder) Release() {
	r.Close()
}
