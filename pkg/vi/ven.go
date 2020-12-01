package vi

import (
	"github.com/bkthomps/Ven/screen"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"log"
	"os"
)

func RunVendorVi() {

	if len(os.Args) != 2 {
		print("Usage: ./cdk vi <file_name>\n")
		return
	}
	userArg := os.Args[1]
	tCellScreen, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}
	encoding.Register()
	quit := make(chan struct{})
	s := &screen.Screen{}
	s.Init(tCellScreen, quit, userArg)
	<-quit
	tCellScreen.Fini()
}
