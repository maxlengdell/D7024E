//package go1
package d7024e

import "testing"
import "math"
import "net"
import "reflect"
import "fmt"

type Go1 testing.T

func (g *Go1) AssertEquals(expected, value interface{}) bool {
  if !reflect.DeepEqual(expected, value) {
    fmt.Printf("%v != %v (expected %v, got %v)\n", expected, value, expected, value)
	g.Fail()
    return false
  }
  return true
}

func AssertEquals(t *testing.T, expected, value interface{}) bool {
  g := Go1(*t)
  return g.AssertEquals(expected, value)
}

func AssertNotEquals(t *testing.T, expected, value interface{}) bool {
  if reflect.DeepEqual(expected, value) {
    fmt.Printf("%v == %v (got unexpected %v)\n", expected, value, value)
	t.Fail()
    return false
  }
  return true
}

func AssertEqualsApprox(t *testing.T, expected, value float64, delta float64) bool {
  diff := math.Abs(expected - value)
  if diff > delta {
    fmt.Printf("%v != %v (+/- %f)", expected, value, delta)
    t.Fail()
    return false
  }
  return true
}

func AssertTrue(t *testing.T, value bool) bool {
  g := Go1(*t)
  return g.AssertEquals(true, value)
}

func AssertFalse(t *testing.T, value bool) bool {
  return AssertTrue(t, !value)
}

func AssertNoError(t *testing.T, f func(... interface{}) (interface{}, error), args ... interface{}) interface{} {
  result, err := f(args...)
  AssertTrue(t, err == nil)
  return result
}

// AssertReceived takes an address (of the form IP:PORT)
// and then checks that the expected bytes are received from that connection.
func AssertReceived(t *testing.T, addr string, expected []byte) bool {
  udpAddr, err := net.ResolveUDPAddr("udp", addr)
  if !AssertTrue(t, err == nil) {
    t.Errorf("Error not nil: %v", err)
    return false
  }
  conn, err := net.ListenUDP("udp", udpAddr)
  if !AssertTrue(t, err == nil) {
    t.Errorf("Error not nil: %v", err)
    return false
  }
  var buf [2000]byte
  n, _, err := conn.ReadFromUDP(buf[0:])
  if !AssertTrue(t, err == nil) {
    t.Errorf("Error not nil: %v", err)
    return false
  }
  if !AssertEquals(t, len(expected), n) {
    t.Errorf("Length not equal: %d != %d", len(expected), n)
    return false
  }
  if !AssertEquals(t, expected, buf[:n]) {
    t.Errorf("Buffers not equal: %v != %v", expected, buf[:n])
    return false
  }
  return true
}

// AssertReceivedResponse takes an address (of the form IP:PORT), sends the
// given message msg, and then checks that the expected reply is received.
func AssertReceivedResponse(t *testing.T, addr string, msg []byte, expected []byte) bool {
  udpAddr, err := net.ResolveUDPAddr("udp", addr)
  if !AssertTrue(t, err == nil) {
    t.Errorf("Error not nil: %v", err)
    return false
  }
  conn, err := net.DialUDP("udp", nil, udpAddr)
  if !AssertTrue(t, err == nil) {
    t.Errorf("Error not nil: %v", err)
    return false
  }
  var buf [2000]byte
  n, err := conn.Write(msg)
  if !AssertTrue(t, err == nil) {
    t.Errorf("Error not nil: %v", err)
    return false
  }
  if !AssertEquals(t, len(msg), n) {
    t.Errorf("Failed to send the message. (Sent %d out of %d bytes.)", n, len(msg))
    return false
  }
  n, _, err = conn.ReadFromUDP(buf[0:])
  if !AssertTrue(t, err == nil) {
    t.Errorf("Error not nil: %v", err)
    return false
  }
  if !AssertEquals(t, len(expected), n) {
    t.Errorf("Length not equal: %d != %d", len(expected), n)
    return false
  }
  if !AssertEquals(t, expected, buf[:n]) {
    t.Errorf("Buffers not equal: %v != %v", expected, buf[:n])
    return false
  }
  return true
}
