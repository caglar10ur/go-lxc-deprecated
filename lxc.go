/*
 * lxc.go: Go bindings for lxc
 *
 * Copyright © 2013, S.Çağlar Onur
 *
 * Authors:
 * S.Çağlar Onur <caglar@10ur.org>
 *
 * This library is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 2, as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License along
 * with this program; if not, write to the Free Software Foundation, Inc.,
 * 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
 */

//Go (golang) Bindings for LXC (Linux Containers)
//
//This package implements Go bindings for the LXC C API.
package lxc

// #cgo linux LDFLAGS: -llxc -lutil
// #include <lxc/lxc.h>
// #include <lxc/lxccontainer.h>
// #include "lxc.h"
import "C"

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

type State int

const (
	// Timeout
	WAIT_FOREVER int = iota
	DONT_WAIT
	// State
	STOPPED  = C.STOPPED
	STARTING = C.STARTING
	RUNNING  = C.RUNNING
	STOPPING = C.STOPPING
	ABORTING = C.ABORTING
	FREEZING = C.FREEZING
	FROZEN   = C.FROZEN
	THAWED   = C.THAWED
)

var stateMap = map[string]State{
	"STOPPED":  STOPPED,
	"STARTING": STARTING,
	"RUNNING":  RUNNING,
	"STOPPING": STOPPING,
	"ABORTING": ABORTING,
	"FREEZING": FREEZING,
	"FROZEN":   FROZEN,
	"THAWED":   THAWED,
}

// State as string
func (t State) String() string {
	switch t {
	case STOPPED:
		return "STOPPED"
	case STARTING:
		return "STARTING"
	case RUNNING:
		return "RUNNING"
	case STOPPING:
		return "STOPPING"
	case ABORTING:
		return "ABORTING"
	case FREEZING:
		return "FREEZING"
	case FROZEN:
		return "FROZEN"
	case THAWED:
		return "THAWED"
	}
	return "<INVALID>"
}

func makeArgs(args []string) []*C.char {
	ret := make([]*C.char, len(args))
	for i, s := range args {
		ret[i] = C.CString(s)
	}
	return ret
}

func freeArgs(cArgs []*C.char) {
	for _, s := range cArgs {
		C.free(unsafe.Pointer(s))
	}
}

type Container struct {
	container *C.struct_lxc_container
}

func (lxc *Container) Error() string {
	return C.GoString(lxc.container.error_string)
}

func (lxc *Container) GetError() error {
	return syscall.Errno(int(lxc.container.error_num))
}

func NewContainer(name string) Container {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return Container{C.lxc_container_new(cname)}
}

// Returns container's name
func (lxc *Container) GetName() string {
	return C.GoString(lxc.container.name)
}

// Returns whether the container is already defined or not
func (lxc *Container) Defined() bool {
	return bool(C.lxc_container_defined(lxc.container))
}

// Returns whether the container is already running or not
func (lxc *Container) Running() bool {
	return bool(C.lxc_container_running(lxc.container))
}

// Returns the container's state
func (lxc *Container) GetState() State {
	return stateMap[C.GoString(C.lxc_container_state(lxc.container))]
}

// Returns the container's PID
func (lxc *Container) GetInitPID() int {
	return int(C.lxc_container_init_pid(lxc.container))
}

// Returns whether the daemonize flag is set
func (lxc *Container) GetDaemonize() bool {
	return bool(lxc.container.daemonize != 0)
}

// Sets the daemonize flag
func (lxc *Container) SetDaemonize() {
	C.lxc_container_want_daemonize(lxc.container)
}

// Freezes the running container
func (lxc *Container) Freeze() bool {
	return bool(C.lxc_container_freeze(lxc.container))
}

// Unfreezes the frozen container
func (lxc *Container) Unfreeze() bool {
	return bool(C.lxc_container_unfreeze(lxc.container))
}

// Creates the container using given template and arguments
func (lxc *Container) Create(template string, args []string) bool {
	ctemplate := C.CString(template)
	defer C.free(unsafe.Pointer(ctemplate))
	if args != nil {
		cargs := makeArgs(args)
		defer freeArgs(cargs)
		return bool(C.lxc_container_create(lxc.container, ctemplate, &cargs[0]))
	}
	return bool(C.lxc_container_create(lxc.container, ctemplate, nil))
}

// Starts the container
func (lxc *Container) Start(useinit bool, args []string) bool {
	cuseinit := 0
	if useinit {
		cuseinit = 1
	}
	if args != nil {
		cargs := makeArgs(args)
		defer freeArgs(cargs)
		return bool(C.lxc_container_start(lxc.container, C.int(cuseinit), &cargs[0]))
	}
	return bool(C.lxc_container_start(lxc.container, C.int(cuseinit), nil))
}

// Stops the container
func (lxc *Container) Stop() bool {
	return bool(C.lxc_container_stop(lxc.container))
}

// Shutdowns the container
func (lxc *Container) Shutdown(timeout int) bool {
	return bool(C.lxc_container_shutdown(lxc.container, C.int(timeout)))
}

// Destroys the container
func (lxc *Container) Destroy() bool {
	return bool(C.lxc_container_destroy(lxc.container))
}

// Waits till the container changes its state or timeouts
func (lxc *Container) Wait(state State, timeout int) bool {
	cstate := C.CString(state.String())
	defer C.free(unsafe.Pointer(cstate))
	return bool(C.lxc_container_wait(lxc.container, cstate, C.int(timeout)))
}

// Returns the container's configuration file's name
func (lxc *Container) GetConfigFileName() string {
	return C.GoString(C.lxc_container_config_file_name(lxc.container))
}

// Returns the value of the given key
func (lxc *Container) GetConfigItem(key string) []string {
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	return strings.Split(C.GoString(C.lxc_container_get_config_item(lxc.container, ckey)), "\n")
}

// Sets the value of given key
func (lxc *Container) SetConfigItem(key string, value string) bool {
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))
	return bool(C.lxc_container_set_config_item(lxc.container, ckey, cvalue))
}

// Clears the value of given key
func (lxc *Container) ClearConfigItem(key string) bool {
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	return bool(C.lxc_container_clear_config_item(lxc.container, ckey))
}

// Returns the keys
func (lxc *Container) GetKeys(key string) []string {
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	return strings.Split(C.GoString(C.lxc_container_get_keys(lxc.container, ckey)), "\n")
}

// Loads the configuration file from given path
func (lxc *Container) LoadConfigFile(path string) bool {
	// TODO: Remove following code from binding as patch sent to lxc-devel
	// http://sourceforge.net/mailarchive/forum.php?thread_name=1364411217-15616-1-git-send-email-caglar%4010ur.org&forum_name=lxc-devel
	// reject loading config file if it doesn't exist
	// otherwise container starts without a netns
	path, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	return bool(C.lxc_container_load_config(lxc.container, cpath))
}

// Saves the configuration file to given path
func (lxc *Container) SaveConfigFile(path string) bool {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	return bool(C.lxc_container_save_config(lxc.container, cpath))
}
