package tdxcrawl

import (
	// "syscall"
	"unsafe"
)

var (
	// wUser32   = syscall.NewLazyDLL("user32.dll")
	// wKernel32 = syscall.NewLazyDLL("kernel32.dll")

	wGetCursorPos = wUser32.NewProc("GetCursorPos")
	wSetCursorPos = wUser32.NewProc("SetCursorPos")

	// 参数: nums of INPUT,INPUT[0], size of INPUT[0]
	wSendInput = wUser32.NewProc("SendInput")
	// MapVirtualKey函数原型：UINT WINAPI MapVirtualKey ( _In_ UINT uCode,	 _In_ UINT uMapType)
	// MapVirtualKey函数参数详解：
	// uCode：定义一个键的扫描码或虚拟键码。该值如何解释依赖于uMapType参数的值。
	// uMapType：定义将要执行的翻译。该参数的值依赖于uCode参数的值。取值如下：
	// MAPVK_VK_TO_VSC 0：代表uCode是虚拟键码->扫描码。若一虚拟键码不区分左右，则返回左键的扫描码。若未进行翻译，则函数返回O。
	// MAPVK_VSC_TO_VK 1：代表uCode是扫描码->虚拟键码，且此虚拟键码不区分左右。若未进行翻译，则函数返回0。
	// MAPVK_VK_TO_CHAR 2：代表uCode为虚拟键->未被移位的字符值存放于返回值的低序字中。死键（发音符号）则通过设置返回值的最高位来表示。若未进行翻译，则函数返回0。
	// MAPVK_VSC_TO_VK_EX 3：代表uCode为扫描码->区分左右键的一虚拟键码。若未进行翻译，则函数返回0。
	// 返回值：返回值可以是一扫描码，或一虚拟键码，或一字符值，这完全依赖于不同的uCode和uMapType的值。若未进行翻译，则函数返回0。
	wMapVirtualKey = wUser32.NewProc("MapVirtualKeyW")
	// char->vk
	// wVkKeyScanExW = wUser32.NewProc("VkKeyScanExW")

	//复制,粘贴
	wEmptyClipboard   = wUser32.NewProc("EmptyClipboard")
	wOpenClipboard    = wUser32.NewProc("OpenClipboard")
	wSetClipboardData = wUser32.NewProc("SetClipboardData")
	wCloseClipboard   = wUser32.NewProc("CloseClipboard")
	wGetClipboardData = wUser32.NewProc("GetClipboardData")

	// GHND 0x0042 结合GMEM_MOVEABLE和GMEM_ZEROINIT。
	// GMEM_FIXED 0x0000 分配固定内存。返回值是指针。
	// GMEM_MOVEABLE 0x0002 分配可移动内存。内存块永远不会在物理内存中移动，但它们可以在默认堆中移动。
	// GMEM_ZEROINIT 0x0040 将内存内容初始化为零。
	// GPTR 0x0040 结合GMEM_FIXED和GMEM_ZEROINIT。
	// 可移动内存标志GHND和GMEM_MOVABLE添加不必要的开销，需要锁定才能安全使用。应避免使用，除非文件明确规定应使用它们。
	wGlobalAlloc  = wKernel32.NewProc("GlobalAlloc")
	wGlobalFree   = wKernel32.NewProc("GlobalFree")
	wGlobalLock   = wKernel32.NewProc("GlobalLock")
	wGlobalUnlock = wKernel32.NewProc("GlobalUnlock")
	// wStringCchCopyW = wKernel32.NewProc("StringCchCopyW")
	wLstrcpy = wKernel32.NewProc("lstrcpyW")
)

func wAsendInput(inputs ...mAinput) uint32 {
	if len(inputs) == 0 {
		return 0
	}
	ret, _, _ := wSendInput.Call(
		uintptr(len(inputs)),
		uintptr(unsafe.Pointer(&inputs[0])),
		unsafe.Sizeof(inputs[0]),
	)
	return uint32(ret)
}
