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

package client_test

import (
	"bytes"
	"crypto/tls"
	"errors"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"

	"time"

	"github.com/aanujj/fluent-forward-go/fluent/client"
	. "github.com/aanujj/fluent-forward-go/fluent/client"
	fclient "github.com/aanujj/fluent-forward-go/fluent/client"
	"github.com/aanujj/fluent-forward-go/fluent/client/clientfakes"
	"github.com/aanujj/fluent-forward-go/fluent/client/ws"
	"github.com/aanujj/fluent-forward-go/fluent/client/ws/ext"
	"github.com/aanujj/fluent-forward-go/fluent/client/ws/ext/extfakes"
	"github.com/aanujj/fluent-forward-go/fluent/client/ws/wsfakes"
	"github.com/aanujj/fluent-forward-go/fluent/protocol"
	"github.com/gorilla/websocket"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tinylib/msgp/msgp"
)

var _ = Describe("IAMAuthInfo", func() {
	It("gets and sets an IAM token", func() {
		iai := NewIAMAuthInfo("a")
		Expect(iai.IAMToken()).To(Equal("a"))
		iai.SetIAMToken("b")
		Expect(iai.IAMToken()).To(Equal("b"))
	})
})

var _ = Describe("DefaultWSConnectionFactory", func() {
	var (
		svr               *httptest.Server
		ch                chan struct{}
		useTLS, testError bool
	)

	happyHandler := func(ch chan struct{}) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer GinkgoRecover()
			svrOpts := ws.ConnectionOptions{}

			var upgrader websocket.Upgrader
			wc, _ := upgrader.Upgrade(w, r, nil)

			header := r.Header.Get(fclient.AuthorizationHeader)
			Expect(header).To(Equal("oi"))

			svrConnection, err := ws.NewConnection(wc, svrOpts)
			if err != nil {
				Fail("broke")
			}

			ch <- struct{}{}

			svrConnection.Close()
		})
	}

	sadHandler := func(ch chan struct{}) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "broken test", http.StatusInternalServerError)
		})
	}

	JustBeforeEach(func() {
		ch = make(chan struct{})

		if useTLS {
			svr = httptest.NewTLSServer(happyHandler(ch))
		} else if testError {
			svr = httptest.NewTLSServer(sadHandler(ch))
		} else {
			svr = httptest.NewServer(happyHandler(ch))
		}

		time.Sleep(5 * time.Millisecond)
	})

	AfterEach(func() {
		svr.Close()
	})

	It("sends auth headers", func() {
		u := "ws" + strings.TrimPrefix(svr.URL, "http")

		cli := fclient.NewWS(client.WSConnectionOptions{
			Factory: &client.DefaultWSConnectionFactory{
				URL: u,
				TLSConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
				AuthInfo: NewIAMAuthInfo("oi"),
			},
		})

		Expect(cli.Connect()).ToNot(HaveOccurred())
		Eventually(ch).Should(Receive())
		Expect(cli.Disconnect()).ToNot(HaveOccurred())
	})

	When("sends wrong url, expects error", func() {

		BeforeEach(func() {
			testError = true
		})

		It("returns error", func() {
			u := "ws" + strings.TrimPrefix(svr.URL, "http")

			cli := fclient.NewWS(client.WSConnectionOptions{
				Factory: &client.DefaultWSConnectionFactory{
					URL: u + "/test",
					TLSConfig: &tls.Config{
						InsecureSkipVerify: true,
					},
					AuthInfo: NewIAMAuthInfo("oi"),
				},
			})

			err := cli.Connect()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("websocket: bad handshake. 500:broken test"))
		})

	})

	When("the factory is configured for TLS", func() {
		BeforeEach(func() {
			useTLS = true
		})

		It("works", func() {
			u := "wss" + strings.TrimPrefix(svr.URL, "https")

			cli := fclient.NewWS(client.WSConnectionOptions{
				Factory: &client.DefaultWSConnectionFactory{
					URL: u,
					TLSConfig: &tls.Config{
						InsecureSkipVerify: true,
					},
					AuthInfo: NewIAMAuthInfo("oi"),
				},
			})

			err := cli.Connect()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})

var _ = Describe("WSClient", func() {
	var (
		factory    *clientfakes.FakeWSConnectionFactory
		client     *WSClient
		clientSide ext.Conn
		conn       *wsfakes.FakeConnection
		session    *WSSession
	)

	BeforeEach(func() {
		factory = &clientfakes.FakeWSConnectionFactory{}
		client = fclient.NewWS(fclient.WSConnectionOptions{
			Factory: factory,
		})
		clientSide = &extfakes.FakeConn{}
		conn = &wsfakes.FakeConnection{}
		session = &WSSession{Connection: conn}

		Expect(factory.NewCallCount()).To(Equal(0))
		Expect(client.Session()).To(BeNil())
	})

	JustBeforeEach(func() {
		factory.NewReturns(clientSide, nil)
		factory.NewSessionReturns(session)
	})

	Describe("Connect", func() {
		It("Does not return an error", func() {
			Expect(client.Connect()).ToNot(HaveOccurred())
		})

		It("Gets the connection from the ConnectionFactory", func() {
			err := client.Connect()
			Expect(err).NotTo(HaveOccurred())
			Expect(factory.NewCallCount()).To(Equal(1))
			Expect(factory.NewSessionCallCount()).To(Equal(1))
			Expect(client.Session()).To(Equal(session))
			Expect(client.Session().Connection).To(Equal(conn))
		})

		When("the factory returns an error", func() {
			var (
				connectionError error
			)

			JustBeforeEach(func() {
				connectionError = errors.New("Nope")
				factory.NewReturns(nil, connectionError)
			})

			It("Returns an error", func() {
				err := client.Connect()
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeIdenticalTo(connectionError))
			})
		})

	})

	Describe("Disconnect", func() {
		When("the session is not nil", func() {
			JustBeforeEach(func() {
				err := client.Connect()
				Expect(err).NotTo(HaveOccurred())
				time.Sleep(100 * time.Millisecond)
			})

			It("closes the connection", func() {
				Expect(client.Disconnect()).ToNot(HaveOccurred())
				Expect(conn.CloseCallCount()).To(Equal(1))
			})
		})

		// When("the session is nil", func() {
		// 	JustBeforeEach(func() {
		// 		client.Session = nil
		// 	})

		// 	It("does not error or panic", func() {
		// 		Expect(func() {
		// 			Expect(client.Disconnect()).ToNot(HaveOccurred())
		// 		}).ToNot(Panic())
		// 	})
		// })
	})

	Describe("Reconnect", func() {
		JustBeforeEach(func() {
			err := client.Connect()
			Expect(err).NotTo(HaveOccurred())
			time.Sleep(100 * time.Millisecond)
		})

		It("calls Disconnect and creates a new Session", func() {
			Expect(client.Reconnect()).ToNot(HaveOccurred())

			Expect(conn.CloseCallCount()).To(Equal(1))

			Expect(factory.NewSessionCallCount()).To(Equal(2))
			Expect(client.Session().Connection).ToNot(BeNil())
		})

		It("works if no active session", func() {
			Expect(client.Disconnect()).ToNot(HaveOccurred())
			Expect(conn.CloseCallCount()).To(Equal(1))

			Expect(client.Reconnect()).ToNot(HaveOccurred())
			Expect(conn.CloseCallCount()).To(Equal(1))

			Expect(factory.NewSessionCallCount()).To(Equal(2))
			Expect(client.Session().Connection).ToNot(BeNil())
		})
	})

	Describe("Send", func() {
		var (
			msg protocol.MessageExt
		)

		BeforeEach(func() {
			msg = protocol.MessageExt{
				Tag:       "foo.bar",
				Timestamp: protocol.EventTime{time.Now()}, //nolint
				Record:    map[string]interface{}{},
				Options:   &protocol.MessageOptions{},
			}
		})

		JustBeforeEach(func() {
			err := client.Connect()
			Expect(err).ToNot(HaveOccurred())
			time.Sleep(100 * time.Millisecond)
		})

		It("Sends the message", func() {
			bits, _ := msg.MarshalMsg(nil)
			Expect(client.Send(&msg)).ToNot(HaveOccurred())

			writtenbits := conn.WriteArgsForCall(0)
			Expect(bytes.Equal(bits, writtenbits)).To(BeTrue())
		})

		It("Sends the message", func() {
			msgBytes, _ := msg.MarshalMsg(nil)
			Expect(client.Send(&msg)).ToNot(HaveOccurred())

			writtenBytes := conn.WriteArgsForCall(0)
			Expect(bytes.Equal(msgBytes, writtenBytes)).To(BeTrue())
		})

		When("The message is large", func() {
			const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

			var (
				expectedBytes int
				messageSize   = 65536
			)

			JustBeforeEach(func() {
				seededRand := rand.New(
					rand.NewSource(time.Now().UnixNano()))
				m := make([]byte, messageSize)
				for i := range m {
					m[i] = charset[seededRand.Intn(len(charset))]
				}
				msg.Record = m

				var b bytes.Buffer
				Expect(msgp.Encode(&b, &msg)).ToNot(HaveOccurred())
				expectedBytes = len(b.Bytes())
			})

			It("Sends the correct number of bits", func() {
				Expect(client.Send(&msg)).ToNot(HaveOccurred())
				Expect(conn.WriteCallCount()).To(Equal(1))
				writtenBytes := len(conn.WriteArgsForCall(0))
				Expect(writtenBytes).To(Equal(expectedBytes))
			})
		})

		When("the connection is disconnected", func() {
			JustBeforeEach(func() {
				err := client.Disconnect()
				Expect(err).ToNot(HaveOccurred())
			})

			It("returns an error", func() {
				Expect(client.Send(&msg)).To(MatchError("no active session"))
			})
		})

		When("the connection is closed with an error", func() {
			BeforeEach(func() {
				conn.ListenReturns(errors.New("BOOM"))
			})

			It("returns the error", func() {
				Expect(client.Send(&msg)).To(MatchError("BOOM"))
			})
		})
	})

	Describe("SendRaw", func() {
		var (
			bits []byte
		)

		BeforeEach(func() {
			bits = []byte("oi")
		})

		JustBeforeEach(func() {
			err := client.Connect()
			Expect(err).ToNot(HaveOccurred())
			time.Sleep(100 * time.Millisecond)
		})

		It("Sends the message", func() {
			Expect(client.SendRaw(bits)).ToNot(HaveOccurred())

			writtenbits := conn.WriteArgsForCall(0)
			Expect(bytes.Equal(bits, writtenbits)).To(BeTrue())
		})

		When("the connection is disconnected", func() {
			JustBeforeEach(func() {
				err := client.Disconnect()
				Expect(err).ToNot(HaveOccurred())
			})

			It("returns an error", func() {
				Expect(client.SendRaw(bits)).To(MatchError("no active session"))
			})
		})

		When("the connection is closed with an error", func() {
			BeforeEach(func() {
				conn.ListenReturns(errors.New("BOOM"))
			})

			It("returns the error", func() {
				Expect(client.SendRaw(bits)).To(MatchError("BOOM"))
			})
		})
	})
})
