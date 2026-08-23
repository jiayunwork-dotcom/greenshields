package api

// qkRecorder wraps the single-point q(k) response so handleQK can
// encode and then release the handle. Close is idempotent so a
// deferred Close after an explicit Release cannot panic.
type qkRecorder struct {
	closed bool
}

func newQKRecorder() *qkRecorder {
	return &qkRecorder{}
}

func (r *qkRecorder) Close() {
	if r.closed {
		return
	}
	r.closed = true
}

func (r *qkRecorder) Release() {
	r.Close()
}
