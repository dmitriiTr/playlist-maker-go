package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Playlist struct {
	XMLName   xml.Name  `xml:"playlist"`
	Tracklist TrackList `xml:"trackList"`
}

type TrackList struct {
	XMLName xml.Name `xml:"trackList"`
	Track   []Track  `xml:"track"`
}

type Track struct {
	XMLName  xml.Name `xml:"track"`
	Location string   `xml:"location"`
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func getVideosNames(pathToVideos string) []string {
	ex, _ := exists(pathToVideos)
	dirs := []string{}
	videos := []string{}
	rPtr := flag.Lookup("r").Value.(flag.Getter).Get().(bool)

	if !ex {
		fmt.Println("not found")
	} else {
		fmt.Println("found")
		entries, _ := os.ReadDir(pathToVideos)
		for _, v := range entries {
			name := v.Name()
			if v.IsDir() {
				dirs = append(dirs, name)
			} else {
				ext := filepath.Ext(name)
				if ext == ".mkv" || ext == ".mp4" {
					videos = append(videos, "file:///"+pathToVideos+"\\"+name)
				}
			}
		}
		if rPtr && len(dirs) != 0 {
			// recursively creating playlists for subfolders
			dirPtr := flag.Lookup("dir").Value.(flag.Getter).Get().(string)
			for _, dir := range dirs {
				createPlaylistFiles(dirPtr + "\\" + dir)
			}
		}
	}
	return videos
}

func createPlaylistXML(fileNames []string) string {
	var allTracks []Track

	for _, fileName := range fileNames {
		allTracks = append(allTracks, Track{Location: fileName})
	}

	tracklist := &TrackList{Track: allTracks}
	playlist := &Playlist{Tracklist: *tracklist}

	var out, _ = xml.MarshalIndent(playlist, " ", "  ")
	fmt.Println(xml.Header + string(out))
	return string(out)
}

func createPlaylistFiles(path string) {
	names := getVideosNames(path)
	if len(names) != 0 {
		file := createPlaylistXML(names)
		pathArray := strings.Split(path, "\\")
		fileName := pathArray[len(pathArray)-1]
		fileNameWithPath := path + "\\" + fileName + ".xspf"
		f, e := os.Create(fileNameWithPath)

		if e != nil {
			panic(e)
		}

		f.WriteString(file)
		defer f.Close()

		fmt.Println("File \"" + fileNameWithPath + "\" is created")
	}
}

func main() {
	dirPtr := flag.String("dir", "", "a string")
	audioPtr := flag.Int("audioTrack", 0, "a number")
	subPtr := flag.String("subFile", "", "a string")
	flag.Bool("r", false, "a boolean")

	flag.Parse()

	fmt.Println("dirPtr:", *dirPtr)
	fmt.Println("audioPtr:", *audioPtr)
	fmt.Println("subPtr:", *subPtr)

	fmt.Println("tail:", flag.Args())
	if *dirPtr != "" {
		createPlaylistFiles(*dirPtr)
	}
}
