// tdxCrawl
package tdxcrawl

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

// 大显示器1920 x 1080, 分时分笔tic小窗口鼠标点击点: 分笔, 操作, 明细导出, txt格式, 导出保存;  行情(右上角)
var (
	行情 = mEpos{1717, 10} // 主窗口: 行情(右上角)

	分笔       = mEpos{999, 615}  // 分时分笔-小窗口-1: 分笔
	操作       = mEpos{915, 615}  // 分时分笔-小窗口-1: 操作
	关闭分时分笔窗口 = mEpos{1170, 170} // 分时分笔-小窗口-1:

	明细数据导出 = mEpos{915, 815} // 操作-小窗口-2: 明细数据导出

	格式文本文件   = mEpos{666, 300}  // 数据导出-小窗口-3: 格式TXT文本
	导出保存     = mEpos{985, 505}  // 数据导出-小窗口-3: 保存导出
	关闭数据导出窗口 = mEpos{1080, 270} // 数据导出-小窗口-3:

	关闭导出取消查看 = mEpos{1018, 318} // 确认-小窗口-4: 取消查看保存的文件

	日线  = mEpos{490, 33}   // 日线主窗口: 日线(顶,中)
	今日线 = mEpos{1630, 100} // 日线主窗口: 右边线, 避免没有退出鼠标十字线导致的cFpreTic()
)

func cBsleep(i int) {
	time.Sleep(time.Millisecond * time.Duration(i))
}

// 检查saveDir是否存在, 没有就创建; 并读取已经存在的文件, 返回存储名, 已有文件名;
func cCsaveDir(shareCode, saveDir string) (string, []string, error) {
	//sharecode
	saveDir = fmt.Sprintf("%s\\tics_%s", saveDir, shareCode)
	os.MkdirAll(saveDir, 0777)
	fs, err := ioutil.ReadDir(saveDir)
	if err != nil {
		return "", nil, err
	}
	exist := make([]string, len(fs))
	for i, f := range fs {
		exist[i] = f.Name()
	}
	return saveDir, exist, err
}

// 切换到某只shareCode的日线主窗口;
func cDtickWindow(shareCode string) {
	cEclickAt(行情) // 回到主窗口: 点击行情(右上角)
	for _, r := range shareCode {
		cBsleep(100)
		kDtype(uint16(r))
	}
	cBsleep(200)
	kDtype(VK_RETURN)
	cEclickAt(日线)
	cEclickAt(今日线) // 日线主窗口, 右边线: 避免没有退出鼠标十字线导致的cFpreTic()
	cBsleep(200)
	kDtype(VK_RIGHT)
}

// for cF,cH... 在屏幕某点, 点击鼠标左键
func cEclickAt(p mEpos) {
	cBsleep(200)
	mDmoveTo(p.x, p.y) // 移到(x,y)
	cBsleep(200)
	mFclickLMR(VK_LBUTTON) // 点击左键
}

// for cH...上一天的tic界面; 返回 日期和shareCode, 组成的默认保存文件名
// 未关闭界面
func cFpreTic() string {
	cBsleep(50)
	kFtypeVk(VK_LEFT) // 前一天
	return cFticNow()
}

// 当天的Tic界面;  返回 日期和shareCode, 组成的默认保存文件名
func cFticNow() string {
	cBsleep(200)
	kFtypeVk(VK_RETURN) // 进入分时小窗口
	cEclickAt(分笔)       // 点击分笔
	cEclickAt(操作)       // 点击操作
	cEclickAt(明细数据导出)   // 点击明细数据导出
	cBsleep(200)
	kDtype(VK_CONTROL, VK_C) // 复制文件名
	cBsleep(200)
	fileName := cBpaste() // 如: C:\Users\catsalt\Downloads\20200420_000001.xls
	cBsleep(200)
	return fileName[len(fileName)-19:] // 20200420_000001.xls
}

// 下一天的Tic界面;  返回 日期和shareCode, 组成的默认保存文件名
func cFnextTic() string {
	cBsleep(50)
	kDtype(VK_RIGHT)
	return cFticNow()
}

// 关闭tic窗口(1,3,4)
func cFcloseTic() {
	cBsleep(100)
	kDtype(VK_ESCAPE)
	cBsleep(100)
	kDtype(VK_ESCAPE)
}

// for cH... 检查文件名(日期+Code) 与 date之差, 天数;
func cGsub(fileName, date string) (int, error) {
	f, err := time.Parse("20060102", fileName[:8])
	if err != nil {
		return 0, fmt.Errorf("cGsub(1):  %w;", err)
	}
	d, err := time.Parse("20060102", date)
	if err != nil {
		return 0, fmt.Errorf("cGsub(2):  %w;", err)
	}
	return int(f.Sub(d).Hours() / 24), nil
}

// 找到时间段的右端end, 节假日回退
func cHrightEnd(end string) {
	first := true //用于标识end一开始, 是否已经在最右端第二个;
	for {
		f := cFpreTic() // 第一个有可能没有数据,左移一次到第二个;
		l, _ := cGsub(f, end)
		cFcloseTic()
		switch {
		case l > 0:
			for i := 0; i < l/2; i++ {
				cBsleep(50)
				kDtype(VK_LEFT)
			}
			first = false // 发生第二次左移
		case l < 0:
			if first { // 是第一次左移, 推出由不合理的end时间,导致 f-end <0; 保持此状态即可;
				return
			} else {
				fmt.Println("11")
				for {
					f = cFnextTic()
					cFcloseTic()
					l, _ = cGsub(f, end)
					if l >= 0 {
						return
					}
				}
			}
		default:
			return
		}
	}
}

// 存储每天的分时tick数据, 如果存在略过. begin~end 之间
func cIrangeSave(begin, end, saveDir string, exist []string) []string {
	cHrightEnd(end)
	cBsleep(300)
	f := cFticNow()
	for {
		l, err := cGsub(f, begin)
		if err != nil {
			fmt.Println(err)
			break
		}
		if l < 0 {
			break
		}
		hasno := true
		for _, v := range exist { // 判断有没有在目标存储文件夹里面
			if v == f {
				hasno = false
				break
			}
		}
		if hasno {
			kEwrite(saveDir + "\\" + f) // 输入新的存储地址
			cEclickAt(格式文本文件)           // 点击格式txt
			cEclickAt(导出保存)             // 点击导出保存
			exist = append(exist, f)    // 添加到exist中
		}
		cFcloseTic()
		today := f
		f = cFpreTic()
		if f == today {
			break
		}
	}
	cFcloseTic()
	return exist
}

// 切换到tdx主窗口, 注意必须打开了tdx.exe;
func ZcTdxWindow() (pid uint32, hwnd uintptr, err error) {
	pid, err = wApidOf("TdxW.exe")
	if err != nil {
		return pid, hwnd, err
	}
	hwnd = wEmainRoot(wChwndsOf(pid))
	return pid, hwnd, err
}

// 参数输入要正确!
func ZcSaveTic(shareCode, saveDir, begin, endDate string) {
	if begin >= endDate {
		return
	}
	saveDir, outed, err := cCsaveDir(shareCode, saveDir)
	if err != nil {
		fmt.Println(err, outed)
		return
	}
	cDtickWindow(shareCode)
	cIrangeSave(begin, endDate, saveDir, outed)
}

