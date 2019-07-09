package collect

import (
	"context"
	"fmt"
	"time"
)

// Collector collect all thing
type Collector struct {
	buf           chan interface{}
	outputCh      chan []interface{}
	maxOutputSize int
	timeout       time.Duration
	ctx           context.Context
	stopFunc      func()
	stopFlag      chan bool
}

// NewCollector create a data collector, return msg list if collect
// enough msg or enough time.
// Param:
// bufferSize - input buffer size,
// maxOutputSize - max size of ONE msg list return by Get(),
// timeout - longest time to collect.
func NewCollector(bufferSize int, maxOutputSize int, timeout time.Duration) *Collector {
	ctx, stopFunc := context.WithCancel(context.Background())

	ret := &Collector{
		buf:           make(chan interface{}, bufferSize),
		outputCh:      make(chan []interface{}),
		maxOutputSize: maxOutputSize,
		timeout:       timeout,
		ctx:           ctx,
		stopFunc:      stopFunc,
		stopFlag:      make(chan bool),
	}
	ret.goCollect()

	return ret
}

func (c *Collector) goCollect() {
	go func() {
	WORKING_STATUS:
		for {
			var msgList []interface{}

		TRY_GET_FIRST_MSG:
			for {
				select {
				case msg := <-c.buf:
					msgList = append(msgList, msg)
					break TRY_GET_FIRST_MSG
				default:
					select {
					case <-c.ctx.Done():
						break WORKING_STATUS
					case msg := <-c.buf:
						msgList = append(msgList, msg)
						break TRY_GET_FIRST_MSG
					}
				}
			}

			timer := time.NewTimer(c.timeout)
		IN_TIME:
			for len(msgList) < c.maxOutputSize {
				// priority: timer > msg buf > ctx done
				select {
				case <-timer.C:
					break IN_TIME
				default:
					select {
					case <-timer.C:
						break IN_TIME
					case msg := <-c.buf:
						msgList = append(msgList, msg)
					default:
						select {
						case <-c.ctx.Done():
							timer.Stop()
							break WORKING_STATUS
						case <-timer.C:
							break IN_TIME
						case msg := <-c.buf:
							msgList = append(msgList, msg)
						}
					}
				}
			}

			c.outputCh <- msgList
		}

		c.stopFlag <- true
	}()
}

// Add a msg to buffer.
// Return nil if success, return err if collector already stopped or
// buffer is full.
func (c *Collector) Add(msg interface{}) error {
	select {
	case <-c.ctx.Done():
		return fmt.Errorf("collector already stopped")
	default:
		break
	}

	select {
	case c.buf <- msg:
		return nil
	default:
		return fmt.Errorf("buf is full")
	}
}

// Get msg List from collector.
// This func will blocking until get msg list. Returning msg list's size
// always greater than 0. Only if this collector is stopped and already
// get rest of msg, will return err.
func (c *Collector) Get() ([]interface{}, error) {
	select {
	case ret := <-c.outputCh:
		return ret, nil
	default:
		select {
		case ret := <-c.outputCh:
			return ret, nil
		case <-c.stopFlag:
			return nil, fmt.Errorf("collector already stopped")
		}
	}
}

// GetNonBlock try to get msg List from collector.
// This func will NOT blocking if get no msg. Returning msg list's size == 0
// if no msg. Only if this collector is stopped and already get rest of msg,
// will return err.
func (c *Collector) GetNonBlock() ([]interface{}, error) {
	select {
	case ret := <-c.outputCh:
		return ret, nil
	default:
		select {
		case ret := <-c.outputCh:
			return ret, nil
		case <-c.stopFlag:
			return nil, fmt.Errorf("collector already stopped")
		default:
			return nil, nil
		}
	}
}

// Stop the collector.
func (c *Collector) Stop() {
	c.stopFunc()
}
