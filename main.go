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
	XMLName   xml.Name  `xml:"track"`
	Location  string    `xml:"location"`
	Extension Extension `xml:"extension"`
}

type Extension struct {
	XMLName     xml.Name `xml:"extension"`
	Application string   `xml:"application,attr"`
	Id          string   `xml:"vlc:id"`
	Option      []string `xml:"vlc:option"`
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
	recFlag := flag.Lookup("r").Value.(flag.Getter).Get().(bool)
	dirFlag := flag.Lookup("dir").Value.(flag.Getter).Get().(string)

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
		if recFlag && len(dirs) != 0 {
			// recursively creating playlists for subfolders
			for _, dir := range dirs {
				createPlaylistFiles(dirFlag + "\\" + dir)
			}
		}
	}
	return videos
}

func createPlaylistXML(fileNames []string) string {
	var allTracks []Track

	audioFlag := flag.Lookup("audioTrack").Value.(flag.Getter).Get().(int)
	subFlag := flag.Lookup("subTrack").Value.(flag.Getter).Get().(int)
	subFileFlag := flag.Lookup("subFile").Value.(flag.Getter).Get().(string)
	noSubFlag := flag.Lookup("noSub").Value.(flag.Getter).Get().(bool)
	var subValue int
	if noSubFlag {
		// Setting sub track to fake number to turn off subtitles
		subValue = 99
	} else {
		subValue = subFlag
	}

	for i, fileName := range fileNames {
		option := []string{"audio-track=" + fmt.Sprint(audioFlag)}
		if subFileFlag != "" {
			trackNumber := fmt.Sprintf("%02d", i+1)
			option = append(option, "sub-file="+strings.Replace(subFileFlag, "$", trackNumber, 1))
		} else {
			option = append(option, "sub-track="+fmt.Sprint(subValue))
		}
		extension := &Extension{
			Id:          fmt.Sprint(i),
			Option:      option,
			Application: "http://www.videolan.org/vlc/playlist/0",
		}
		allTracks = append(allTracks, Track{Location: fileName, Extension: *extension})
	}

	tracklist := &TrackList{Track: allTracks}
	playlist := &Playlist{Tracklist: *tracklist}

	var out, _ = xml.MarshalIndent(playlist, "", "  ")
	return xml.Header + string(out)
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
	flag.Int("audioTrack", 0, "a number")
	flag.Int("subTrack", 0, "a number")
	flag.Bool("noSub", false, "a bool")
	flag.String("subFile", "", "a string")
	flag.Bool("r", false, "a boolean")

	flag.Parse()

	fmt.Println("tail:", flag.Args())
	if *dirPtr != "" {
		createPlaylistFiles(*dirPtr)
	}
}
