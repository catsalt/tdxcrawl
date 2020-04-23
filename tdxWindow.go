// tdxWindow.go
package tdxcrawl

import (
	"fmt"
	"sort"

	"syscall"
	"unsafe"
)

// 调用w32 api, 因为tdx有多个窗口, 以下函数经验证均无法找到tdx的pid,hwnd,handle
// GetCurrentProcessId()
// GetCurrentThread()
// GetDesktopWindow()
// GetForegroundWindow()
// GetActiveWindow()
// GetFocus()
// GetCurrentActCtx()

// 验证之后, 采用根据TdxW.exe名称, 遍历系统Pid,找到tdx的pid,
// 根据pid, 遍历所有的窗口, 找出tdx的窗口句柄hwnd(不止一个)
// 在所有的tdx窗口找出主窗口hwnd
var (
	wKernel32                 = syscall.NewLazyDLL("kernel32.dll")
	wCreateToolhelp32Snapshot = wKernel32.NewProc("CreateToolhelp32Snapshot")
	wProcess32Next            = wKernel32.NewProc("Process32Next")
	wCloseHandle              = wKernel32.NewProc("CloseHandle")

	wUser32                   = syscall.NewLazyDLL("user32.dll")
	wGetWindowThreadProcessId = wUser32.NewProc("GetWindowThreadProcessId")
	wEnumWindows              = wUser32.NewProc("EnumWindows")
	wGetAncestor              = wUser32.NewProc("GetAncestor")
	wShowWindow               = wUser32.NewProc("ShowWindow")
)

// 顺序不可变!
type wPROCESSENTRY32 struct {
	dwSize, cntUsage, th32ProcessID                                        int32
	th32DefaultHeapID                                                      uintptr
	th32ModuleID, cntThreads, th32ParentProcessID, pcPriClassBase, dwFlags int32
	szExeFile                                                              [260]byte
}

// 根据进程名遍历查找Pid
func wApidOf(pName string) (pid uint32, err error) {
	pHandle, _, _ := wCreateToolhelp32Snapshot.Call(uintptr(0x2), uintptr(0x0))
	if int(pHandle) != -1 {
		n := 0
		var proc wPROCESSENTRY32
		proc.dwSize = int32(unsafe.Sizeof(proc))
		for {
			if rt, _, _ := wProcess32Next.Call(pHandle, uintptr(unsafe.Pointer(&proc))); int(rt) == 1 {
				has := true
				for i := 0; i < len(pName); i++ {
					if pName[i] != proc.szExeFile[i] {
						has = false
						break
					}
				}
				if has {
					// fmt.Println(string(proc.szExeFile[:]))
					pid = uint32(proc.th32ProcessID)
					n++
				}
			} else {
				break
			}
		}
		if n != 1 {
			err = fmt.Errorf("wApidOf(): %v 没有对应的Pid;", pName)
		}
	}
	wCloseHandle.Call(pHandle)
	return pid, err
}

// for wC...窗口句柄Hwnd->进程句柄Handle,进程PID
func wBhandlePidOf(hwnd uintptr) (handle uintptr, pid uint32) {
	handle, _, _ = wGetWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&pid)))
	return handle, pid
}

// 进程PID->所有的窗口句柄Hwnds
func wChwndsOf(pid uint32) (hwnds []uintptr) {
	f := syscall.NewCallback(func(w, _ uintptr) uintptr {
		_, id := wBhandlePidOf(w)
		if pid == id {
			hwnds = append(hwnds, w)
		}
		return 1
	})
	wEnumWindows.Call(f, 0)
	sort.SliceStable(hwnds, func(i, j int) bool {
		if hwnds[i] > hwnds[j] {
			return true
		}
		return false
	})
	return hwnds
}

// for wE... 参数1表示父窗口(多次为65552), 2表示根窗口, 3表示根据1,2的父链追随根窗口, 根据观察设置 3
func wDrootOf(hwnd uintptr) uintptr {
	root, _, _ := wGetAncestor.Call(hwnd, 3)
	return root
}

// 此处m>10 是观察tdx实际窗口数确定的, 其他程序另行设置
func wEmainRoot(hwnds []uintptr) (hwnd uintptr) {
	m := 0
	for _, v := range hwnds {
		if wDrootOf(v) == hwnd {
			m++
		}
		hwnd = wDrootOf(v)
		if m > 10 {
			break
		}
	}
	wShowWindow.Call(hwnd, 2)
	wShowWindow.Call(hwnd, 3)
	return hwnd
}
