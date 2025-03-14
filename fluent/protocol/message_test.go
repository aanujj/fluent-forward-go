/*
MIT License

Copyright contributors to the fluent-forward-go project

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package protocol_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tinylib/msgp/msgp"

	. "github.com/aanujj/fluent-forward-go/fluent/protocol"
)

func TestMarshalNewMessage(t *testing.T) {
	record := map[string]interface{}{
		"first": "Sir",
		"last":  "Gawain",
		"equipment": []string{
			"sword",
		},
	}
	msg := NewMessage("tag", record)
	assert.Equal(t, msg.Tag, "tag")
	assert.Equal(t, msg.Record, record)
	assert.Greater(t, msg.Timestamp, int64(0))

	msgext := NewMessageExt("tag", record)
	assert.Equal(t, msgext.Tag, "tag")
	assert.Equal(t, msgext.Record, record)
	assert.Greater(t, msgext.Timestamp.Time.UTC().Nanosecond(), 0)
}

func TestMarshalUnmarshalMessage(t *testing.T) {
	v := Message{
		Options: &MessageOptions{},
	}
	bts, err := v.MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}
	left, err := v.UnmarshalMsg(bts)
	if err != nil {
		t.Fatal(err)
	}
	if len(left) > 0 {
		t.Errorf("%d bytes left over after UnmarshalMsg(): %q", len(left), left)
	}

	left, err = msgp.Skip(bts)
	if err != nil {
		t.Fatal(err)
	}
	if len(left) > 0 {
		t.Errorf("%d bytes left over after Skip(): %q", len(left), left)
	}
}

func BenchmarkMarshalMsgMessage(b *testing.B) {
	v := Message{
		Options: &MessageOptions{},
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := v.MarshalMsg(nil)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkAppendMsgMessage(b *testing.B) {
	v := Message{
		Options: &MessageOptions{},
	}
	bts := make([]byte, 0, v.Msgsize())
	bts, _ = v.MarshalMsg(bts[0:0])
	b.SetBytes(int64(len(bts)))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bts, _ = v.MarshalMsg(bts[0:0])
	}
}

func BenchmarkUnmarshalMessage(b *testing.B) {
	v := Message{
		Options: &MessageOptions{},
	}
	bts, _ := v.MarshalMsg(nil)
	b.ReportAllocs()
	b.SetBytes(int64(len(bts)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := v.UnmarshalMsg(bts)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestEncodeDecodeMessage(t *testing.T) {
	v := Message{
		Options: &MessageOptions{},
	}
	var buf bytes.Buffer
	err := msgp.Encode(&buf, &v)
	if err != nil {
		t.Error(err)
	}

	m := v.Msgsize()
	if buf.Len() > m {
		t.Log("WARNING: TestEncodeDecodeMessage Msgsize() is inaccurate")
	}

	vn := Message{}
	err = msgp.Decode(&buf, &vn)
	if err != nil {
		t.Error(err)
	}

	buf.Reset()
	err = msgp.Encode(&buf, &v)
	if err != nil {
		t.Error(err)
	}
	err = msgp.NewReader(&buf).Skip()
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkEncodeMessage(b *testing.B) {
	v := Message{}
	var buf bytes.Buffer
	err := msgp.Encode(&buf, &v)
	if err != nil {
		b.Error(err)
	}
	b.SetBytes(int64(buf.Len()))
	en := msgp.NewWriter(msgp.Nowhere)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := v.EncodeMsg(en)
		if err != nil {
			b.Error(err)
		}
	}
	en.Flush()
}

func BenchmarkDecodeMessage(b *testing.B) {
	v := Message{}
	var buf bytes.Buffer
	err := msgp.Encode(&buf, &v)
	if err != nil {
		b.Error(err)
	}
	b.SetBytes(int64(buf.Len()))
	rd := msgp.NewEndlessReader(buf.Bytes(), b)
	dc := msgp.NewReader(rd)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := v.DecodeMsg(dc)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestMarshalUnmarshalMessageExt(t *testing.T) {
	v := MessageExt{
		Options: &MessageOptions{},
	}
	bts, err := v.MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}
	left, err := v.UnmarshalMsg(bts)
	if err != nil {
		t.Fatal(err)
	}
	if len(left) > 0 {
		t.Errorf("%d bytes left over after UnmarshalMsg(): %q", len(left), left)
	}

	left, err = msgp.Skip(bts)
	if err != nil {
		t.Fatal(err)
	}
	if len(left) > 0 {
		t.Errorf("%d bytes left over after Skip(): %q", len(left), left)
	}
}

func BenchmarkMarshalMsgMessageExt(b *testing.B) {
	v := MessageExt{
		Options: &MessageOptions{},
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := v.MarshalMsg(nil)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkAppendMsgMessageExt(b *testing.B) {
	v := MessageExt{
		Options: &MessageOptions{},
	}
	bts := make([]byte, 0, v.Msgsize())
	bts, _ = v.MarshalMsg(bts[0:0])
	b.SetBytes(int64(len(bts)))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bts, _ = v.MarshalMsg(bts[0:0])
	}
}

func BenchmarkUnmarshalMessageExt(b *testing.B) {
	v := MessageExt{
		Options: &MessageOptions{},
	}
	bts, _ := v.MarshalMsg(nil)
	b.ReportAllocs()
	b.SetBytes(int64(len(bts)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := v.UnmarshalMsg(bts)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestEncodeDecodeMessageExt(t *testing.T) {
	v := MessageExt{
		Options: &MessageOptions{},
	}
	var buf bytes.Buffer
	err := msgp.Encode(&buf, &v)
	if err != nil {
		t.Error(err)
	}

	m := v.Msgsize()
	if buf.Len() > m {
		t.Log("WARNING: TestEncodeDecodeMessageExt Msgsize() is inaccurate")
	}

	vn := MessageExt{}
	err = msgp.Decode(&buf, &vn)
	if err != nil {
		t.Error(err)
	}

	buf.Reset()
	err = msgp.Encode(&buf, &v)
	if err != nil {
		t.Error(err)
	}
	err = msgp.NewReader(&buf).Skip()
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkEncodeMessageExt(b *testing.B) {
	v := MessageExt{
		Options: &MessageOptions{},
	}
	var buf bytes.Buffer
	err := msgp.Encode(&buf, &v)
	if err != nil {
		b.Error(err)
	}
	b.SetBytes(int64(buf.Len()))
	en := msgp.NewWriter(msgp.Nowhere)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := v.EncodeMsg(en)
		if err != nil {
			b.Error(err)
		}
	}
	en.Flush()
}

func BenchmarkDecodeMessageExt(b *testing.B) {
	v := MessageExt{
		Options: &MessageOptions{},
	}
	var buf bytes.Buffer
	err := msgp.Encode(&buf, &v)
	if err != nil {
		b.Error(err)
	}
	b.SetBytes(int64(buf.Len()))
	rd := msgp.NewEndlessReader(buf.Bytes(), b)
	dc := msgp.NewReader(rd)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := v.DecodeMsg(dc)
		if err != nil {
			b.Fatal(err)
		}
	}
}
