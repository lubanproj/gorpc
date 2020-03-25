package log

import (
	"testing"
)

func TestLog(t *testing.T) {
	Trace("test....")
	Tracef("testTracef...")
	Debug("test....")
	Debugf("testDebugf...")
	Info("test....")
	Infof("testInfof...")
	Warning("test....")
	Warningf("testWarningf...")
	Error("test....")
	Errorf("testErrorf...")
	Fatal("test....")
	Fatalf("testFatalf...")
}