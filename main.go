package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

// This tells Go: "Take the 'dist' folder and bake it into this binary"
//
//go:embed dist/*
var content embed.FS

func main() {
	distDir, err := fs.Sub(content, "dist")
	if err != nil {
		log.Fatal(err)
	}

	fileServer := http.FileServer(http.FS(distDir))

	http.Handle("/", fileServer)

	log.Println(`
⢕⢕⢕⢕⢕⢕⢕⢕⢕⢕⢕⢕⢕⢕⢕⢕⢕⢕⢕⢕⢕⢕⠕⠕⠕⠕⢕⢕
⢕⢕⢕⢕⢕⠕⠕⢕⢕⢕⢕⢕⢕⢕⢕⢕⢕⠕⠁⣁⣠⣤⣤⣤⣶⣦⡄⢑
⢕⢕⢕⠅⢁⣴⣤⠀⣀⠁⠑⠑⠁⢁⣀⣀⣀⣀⣘⢻⣿⣿⣿⣿⣿⡟⢁⢔
⢕⢕⠕⠀⣿⡁⠄⠀⣹⣿⣿⣿⡿⢋⣥⠤⠙⣿⣿⣿⣿⣿⡿⠿⡟⠀⢔⢕
⢕⠕⠁⣴⣦⣤⣴⣾⣿⣿⣿⣿⣇⠻⣇⠐⠀⣼⣿⣿⣿⣿⣿⣄⠀⠐⢕⢕
⠅⢀⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⣶⣶⣿⣿⣿⣿⣿⣿⣿⣿⣷⡄⠐⢕
⠅⢸⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡄⠐
⢄⠈⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡆
⢕⢔⠀⠈⠛⠿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿
⢕⢕⢄⠈⠳⣶⣶⣶⣤⣤⣤⣤⣭⡍⢭⡍⢨⣯⡛⢿⣿⣿⣿⣿⣿⣿⣿⣿
⢕⢕⢕⢕⠀⠈⠛⠿⢿⣿⣿⣿⣿⣿⣦⣤⣿⣿⣿⣦⣈⠛⢿⢿⣿⣿⣿⣿
⢕⢕⢕⠁⢠⣾⣶⣾⣭⣖⣛⣿⠿⣿⣿⣿⣿⣿⣿⣿⣿⣷⡆⢸⣿⣿⣿⡟
⢕⠅⢀⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠟⠈⢿⣿⣿⡇⡇
⢕⠕⠀⠼⠟⢉⣉⡙⠻⠿⢿⣿⣿⣿⣿⣿⡿⢿⣛⣭⡴⠶⠶⠂⠀⠿⠿⠇

Starting server
	`)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
