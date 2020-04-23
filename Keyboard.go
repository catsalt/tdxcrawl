package tdxcrawl

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf16"
	"unsafe"
)

//必须用 mAinput, mouse, 否则出错, why? 可能是Win API 的源码设置如此.

// Flags 参数值: (可以组合)
// 0x0001 KEYEVENTF_EXTENDEDKEY 如果指定，扫描代码前面有一个前缀字节，该字节的值为 0xE0 （224）。
// 0x0002 KEYEVENTF_KEYUP 如果指定，则释放Key。如果未指定，则按下该键。
// 0x0008 KEYEVENTF_SCANCODE 如果指定，wScan标识key，并且Vk将被忽略。根据扫描代码定义键盘输入。键盘不同->虚拟键值不同->扫描码相同! 这对于模拟物理击键非常有用。
// 0x0004 KEYEVENTF_UNICODE 如果指定，系统将合成VK_PACKET击键。Vk参数必须为零。只能与KEYEVENTF_KEYUP标志组合。有关详细信息，请参阅备注部分。
type kBinput struct {
	Vk        uint16 //虚拟键值
	Scan      uint16 //扫描键值
	Flags     uint32 //见前面注释
	Time      uint32 //如果为零, 系统提供时间戳
	ExtraInfo uintptr
}

// INPUT_KEYBOARD对应 ki structure 取值 1; 必须用 mAinput, mouse, 否则出错, why?
// 生成键盘类型Struct;
func kCinput(input kBinput) mAinput {
	return mAinput{Type: 1, mouse: *((*mBinput)(unsafe.Pointer(&input)))}
}

// 组合按键(ctr,alt,shift), 顺序重要! 最多前3个有效.
// 一些系统功能, 如ctr+alt+del, win+l, 无法使用
func kDtype(vks ...uint16) error {
	if len(vks) > 3 {
		fmt.Println("kDcombind(1): 组合键过多取前3个, ", vks)
		vks = vks[:2]
	}
	inputs := make([]mAinput, len(vks)*2)
	for i, vk := range vks {
		inputs[i] = kCinput(kBinput{Vk: vk})
		inputs[len(inputs)-1-i] = kCinput(kBinput{Vk: vk, Flags: 2})
	}
	if wAsendInput(inputs...) == 0 {
		return fmt.Errorf("kDcombind(2): %v;", vks)
	}
	time.Sleep(time.Millisecond * 300)
	return nil
}

// 所见即所得
// 输入文本字符串(包括字母,符号,数字,\",\\,'不用转义)用Unicode, (\b,\t,\n,\v,\f,\r)用vk;
func kEwrite(str string) (err error) {
	s := strings.ReplaceAll(str, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	for i, u := range utf16.Encode([]rune(s)) {
		switch u {
		case 8:
			err = kFtypeVk(VK_BACK)
		case 9:
			err = kFtypeVk(VK_TAB)
		case 10:
			err = kFtypeVk(VK_RETURN)
		case 7, 11, 12: //未定义!!
		default:
			err = kFtypeUnicode(u)
		}
		if err != nil {
			return fmt.Errorf("kEwrite(): %s, %d, %w;", str, i, err)
			fmt.Println("kEwrite(): ", str, i, err)
		}
	}
	return nil
}

// 按下按键,然后释放, 直接输入 unicode字符, 非vk
func kFtypeUnicode(u uint16) error {
	if wAsendInput(
		kCinput(kBinput{Scan: u, Flags: 4}),
		kCinput(kBinput{Scan: u, Flags: 6}), //4+2,
	) == 0 {
		return fmt.Errorf("kEtypeUnicode(): %U;", u)
	}
	return nil
}

// 0x30-0x39(数字键0~9) 输入'0'~'9',得到0~9
// 0x41-0x51(字母键A~Z) 输入'A'~'Z',得到a~z, 因键盘是以大写字母表示
// 按下按键,随后释放,
func kFtypeVk(vk uint16) error {
	if wAsendInput(
		kCinput(kBinput{Vk: vk}),
		kCinput(kBinput{Vk: vk, Flags: 2}),
	) == 0 {
		return fmt.Errorf("kFtypeVk(): %v;", vk)
	}
	return nil
}

// 按下按键, 不释放, 第二次按这个键时, 无效
func kGpressHold(vk uint16) error {
	if wAsendInput(kCinput(kBinput{Vk: vk})) == 0 {
		return fmt.Errorf("kFpress(): %v;", vk)
	}
	return nil
}

// 释放按键
func kGrelease(vk uint16) error {
	if wAsendInput(kCinput(kBinput{Vk: vk, Flags: 2})) == 0 {
		return fmt.Errorf("kFrelease(): %v;", vk)
	}
	return nil
}
