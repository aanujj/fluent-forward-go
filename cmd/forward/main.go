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

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aanujj/fluent-forward-go/fluent/client"
	"github.com/aanujj/fluent-forward-go/fluent/protocol"
)

var (
	tagVar string
)

func init() {
	flag.StringVar(&tagVar, "tag", "test.message", "-tag <dot-delimited tag>")
	flag.StringVar(&tagVar, "t", "test.message", "-t <dot-delimited tag> (shorthand for -tag)")
}

func main() {
	flag.Parse()

	c := client.New(client.ConnectionOptions{
		RequireAck: true,
	})

	err := c.Connect()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to connect, exiting", err)
		os.Exit(-1)
	}

	record := map[string]interface{}{
		"first": "Sir",
		"last":  "Gawain",
		"enemy": "Green Knight",
		"equipment": []string{
			"sword",
			"lance",
			"full plate",
		},
	}

	entries := []protocol.EntryExt{
		{
			Timestamp: protocol.EventTimeNow(),
			Record: map[string]interface{}{
				"first": "Edgar",
				"last":  "Winter",
				"enemy": "wimpy music",
			},
		},
		{
			Timestamp: protocol.EventTimeNow(),
			Record: map[string]interface{}{
				"first": "George",
				"last":  "Clinton",
				"enemy": "Sir Nose D Voidoffunk",
			},
		},
	}

	err = c.SendMessage(tagVar, record)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = c.SendMessageExt(tagVar, record)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = c.SendForward(tagVar, entries)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = c.SendPacked(tagVar+".packed", entries)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = c.SendCompressed(tagVar+".compressed", entries)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Messages sent")

	os.Exit(0)
}
