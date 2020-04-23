package tdxcrawl

import (
	"fmt"
	"unsafe"
	// "unsafe"
)

// mAinput Type 参数值:
// INPUT_MOUSE 0,mouse, mi structure
// INPUT_KEYBOARD 1,keyboard, ki structure.
// INPUT_HARDWARE 2,hardware, hi structure.
type mAinput struct {
	Type uint32
	// use mAinput for the union because it is the largest of all allowed structures
	mouse mBinput
}

// type HARDWAREINPUT struct {
// 	Msg    uint32
// 	ParamL uint16
// 	ParamH uint16
// }
// func HardwareInput(input HARDWAREINPUT) mAinput {
// 	return mAinput{
// 		Type:  2,
// 		mouse: *((*mBinput)(unsafe.Pointer(&input))),
// 	}
// }

// mouseData:
// 如果dwFlags包含MOUSEEVENTF_WHEEL，则鼠标数据指定车轮移动量。正值向前旋转，远离焦点,负值向后旋转。一个车轮单击定义为WHEEL_DELTA, 120。
// 如果dwFlags不包含MOUSEEVENTF_WHEEL、MOUSEEVENTF_XDOWN或MOUSEEVENTF_XUP，则鼠标数据应为零。
// 如果dwFlags包含MOUSEEVENTF_XDOWN或MOUSEEVENTF_XUP，则指定按下或释放的 X 按钮。此值可以是以下标志的任意组合。
// 0x0001	XBUTTON1  设置是否按下或释放第一个 X 按钮。
// 0x0002	XBUTTON2  设置是否按下或释放第二个 X 按钮。
// Flags:
// 一组位标志，用于指定鼠标运动和按钮单击的各个方面。此成员中的位可以是以下值的任何合理组合。
// 指定鼠标按钮状态的位标志设置为指示状态更改，而不是正在进行的条件。例如，如果按下并按住鼠标左键，则在首次按下左键时设置MOUSEEVENTF_LEFTDOWN，但不针对后续动作设置。同样，只有在首次发布按钮时才设置MOUSEEVENTF_LEFTUP。
// 不能同时在dwFlags参数中指定MOUSEEVENTF_WHEEL标志和MOUSEEVENTF_XDOWN或MOUSEEVENTF_XUP标志，因为它们都需要使用鼠标数据字段。
// 0x0001 MOUSEEVENTF_MOVE 发生了移动。
// 0x2000 MOUSEEVENTF_MOVE_NOCOALESCE 不会合并WM_MOUSEMOVE消息。默认行为是合并WM_MOUSEMOVE消息。 视窗XP/2000：不支持此值。
// 0x0002 MOUSEEVENTF_LEFTDOWN 按下了左侧按钮。
// 0x0004 MOUSEEVENTF_LEFTUP 左键已释放。
// 0x0008 MOUSEEVENTF_RIGHTDOWN 按下了右侧按钮。
// 0x0010 MOUSEEVENTF_RIGHTUP 右键已释放。
// 0x0020 MOUSEEVENTF_MIDDLEDOWN 按下中间按钮。
// 0x0040 MOUSEEVENTF_MIDDLEUP 中间按钮已释放。
// 0x8000 MOUSEEVENTF_ABSOLUTE dx和dy成员包含规范化的绝对坐标。如果未设置标志，则 dx和dy包含相对数据（自上次报告位置以来位置的变化）。可以设置或不设置此标志，而不管连接到系统的鼠标或其他指针设备（如果有）。有关相对鼠标运动的详细信息，请参阅以下备注部分。
// 0x4000 MOUSEEVENTF_VIRTUALDESK 将坐标映射到整个桌面。必须与MOUSEEVENTF_ABSOLUTE一起使用。
// 0x1000 MOUSEEVENTF_HWHEEL 如果鼠标有滚轮，则车轮是水平移动的。移动量在鼠标数据中指定。 视窗XP/2000：不支持此值。
// 0x0800 MOUSEEVENTF_WHEEL 如果鼠标有轮子，则车轮被移动。移动量在鼠标数据中指定。
// 0x0080 MOUSEEVENTF_XDOWN 按下 X 按钮。
// 0x0100 MOUSEEVENTF_XUP 释放了 X 按钮。
type mBinput struct {
	Dx        int32  //鼠标位置 x轴
	Dy        int32  //鼠标位置 y轴
	MouseData uint32 //见前面注释
	Flags     uint32 //见前面注释
	Time      uint32 //如果为零,系统提供时间戳
	ExtraInfo uintptr
}

// 生成鼠标类型;
func mCinput(input mBinput) mAinput {
	return mAinput{
		Type:  0,
		mouse: input,
	}
}

// 鼠标移动到(x,y)点处;左上角为坐标原点, 正数
func mDmoveTo(x, y int32) error {
	if n, _, _ := wSetCursorPos.Call(uintptr(x), uintptr(y)); n == 0 {
		return fmt.Errorf("mDmoveTo(): %d, %d;", x, y)
	}
	return nil
}

// 鼠标以现在位置为原点,沿x,y轴移动x,y距离, 可以为负数;
func mDmove(x, y int32) (err error) {
	var p mEpos
	if n, _, _ := wGetCursorPos.Call(uintptr(unsafe.Pointer(&p))); n == 0 {
		return fmt.Errorf("mDmove(1): ;")
	}
	if err = mDmoveTo(p.x+x, p.y+y); err != nil {
		err = fmt.Errorf("mDmove(2): %w;", err)
	}
	return err
}

type mEpos struct {
	x, y int32
}

// 获取当前鼠标位置;
func mEgetPos() (pos mEpos, err error) {
	if n, _, _ := wGetCursorPos.Call(uintptr(unsafe.Pointer(&pos))); n == 0 {
		err = fmt.Errorf("mEgetPos(): ;")
	}
	return pos, err
}

// 鼠标点击(点按,然后释放); 左键 VK_LBUTTON, 右键VK_RBUTTON, 中键VK_MBUTTON;
func mFclickLMR(vk uint16) (err error) {
	switch vk {
	case VK_LBUTTON:
		err = mHbutton(2, 4) //不需要2+4
	case VK_RBUTTON:
		err = mHbutton(8, 18) //8+10, 否则成右键选择
	case VK_MBUTTON:
		err = mHbutton(20, 40) //20+40
	default:
		err = fmt.Errorf("mFclickLMR(1): %v;", vk)
	}
	if err != nil {
		err = fmt.Errorf("mFclickLMR(2): %v, %w;", vk, err)
	}
	return err
}

// 鼠标点按, 不释放; 左键 VK_LBUTTON, 右键VK_RBUTTON, 中键VK_MBUTTON;
func mGpressHold(vk uint16) (err error) {
	switch vk {
	case VK_LBUTTON:
		err = mHbutton(2)
	case VK_RBUTTON:
		err = mHbutton(8)
	case VK_MBUTTON:
		err = mHbutton(20)
	default:
		err = fmt.Errorf("mGpressHold(1): %v;", vk)
	}
	if err != nil {
		err = fmt.Errorf("mGpressHold(2): %v, %w;", vk, err)
	}
	return err
}

// 鼠标释放; 左键 VK_LBUTTON, 右键VK_RBUTTON, 中键VK_MBUTTON;
func mGreleaseVk(vk uint16) (err error) {
	switch vk {
	case VK_LBUTTON:
		err = mHbutton(6)
	case VK_RBUTTON:
		err = mHbutton(18)
	case VK_MBUTTON:
		err = mHbutton(60)
	default:
		err = fmt.Errorf("mGreleaseVK(1): %v,", vk)
	}
	if err != nil {
		err = fmt.Errorf("mGreleaseVK(2): %v, %w;", vk, err)
	}
	return err
}

// 鼠标状态, down or up, 最多两个参数; 参看前面mBinput参数;
func mHbutton(downUp ...uint32) (err error) {
	if len(downUp) < 1 || len(downUp) > 2 {
		return fmt.Errorf("mHbutton(1): %v;", downUp)
	}
	if wAsendInput(mCinput(mBinput{Flags: downUp[0]})) == 0 {
		return fmt.Errorf("mHbutton(2): %v;", downUp)
	}
	if len(downUp) == 2 {
		if wAsendInput(mCinput(mBinput{Flags: downUp[1]})) == 0 {
			return fmt.Errorf("mHbutton(3): %v;", downUp)
		}
	}
	return nil
}
