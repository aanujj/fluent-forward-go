// Code generated by counterfeiter. DO NOT EDIT.
package clientfakes

import (
	"sync"

	"github.com/aanujj/fluent-forward-go/fluent/client"
	"github.com/aanujj/fluent-forward-go/fluent/client/ws"
	"github.com/aanujj/fluent-forward-go/fluent/client/ws/ext"
)

type FakeClientFactory struct {
	NewStub        func() (ext.Conn, error)
	newMutex       sync.RWMutex
	newArgsForCall []struct {
	}
	newReturns struct {
		result1 ext.Conn
		result2 error
	}
	newReturnsOnCall map[int]struct {
		result1 ext.Conn
		result2 error
	}
	NewSessionStub        func(ws.Connection) *client.WSSession
	newSessionMutex       sync.RWMutex
	newSessionArgsForCall []struct {
		arg1 ws.Connection
	}
	newSessionReturns struct {
		result1 *client.WSSession
	}
	newSessionReturnsOnCall map[int]struct {
		result1 *client.WSSession
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeClientFactory) New() (ext.Conn, error) {
	fake.newMutex.Lock()
	ret, specificReturn := fake.newReturnsOnCall[len(fake.newArgsForCall)]
	fake.newArgsForCall = append(fake.newArgsForCall, struct {
	}{})
	stub := fake.NewStub
	fakeReturns := fake.newReturns
	fake.recordInvocation("New", []interface{}{})
	fake.newMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClientFactory) NewCallCount() int {
	fake.newMutex.RLock()
	defer fake.newMutex.RUnlock()
	return len(fake.newArgsForCall)
}

func (fake *FakeClientFactory) NewCalls(stub func() (ext.Conn, error)) {
	fake.newMutex.Lock()
	defer fake.newMutex.Unlock()
	fake.NewStub = stub
}

func (fake *FakeClientFactory) NewReturns(result1 ext.Conn, result2 error) {
	fake.newMutex.Lock()
	defer fake.newMutex.Unlock()
	fake.NewStub = nil
	fake.newReturns = struct {
		result1 ext.Conn
		result2 error
	}{result1, result2}
}

func (fake *FakeClientFactory) NewReturnsOnCall(i int, result1 ext.Conn, result2 error) {
	fake.newMutex.Lock()
	defer fake.newMutex.Unlock()
	fake.NewStub = nil
	if fake.newReturnsOnCall == nil {
		fake.newReturnsOnCall = make(map[int]struct {
			result1 ext.Conn
			result2 error
		})
	}
	fake.newReturnsOnCall[i] = struct {
		result1 ext.Conn
		result2 error
	}{result1, result2}
}

func (fake *FakeClientFactory) NewSession(arg1 ws.Connection) *client.WSSession {
	fake.newSessionMutex.Lock()
	ret, specificReturn := fake.newSessionReturnsOnCall[len(fake.newSessionArgsForCall)]
	fake.newSessionArgsForCall = append(fake.newSessionArgsForCall, struct {
		arg1 ws.Connection
	}{arg1})
	stub := fake.NewSessionStub
	fakeReturns := fake.newSessionReturns
	fake.recordInvocation("NewSession", []interface{}{arg1})
	fake.newSessionMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeClientFactory) NewSessionCallCount() int {
	fake.newSessionMutex.RLock()
	defer fake.newSessionMutex.RUnlock()
	return len(fake.newSessionArgsForCall)
}

func (fake *FakeClientFactory) NewSessionCalls(stub func(ws.Connection) *client.WSSession) {
	fake.newSessionMutex.Lock()
	defer fake.newSessionMutex.Unlock()
	fake.NewSessionStub = stub
}

func (fake *FakeClientFactory) NewSessionArgsForCall(i int) ws.Connection {
	fake.newSessionMutex.RLock()
	defer fake.newSessionMutex.RUnlock()
	argsForCall := fake.newSessionArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeClientFactory) NewSessionReturns(result1 *client.WSSession) {
	fake.newSessionMutex.Lock()
	defer fake.newSessionMutex.Unlock()
	fake.NewSessionStub = nil
	fake.newSessionReturns = struct {
		result1 *client.WSSession
	}{result1}
}

func (fake *FakeClientFactory) NewSessionReturnsOnCall(i int, result1 *client.WSSession) {
	fake.newSessionMutex.Lock()
	defer fake.newSessionMutex.Unlock()
	fake.NewSessionStub = nil
	if fake.newSessionReturnsOnCall == nil {
		fake.newSessionReturnsOnCall = make(map[int]struct {
			result1 *client.WSSession
		})
	}
	fake.newSessionReturnsOnCall[i] = struct {
		result1 *client.WSSession
	}{result1}
}

func (fake *FakeClientFactory) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.newMutex.RLock()
	defer fake.newMutex.RUnlock()
	fake.newSessionMutex.RLock()
	defer fake.newSessionMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeClientFactory) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ client.WSConnectionFactory = new(FakeClientFactory)
