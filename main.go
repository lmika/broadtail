package main

import (
	"flag"
	"fmt"
	"github.com/lmika/broadtail/handlers"
	"io/fs"
	"log"
	"net/http"
	"os"
)

func main() {
	flagBindAddr := flag.String("bind", "", "bind address")
	flagDevMode := flag.Bool("dev", false, "dev mode")
	flagPort := flag.Int("p", 3690, "port")
	flagLibraryDir := flag.String("library", "", "library dir")
	flag.Parse()

	var templateFS, assetsFS fs.FS
	if *flagDevMode {
		templateFS = os.DirFS("templates")
		assetsFS = os.DirFS("build/assets")
	} else {
		var err error
		templateFS, err = fs.Sub(embeddedTemplates, "templates")
		if err != nil {
			log.Fatal("cannot get sub-fs of templates")
		}

		assetsFS, err = fs.Sub(embeddedAssets, "build/assets")
		if err != nil {
			log.Fatal("cannot get sub-fs of templates")
		}
	}

	server, err := handlers.Server(handlers.Config{
		LibraryDir: *flagLibraryDir,
		TemplateFS: templateFS,
		AssetFS: assetsFS,
	})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Listening on %v:%v", *flagBindAddr, *flagPort)
	log.Println(http.ListenAndServe(fmt.Sprintf("%v:%v", *flagBindAddr, *flagPort), server))
}
