package syswriter

import (
	"bytes"
	"fmt"
	"log"
	"log/syslog"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type SysWriter struct {
	*syslog.Writer
	closed int32
	opt    option
	stdlog *log.Logger
	ch     chan *bytes.Buffer
	pool   *sync.Pool
	wg     sync.WaitGroup
}

func New(writer *syslog.Writer, wns ...Option) (*SysWriter, error) {
	opt := defaultOption
	for _, fn := range wns {
		fn(&opt)
	}

	stdlog := log.New(os.Stderr, "flog ", log.LstdFlags)
	ch := make(chan *bytes.Buffer, opt.ChanSize)
	sw := &SysWriter{
		Writer: writer,
		ch:     ch,
		stdlog: stdlog,
		opt:    opt,
		pool:   &sync.Pool{New: func() interface{} { return new(bytes.Buffer) }},
	}

	if opt.Batch {
		sw.wg.Add(1)
		go sw.daemon()
	}
	return sw, nil
}

func (s *SysWriter) write(p string) (err error) {
	switch s.opt.Priority {
	case syslog.LOG_EMERG:
		return s.Writer.Emerg(p)
	case syslog.LOG_ALERT:
		return s.Writer.Alert(p)
	case syslog.LOG_CRIT:
		return s.Writer.Crit(p)
	case syslog.LOG_ERR:
		return s.Writer.Err(p)
	case syslog.LOG_WARNING:
		return s.Writer.Warning(p)
	case syslog.LOG_NOTICE:
		return s.Writer.Notice(p)
	case syslog.LOG_INFO:
		return s.Writer.Info(p)
	case syslog.LOG_DEBUG:
		return s.Writer.Debug(p)
	default:
		_, err = s.Writer.Write([]byte(p))
		return err
	}
}

func (s *SysWriter) daemon() {
	// TODO: check aggsbuf size prevent it too big
	aggsbuf := &bytes.Buffer{}
	// TODO: make it configrable
	aggstk := time.NewTicker(10 * time.Millisecond)
	var err error
	for {
		select {
		case buf, ok := <-s.ch:
			if ok {
				aggsbuf.Write(buf.Bytes())
				s.putBuf(buf)
			}
		case <-aggstk.C:
			if aggsbuf.Len() > 0 {
				if err = s.write(aggsbuf.String()); err != nil {
					s.stdlog.Printf("write log error: %s", err)
				}
				aggsbuf.Reset()
			}
		}
		if atomic.LoadInt32(&s.closed) != 1 {
			continue
		}
		// read all buf from channel and break loop
		if err = s.write(aggsbuf.String()); err != nil {
			s.stdlog.Printf("write log error: %s", err)
		}
		for buf := range s.ch {
			if err = s.write(buf.String()); err != nil {
				s.stdlog.Printf("write log error: %s", err)
			}
			s.putBuf(buf)
		}
		break
	}
	s.wg.Done()
	return
}

func (s *SysWriter) Write(p []byte) (n int, err error) {
	if !s.opt.Batch {
		s.write(string(p))
		return
	}

	// atomic is not necessary
	if atomic.LoadInt32(&s.closed) == 1 {
		return 0, fmt.Errorf("syswriter already closed")
	}

	// because write to file is asynchronousc,
	// copy p to internal buf prevent p be change on outside
	buf := s.getBuf()
	buf.Write(p)

	if s.opt.WriteTimeout == 0 {
		select {
		case s.ch <- buf:
			return len(p), nil
		default:
			// TODO: write discard log to to stdout?
			return 0, fmt.Errorf("log channel is full, discard log")
		}
	}

	// write log with timeout
	timeout := time.NewTimer(s.opt.WriteTimeout)
	select {
	case s.ch <- buf:
		return len(p), nil
	case <-timeout.C:
		// TODO: write discard log to to stdout?
		return 0, fmt.Errorf("log channel is full, discard log")
	}

	return
}

func (s *SysWriter) Close() error {
	if s.opt.Batch {
		atomic.StoreInt32(&s.closed, 1)
		close(s.ch)
		s.wg.Wait()
	}

	return nil
}

func (s *SysWriter) putBuf(buf *bytes.Buffer) {
	buf.Reset()
	s.pool.Put(buf)
}

func (s *SysWriter) getBuf() *bytes.Buffer {
	return s.pool.Get().(*bytes.Buffer)
}
