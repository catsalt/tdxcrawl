// Clipboard.go
package tdxcrawl

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"
)

// 拷贝类容到剪贴板
func cAcopy(str string) (clipBuffer uintptr) {
	if n, _, _ := wOpenClipboard.Call(0); n != 0 {
		wEmptyClipboard.Call()
		capBuffer := len(syscall.StringToUTF16(str)) * 2
		clipBuffer, _, _ = wGlobalAlloc.Call(0x0040, uintptr(capBuffer))
		lockBuffer, _, _ := wGlobalLock.Call(clipBuffer)
		wLstrcpy.Call(lockBuffer, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(str))))
		wSetClipboardData.Call(0x000D, clipBuffer)
		wGlobalUnlock.Call(clipBuffer)
		wCloseClipboard.Call()
	} else {
		fmt.Println("cAcopy() failed!")
	}
	return clipBuffer
}

// 从剪贴板的内容
func cBpaste() (str string) {
	for i := 0; i < 5; i++ {
		if n, _, _ := wOpenClipboard.Call(0); n != 0 {
			defer wCloseClipboard.Call()
			r, _, _ := wGetClipboardData.Call(0x000D)
			l, _, _ := wGlobalLock.Call(r)
			str = syscall.UTF16ToString((*[1 << 20]uint16)(unsafe.Pointer(l))[:]) //为啥1<<20
			wGlobalUnlock.Call(r)
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
	return str
}
func cDfree(clipBuffer uintptr) {
	wGlobalFree.Call(clipBuffer) //为啥不能释放??
}
