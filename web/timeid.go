package main

import (
	"fmt"
	"launchpad.net/gozk/zookeeper"
	"sort"
	"strings"
	"time"
)

type TimeID struct {
	lastID int64
}

// Public message id-creater
var PubMID *TimeID

// NewTimeID create a new TimeID struct
func NewTimeID() *TimeID {
	return &TimeID{lastID: 0}
}

// ID generate a time ID
func (t *TimeID) ID() int64 {
	for {
		s := time.Now().UnixNano() / 100
		if t.lastID == s {
			time.Sleep(100 * time.Nanosecond)
		} else {
			return s
		}
	}
	return 0
}

// PubMIDLock public message mid lock, make sure that get the unique mid
func (t *TimeID) Lock() (bool, string, error) {
	prefix := "p"
	splitSign := "@"
	pathCreated, err := zk.Create(fmt.Sprintf("%s/%s%s", Conf.ZKCometPath, prefix, splitSign),
		"0", zookeeper.EPHEMERAL|zookeeper.SEQUENCE, zookeeper.WorldACL(zookeeper.PERM_ALL))
	if err != nil {
		return false, pathCreated, fmt.Errorf("zk.Create(%s/%s%s) error(%v)", Conf.ZKCometPath, prefix, splitSign, err)
	}

	for {
		childrens, stat, err := zk.Children(Conf.ZKCometPath)
		if err != nil {
			return false, pathCreated, fmt.Errorf("zk.Children(%s) error(%v)", Conf.ZKCometPath, err)
		}

		// If node isn`t exist
		if childrens == nil || stat == nil {
			return false, pathCreated, fmt.Errorf("node(%s) is not existent", Conf.ZKCometPath)
		}

		var realQueue []string
		for _, children := range childrens {
			tmp := strings.Split(children, splitSign)
			if prefix == tmp[0] {
				realQueue = append(realQueue, children)
			}
		}

		// Sort sequence nodes
		sort.Strings(realQueue)

		tmp := strings.Split(pathCreated, "/")
		posReal := sort.StringSlice(realQueue).Search(tmp[len(tmp)-1])

		// If does not get lock
		if posReal > 0 {
			// Watch the last one
			watchPath := fmt.Sprintf("%s/%s", Conf.ZKCometPath, realQueue[posReal-1])
			_, watch, err := zk.ExistsW(watchPath)
			if err != nil || zookeeper.IsError(err, zookeeper.ZNONODE) {
				return false, pathCreated, fmt.Errorf("zk.ExistsW(%s) error(%v) or no node", watchPath, err)
			}

			// Watch the lower node
			watchNode := <-watch
			switch watchNode.Type {
			case zookeeper.EVENT_DELETED:
			default:
				return false, pathCreated, fmt.Errorf("zookeeper watch errCode:%d", watchNode.Type)
			}

			return false, pathCreated, nil
		} else {
			return true, pathCreated, nil
		}
	}

	// Never get here
	return false, pathCreated, fmt.Errorf("never get here")
}

// PubMIDLockRelease release the public message id-lock
func (t *TimeID) LockRelease(pathCreated string) error {
	if err := zk.Delete(pathCreated, -1); err != nil {
		return err
	}
	return nil
}
