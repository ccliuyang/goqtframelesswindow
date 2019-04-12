package qframelesswindow

import (
	"unsafe"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"

	win "github.com/akiyosi/w32"
)

func (f *QFramelessWindow) SetNativeEvent(app *widgets.QApplication) {
	filterObj := core.NewQAbstractNativeEventFilter()
	filterObj.ConnectNativeEventFilter(func(eventType *core.QByteArray, message unsafe.Pointer, result int) bool {
		msg := (*win.MSG)(message)
		lparam := msg.LParam
		hwnd := msg.Hwnd

		switch msg.Message {
		case win.WM_CREATE:
		 	style := win.GetWindowLong(hwnd, win.GWL_STYLE)
		 	style = style | win.WS_THICKFRAME | win.WS_CAPTION
		 	win.SetWindowLong(hwnd, win.GWL_STYLE, uint32(style))

			// shadow
			shadow := &win.MARGINS{0, 0, 0, 1}
			win.DwmExtendFrameIntoClientArea(hwnd, shadow)
			var nullptr win.HWND
			var uflag uint
			uflag |= win.SWP_NOZORDER | win.SWP_NOOWNERZORDER | win.SWP_NOMOVE | win.SWP_NOSIZE | win.SWP_FRAMECHANGED
			win.SetWindowPos(hwnd, nullptr, 0, 0, 0, 0, uflag)

			f.borderless = true
		 	return false

		case win.WM_NCCALCSIZE:
			if msg.WParam == 1 && f.borderless {
				// this kills the window frame and title bar we added with WS_THICKFRAME and WS_CAPTION
				result = 0
				if win.IsWindowVisible(hwnd) == false {
		 			win.ShowWindow(hwnd, win.SW_SHOW)
				}
			}
			return true

		case win.WM_GETMINMAXINFO:
			mm := (*win.MINMAXINFO)((unsafe.Pointer)(lparam))
			mm.PtMinTrackSize.X = int32(f.minimumWidth)
			mm.PtMinTrackSize.Y = int32(f.minimumHeight)
			return true

		default:
		}
		return false
	})
	app.InstallNativeEventFilter(filterObj)
}
