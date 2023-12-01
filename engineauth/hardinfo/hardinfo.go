package hardinfo

import (
	"fmt"
	"net"
	"os"

	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/StackExchange/wmi"
	"golang.org/x/sys/windows"
)

var (
	//advapi = syscall.NewLazyDLL("Advapi32.dll")
	kernel = syscall.NewLazyDLL("Kernel32.dll")
)

func PrintHardInfo() {
	fmt.Printf("%s:%s\n", "主板信息", GetMotherboardInfo())
	fmt.Printf("%s:%s\n", "Bios信息", GetBiosInfo())
	fmt.Printf("%s:%s\n", "CPU信息", GetCpuInfo())
	fmt.Printf("%s:%s\n", "内存信息", GetMemory())
	fmt.Printf("%s:%s\n", "开机时间信息", GetStartTime())
	fmt.Printf("%s:%s\n", "系统版本信息", GetSystemVersion())
	fmt.Printf("%s:%s\n", "用户信息", GetUserName())
	fmt.Printf("%s\n", "硬盘信息:")
	fmt.Printf("%8s %16s %16s %16s\n", "盘符", "总空间(GB)", "剩余空间(GB)", "剩余比例(%)")
	for _, hd := range GetDiskInfo() {
		fmt.Printf("%10s %18.2f %18.2f %18.02f\n", hd.Path, float64(hd.Total)/1024/1024/1024, float64(hd.Free)/1024/1024/1024, float64(hd.Free)/float64(hd.Total)*100)
	}
	fmt.Printf("%s\n", "网卡信息:")
	for i, net := range GetIntfs() {
		fmt.Printf("----------------------------%d------------------------------\n", i)
		fmt.Printf("%10s:%d,名称:%s,是以太网:%t \n", "索引", net.Index, net.Name, net.IsEthernet)
		fmt.Printf("%10s:%s\n", "MAC地址", net.MacAddress)
		fmt.Printf("%10s:%v\n", "IPV4", net.Ipv4)
		fmt.Printf("%10s:%v\n", "IPV6", net.Ipv6)
	}
}

// 开机时间
func GetStartTime() string {
	GetTickCount := kernel.NewProc("GetTickCount")
	r, _, _ := GetTickCount.Call()
	if r == 0 {
		return ""
	}
	ms := time.Duration(r * 1000 * 1000)
	return ms.String()
}

// 当前用户名
func GetUserName() string {
	var size uint32 = 128
	var buffer = make([]uint16, size)
	user, _ := syscall.UTF16PtrFromString("USERNAME")
	domain, _ := syscall.UTF16PtrFromString("USERDOMAIN")
	r, err := syscall.GetEnvironmentVariable(user, &buffer[0], size)
	if err != nil {
		return ""
	}
	buffer[r] = '@'
	old := r + 1
	if old >= size {
		return syscall.UTF16ToString(buffer[:r])
	}
	r, err = syscall.GetEnvironmentVariable(domain, &buffer[old], size-old)
	if err != nil {
		return syscall.UTF16ToString(buffer[:old-1])
	}
	return syscall.UTF16ToString(buffer[:old+r])
}

// 系统版本
func GetSystemVersion() string {
	version, err := syscall.GetVersion()
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%d.%d (%d)", byte(version), uint8(version>>8), version>>16)
}

type diskusage struct {
	Path  string `json:"path"`
	Total uint64 `json:"total"`
	Free  uint64 `json:"free"`
}

func usage(getDiskFreeSpaceExW *syscall.LazyProc, path string) (diskusage, error) {
	lpFreeBytesAvailable := int64(0)
	var info = diskusage{Path: path}
	p, _ := syscall.UTF16PtrFromString(info.Path)
	diskret, _, err := getDiskFreeSpaceExW.Call(
		uintptr(unsafe.Pointer(p)),
		uintptr(unsafe.Pointer(&lpFreeBytesAvailable)),
		uintptr(unsafe.Pointer(&(info.Total))),
		uintptr(unsafe.Pointer(&(info.Free))))
	if diskret != 0 {
		err = nil
	}
	return info, err
}

// 硬盘信息
func GetDiskInfo() (infos []diskusage) {
	GetLogicalDriveStringsW := kernel.NewProc("GetLogicalDriveStringsW")
	GetDiskFreeSpaceExW := kernel.NewProc("GetDiskFreeSpaceExW")
	lpBuffer := make([]byte, 254)
	diskret, _, _ := GetLogicalDriveStringsW.Call(
		uintptr(len(lpBuffer)),
		uintptr(unsafe.Pointer(&lpBuffer[0])))
	if diskret == 0 {
		return
	}
	for _, v := range lpBuffer {
		if v >= 65 && v <= 90 {
			path := string(v) + ":"
			if path == "A:" || path == "B:" {
				continue
			}
			info, err := usage(GetDiskFreeSpaceExW, string(v)+":")
			if err != nil {
				continue
			}
			infos = append(infos, info)
		}
	}
	return infos
}

// CPU信息
// 第一个字符串为CPU核心数,第二个字符串为CPU类型
func GetCpuInfo() []string {
	var cpuInfo []string
	var size uint32 = 128
	var buffer = make([]uint16, size)
	nums, _ := syscall.UTF16PtrFromString("NUMBER_OF_PROCESSORS")
	arch, _ := syscall.UTF16PtrFromString("PROCESSOR_ARCHITECTURE")
	r, err := syscall.GetEnvironmentVariable(nums, &buffer[0], size)
	if err != nil {
		return cpuInfo
	}
	cpuInfo = append(cpuInfo, syscall.UTF16ToString(buffer[:r]))

	n, err := syscall.GetEnvironmentVariable(arch, &buffer[r-1], size-r)
	if err != nil {
		return cpuInfo
	}
	cpuInfo = append(cpuInfo, syscall.UTF16ToString(buffer[r-1:r+n]))
	return cpuInfo
}

type memoryStatusEx struct {
	cbSize                  uint32
	dwMemoryLoad            uint32
	ullTotalPhys            uint64 // in bytes
	ullAvailPhys            uint64
	ullTotalPageFile        uint64
	ullAvailPageFile        uint64
	ullTotalVirtual         uint64
	ullAvailVirtual         uint64
	ullAvailExtendedVirtual uint64
}

// 内存信息
func GetMemory() string {
	GlobalMemoryStatusEx := kernel.NewProc("GlobalMemoryStatusEx")
	var memInfo memoryStatusEx
	memInfo.cbSize = uint32(unsafe.Sizeof(memInfo))
	mem, _, _ := GlobalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&memInfo)))
	if mem == 0 {
		return ""
	}
	return fmt.Sprint(memInfo.ullTotalPhys / (1024 * 1024))
}

type intfInfo struct {
	Index      int    //序号
	IsEthernet bool   //是否以太网
	Name       string //名称
	MacAddress string //MAC地址
	Ipv4       []string
	Ipv6       []string
}

// 网卡信息
func GetIntfs() []intfInfo {
	intf, err := net.Interfaces()
	if err != nil {
		return []intfInfo{}
	}
	aas, _ := adapterAddresses()

	var itfs []intfInfo
	for _, v := range intf {
		ips, err := v.Addrs()
		if err != nil {
			continue
		}
		var itf intfInfo
		itf.Index = v.Index
		itf.IsEthernet = isEthernet(v.Index, aas)
		itf.Name = v.Name
		itf.MacAddress = v.HardwareAddr.String()
		for _, ip := range ips {
			if strings.Contains(ip.String(), ":") {
				itf.Ipv6 = append(itf.Ipv6, ip.String())
			} else {
				itf.Ipv4 = append(itf.Ipv4, ip.String())
			}
		}
		itfs = append(itfs, itf)
	}
	return itfs
}

// 根据网卡接口 Index 判断其是否为 Ethernet 网卡
func isEthernet(ifindex int, aas []*windows.IpAdapterAddresses) bool {
	result := false
	for _, aa := range aas {
		index := aa.IfIndex
		if ifindex == int(index) {
			switch aa.IfType {
			case windows.IF_TYPE_ETHERNET_CSMACD:
				result = true
			}

			if result {
				break
			}
		}
	}
	return result
}

// 从 net/interface_windows.go 中复制过来
func adapterAddresses() ([]*windows.IpAdapterAddresses, error) {
	var b []byte
	l := uint32(15000) // recommended initial size
	for {
		b = make([]byte, l)
		err := windows.GetAdaptersAddresses(syscall.AF_UNSPEC, windows.GAA_FLAG_INCLUDE_PREFIX, 0, (*windows.IpAdapterAddresses)(unsafe.Pointer(&b[0])), &l)
		if err == nil {
			if l == 0 {
				return nil, nil
			}
			break
		}

		if err.(syscall.Errno) != syscall.ERROR_BUFFER_OVERFLOW {
			return nil, os.NewSyscallError("getadaptersaddresses", err)
		}

		if l <= uint32(len(b)) {
			return nil, os.NewSyscallError("getadaptersaddresses", err)
		}
	}

	var aas []*windows.IpAdapterAddresses
	for aa := (*windows.IpAdapterAddresses)(unsafe.Pointer(&b[0])); aa != nil; aa = aa.Next {
		aas = append(aas, aa)
	}

	return aas, nil
}

// 主板信息
func GetMotherboardInfo() string {
	var s = []struct {
		Product string
	}{}

	err := wmi.Query("SELECT  Product  FROM Win32_BaseBoard", &s) // WHERE (Product IS NOT NULL)
	if err != nil {
		return ""
	}
	str := ""
	if len(s) > 0 {
		str = s[0].Product
	}
	return str
}

// BIOS信息
func GetBiosInfo() string {
	var s = []struct {
		Name string
	}{}
	err := wmi.Query("SELECT Name FROM Win32_BIOS WHERE (Name IS NOT NULL)", &s) // WHERE (BIOSVersion IS NOT NULL)
	if err != nil {
		return ""
	}
	str := ""
	if len(s) > 0 {
		str = s[0].Name
	}
	return str
}
