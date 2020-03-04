package raw

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows/registry"
)

const (
	tapDriverKey = `SYSTEM\CurrentControlSet\Control\Class\{4D36E972-E325-11CE-BFC1-08002BE10318}`
	netConfigKey = `SYSTEM\CurrentControlSet\Control\Network\{4D36E972-E325-11CE-BFC1-08002BE10318}`
)

type Tap struct {
	fd   syscall.Handle
	name string
}

//If you want to use raw package on windows You must installed tap-win32
func NewTap(name string) (*Tap, error) {
	// find the device in registry.
	deviceid, err := getdeviceid(
		"tap0901",
		name)
	if err != nil {
		return nil, err
	}
	path := "\\\\.\\Global\\" + deviceid + ".tap"
	pathp, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return nil, err
	}
	// type Handle uintptr
	file, err := syscall.CreateFile(
		pathp,
		syscall.GENERIC_READ|syscall.GENERIC_WRITE,
		uint32(syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE),
		nil,
		syscall.OPEN_EXISTING,
		syscall.FILE_ATTRIBUTE_SYSTEM, //|syscall.FILE_FLAG_OVERLAPPED,
		0)
	// if err hanppens, close the interface.
	defer func() {
		if err != nil {
			syscall.Close(file)
		}
		if err := recover(); err != nil {
			syscall.Close(file)
		}
	}()
	if err != nil {
		return nil, err
	}
	var bytesReturned uint32

	// find the mac address of tap device, use this to find the name of interface
	mac := make([]byte, 6)
	err = syscall.DeviceIoControl(
		file,
		//		tap_win_ioctl_get_mac,
		uint32(0x00220004),
		&mac[0],
		uint32(len(mac)),
		&mac[0],
		uint32(len(mac)),
		&bytesReturned,
		nil)
	if err != nil {
		return nil, err
	}

	// bring up device.
	rdbbuf := make([]byte, syscall.MAXIMUM_REPARSE_DATA_BUFFER_SIZE)
	code := []byte{0x01, 0x00, 0x00, 0x00}

	if err := syscall.DeviceIoControl(
		file,
		uint32(0x00220018),
		&code[0],
		uint32(4),
		&rdbbuf[0],
		uint32(len(rdbbuf)),
		&bytesReturned, nil); err != nil {
		log.Fatal(err)
	}

	tap := &Tap{fd: file}
	// find the name of tap interface(u need it to set the ip or other command)
	ifces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, v := range ifces {
		if len(v.HardwareAddr) < 6 {
			continue
		}
		if bytes.Equal(v.HardwareAddr[:6], mac[:6]) {
			tap.name = v.Name
			return tap, nil
		}
	}
	return nil, fmt.Errorf("can't found tap interface")
}

// getdeviceid finds out a TAP device from registry, it *may* requires privileged right to prevent some weird issue.
func getdeviceid(componentID string, interfaceName string) (deviceid string, err error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, tapDriverKey, registry.READ)
	if err != nil {
		return "", fmt.Errorf("Failed to open the adapter registry, TAP driver may be not installed, %v", err)
	}
	defer k.Close()
	// read all subkeys, it should not return an err here
	keys, err := k.ReadSubKeyNames(-1)
	if err != nil {
		return "", err
	}
	// find the one matched ComponentId
	for _, v := range keys {
		key, err := registry.OpenKey(registry.LOCAL_MACHINE, tapDriverKey+"\\"+v, registry.READ)
		if err != nil {
			continue
		}
		val, _, err := key.GetStringValue("ComponentId")
		if err != nil {
			key.Close()
			continue
		}
		if val == componentID {
			val, _, err = key.GetStringValue("NetCfgInstanceId")
			if err != nil {
				key.Close()
				continue
			}
			if len(interfaceName) > 0 {
				key2 := fmt.Sprintf("%s\\%s\\Connection", netConfigKey, val)
				k2, err := registry.OpenKey(registry.LOCAL_MACHINE, key2, registry.READ)
				if err != nil {
					continue
				}
				defer k2.Close()
				val, _, err := k2.GetStringValue("Name")
				if err != nil || val != interfaceName {
					continue
				}
			}
			key.Close()
			return val, nil
		}
		key.Close()
	}
	if len(interfaceName) > 0 {
		return "", fmt.Errorf("Failed to find the tap device in registry with specified ComponentId '%s' and InterfaceName '%s', TAP driver may be not installed or you may have specified an interface name that doesn't exist", componentID, interfaceName)
	}

	return "", fmt.Errorf("Failed to find the tap device in registry with specified ComponentId '%s', TAP driver may be not installed", componentID)
}

func (s *Tap) Read(buf []byte) (int, error) {
	return syscall.Read(s.fd, buf)
}

func (s *Tap) Write(buf []byte) (int, error) {
	return syscall.Write(s.fd, buf)
}

func (s *Tap) Close() error {
	return syscall.CloseHandle(s.fd)
}

func (s *Tap) Addr() ([]byte, error) {
	ap, err := syscall.GetProcAddress(s.fd, s.name)
	if err != nil {
		return nil, err
	}
	return (*[unsafe.Sizeof(unsafe.Pointer(0))]byte)(unsafe.Pointer(ap))[:], nil
}

func (s *Tap) Name() string {
	return s.name
}
