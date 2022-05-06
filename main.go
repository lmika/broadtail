package main

import (
	"context"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/lmika/broadtail/config"
	"github.com/lmika/broadtail/handlers"
)

func main() {
	flagDevMode := flag.Bool("dev", false, "dev mode")
	flagYTDLSimulator := flag.Bool("ytdl-simulator", false, "use youtube-dl simulator (for dev)")
	flagConfigFile := flag.String("config", "", "config file")
	flag.Parse()

	cfg, err := config.Read(*flagConfigFile)
	if err != nil {
		log.Printf("warn: cannot read config file: %v", *flagConfigFile)
	}

	var templateFS, assetsFS fs.FS
	if *flagDevMode {
		log.Println("Starting in dev mode")
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

	handler, closeFn, err := handlers.Server(handlers.Config{
		LibraryDir:          cfg.LibraryDir,
		LibraryOwner:        cfg.LibraryOwner,
		PlexBaseURL:         cfg.PlexBaseURL,
		PlexToken:           cfg.PlexToken,
		BaseDataDir:         cfg.DataDir,
		JobDataFile:         filepath.Join(cfg.DataDir, "jobs.db"),
		VideoDataFile:       filepath.Join(cfg.DataDir, "videos.db"),
		FeedsDataFile:       filepath.Join(cfg.DataDir, "feeds.db"),
		FeedItemsDataFile:   filepath.Join(cfg.DataDir, "feeditem.db"),
		FavouritesDataFile:  filepath.Join(cfg.DataDir, "favourites.db"),
		YTDownloadCommand:   cfg.YoutubeDLCommandAsSlice(),
		YTDownloadSimulator: *flagYTDLSimulator,
		TemplateFS:          templateFS,
		AssetFS:             assetsFS,
	})
	if err != nil {
		log.Fatalln(err)
	}

	server := &http.Server{
		Addr:    fmt.Sprintf("%v:%v", cfg.BindAddr, cfg.Port),
		Handler: handler,
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		<-c

		ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelFn()
		server.Shutdown(ctx)
	}()

	log.Printf("Listening on %v:%v", cfg.BindAddr, cfg.Port)
	server.ListenAndServe()

	log.Printf("Shutting down")
	closeFn()

	log.Printf("All done. Bye.")
}
