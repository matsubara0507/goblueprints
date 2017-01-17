package trace

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	tracer.Trace("こんにちは、trace パッケージ")
	if buf.String() != "こんにちは、trace パッケージ\n" {
		t.Errorf("'%s'という誤った文字列が出力されました", buf.String())
	}

	// case trace Off
	tracer = New(nil)
	tracer.Trace("これは表示されない")
}
