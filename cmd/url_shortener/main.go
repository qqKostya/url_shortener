package main

import (
	"fmt"
	"url_shortener/internal/config"
)

func main() {
	//TODO: init config: cleanenv +
	cfg := config.MustLoad()
	fmt.Println(cfg)
	//TODO: init logger: slog
	//TODO: init storege: sqlite
	//TODO: init router: chi, "chi render"
	//TODO: start server
}
