package filtering

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/golang/snappy"
)

func Snappy() Filtering {
	return New(
		FilterFunc(func(bs []byte) ([]byte, error) {
			buf := &bytes.Buffer{}
			w := snappy.NewBufferedWriter(buf)
			if _, err := io.Copy(w, bytes.NewReader(bs)); err != nil {
				return nil, err
			}
			if err := w.Close(); err != nil {
				return nil, err
			}
			return buf.Bytes(), nil
		}),
		UnfilterFunc(func(bs []byte) ([]byte, error) {
			return ioutil.ReadAll(snappy.NewReader(bytes.NewReader(bs)))
		}))
}
