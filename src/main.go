/*
  MiningHQ Miner Manager - The MiningHQ Miner Manager GUI
  https://mininghq.io

  Copyright (C) 2018  Donovan Solms     <https://github.com/donovansolms>

  This program is free software: you can redistribute it and/or modify
  it under the terms of the GNU General Public License as published by
  the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.

  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU General Public License for more details.

  You should have received a copy of the GNU General Public License
  along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

// Package main implements the main runnable for the miner manager.
// It constructs and launches the Electron front-end
// package main
//
// import (
// 	"flag"
// 	"fmt"
// )
//
// // AppName is injected by the Astilectron packager
// var AppName string
//
// // BuiltAt is injected by the Astilectron packager
// var BuiltAt string
//
// // main implements the main runnable of the application
// func main() {
// 	// Grab the command-line flags
// 	debug := flag.Bool("d", false, "Enable debug mode")
// 	flag.Parse()
//
// 	if *debug {
// 		fmt.Println("RUNNING DEBUG")
// 	}
// 	fmt.Println("Running")
// }

// HACK Testing code below
//
package main

import (
	"github.com/donovansolms/mininghq-miner-manager/src/cmd"
)

// AppName is injected by the Astilectron bundler
var AppName string

func main() {
	cmd.AppName = AppName
	cmd.Asset = Asset
	cmd.RestoreAssets = RestoreAssets
	cmd.Execute()
}

//
// // Constants
// const htmlAbout = `Welcome on <b>Astilectron</b> demo!<br>
// This is using the bootstrap and the bundler. V0.0.2`
//
// // Vars
// var (
// 	AppName string
// 	BuiltAt string
// 	debug   = flag.Bool("d", false, "enables the debug mode")
// 	cli     = flag.Bool("c", false, "enables the cli mode")
// 	w       *astilectron.Window
// )
//
// func main() {
// 	// Init
// 	flag.Parse()
// 	astilog.FlagInit()
//
// 	if *cli == true {
// 		fmt.Println("We are running in cli mode")
// 	} else {
//
// 		// Run bootstrap
// 		astilog.Debugf("Running app built at %s", BuiltAt)
// 		if err := bootstrap.Run(bootstrap.Options{
// 			Asset:    Asset,
// 			AssetDir: AssetDir,
// 			AstilectronOptions: astilectron.Options{
// 				AppName:            AppName,
// 				AppIconDarwinPath:  "resources/icon.icns",
// 				AppIconDefaultPath: "resources/icon.png",
// 			},
// 			Debug: *debug,
// 			MenuOptions: []*astilectron.MenuItemOptions{{
// 				Label: astilectron.PtrStr("File"),
// 				SubMenu: []*astilectron.MenuItemOptions{
// 					{
// 						Label: astilectron.PtrStr("About"),
// 						OnClick: func(e astilectron.Event) (deleteListener bool) {
// 							if err := bootstrap.SendMessage(w, "about", htmlAbout, func(m *bootstrap.MessageIn) {
// 								// Unmarshal payload
// 								var s string
// 								if err := json.Unmarshal(m.Payload, &s); err != nil {
// 									astilog.Error(errors.Wrap(err, "unmarshaling payload failed"))
// 									return
// 								}
// 								astilog.Infof("About modal has been displayed and payload is %s!", s)
// 							}); err != nil {
// 								astilog.Error(errors.Wrap(err, "sending about event failed"))
// 							}
// 							return
// 						},
// 					},
// 					{Role: astilectron.MenuItemRoleClose},
// 				},
// 			}},
// 			OnWait: func(_ *astilectron.Astilectron, ws []*astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
// 				w = ws[0]
// 				go func() {
// 					time.Sleep(5 * time.Second)
// 					if err := bootstrap.SendMessage(w, "check.out.menu", "Don't forget to check out the menu!"); err != nil {
// 						astilog.Error(errors.Wrap(err, "sending check.out.menu event failed"))
// 					}
// 				}()
// 				return nil
// 			},
// 			RestoreAssets: RestoreAssets,
// 			Windows: []*bootstrap.Window{{
// 				Homepage:       "index.html",
// 				MessageHandler: handleMessages,
// 				Options: &astilectron.WindowOptions{
// 					BackgroundColor: astilectron.PtrStr("#333"),
// 					Center:          astilectron.PtrBool(true),
// 					Height:          astilectron.PtrInt(700),
// 					Width:           astilectron.PtrInt(700),
// 				},
// 			}},
// 		}); err != nil {
// 			astilog.Fatal(errors.Wrap(err, "running bootstrap failed"))
// 		}
// 	}
// }
//
// // handleMessages handles messages
// func handleMessages(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
// 	switch m.Name {
// 	case "explore":
// 		// Unmarshal payload
// 		var path string
// 		if len(m.Payload) > 0 {
// 			// Unmarshal payload
// 			if err = json.Unmarshal(m.Payload, &path); err != nil {
// 				payload = err.Error()
// 				return
// 			}
// 		}
//
// 		// Explore
// 		if payload, err = explore(path); err != nil {
// 			payload = err.Error()
// 			return
// 		}
// 	}
// 	return
// }
//
// // Exploration represents the results of an exploration
// type Exploration struct {
// 	Dirs       []Dir              `json:"dirs"`
// 	Files      *astichartjs.Chart `json:"files,omitempty"`
// 	FilesCount int                `json:"files_count"`
// 	FilesSize  string             `json:"files_size"`
// 	Path       string             `json:"path"`
// }
//
// // PayloadDir represents a dir payload
// type Dir struct {
// 	Name string `json:"name"`
// 	Path string `json:"path"`
// }
//
// // explore explores a path.
// // If path is empty, it explores the user's home directory
// func explore(path string) (e Exploration, err error) {
// 	// If no path is provided, use the user's home dir
// 	if len(path) == 0 {
// 		var u *user.User
// 		if u, err = user.Current(); err != nil {
// 			return
// 		}
// 		path = u.HomeDir
// 	}
//
// 	// Read dir
// 	var files []os.FileInfo
// 	if files, err = ioutil.ReadDir(path); err != nil {
// 		return
// 	}
//
// 	// Init exploration
// 	e = Exploration{
// 		Dirs: []Dir{},
// 		Path: path,
// 	}
//
// 	// Add previous dir
// 	if filepath.Dir(path) != path {
// 		e.Dirs = append(e.Dirs, Dir{
// 			Name: "..",
// 			Path: filepath.Dir(path),
// 		})
// 	}
//
// 	// Loop through files
// 	var sizes []int
// 	var sizesMap = make(map[int][]string)
// 	var filesSize int64
// 	for _, f := range files {
// 		if f.IsDir() {
// 			e.Dirs = append(e.Dirs, Dir{
// 				Name: f.Name(),
// 				Path: filepath.Join(path, f.Name()),
// 			})
// 		} else {
// 			var s = int(f.Size())
// 			sizes = append(sizes, s)
// 			sizesMap[s] = append(sizesMap[s], f.Name())
// 			e.FilesCount++
// 			filesSize += f.Size()
// 		}
// 	}
//
// 	// Prepare files size
// 	if filesSize < 1e3 {
// 		e.FilesSize = strconv.Itoa(int(filesSize)) + "b"
// 	} else if filesSize < 1e6 {
// 		e.FilesSize = strconv.FormatFloat(float64(filesSize)/float64(1024), 'f', 0, 64) + "kb"
// 	} else if filesSize < 1e9 {
// 		e.FilesSize = strconv.FormatFloat(float64(filesSize)/float64(1024*1024), 'f', 0, 64) + "Mb"
// 	} else {
// 		e.FilesSize = strconv.FormatFloat(float64(filesSize)/float64(1024*1024*1024), 'f', 0, 64) + "Gb"
// 	}
//
// 	// Prepare files chart
// 	sort.Ints(sizes)
// 	if len(sizes) > 0 {
// 		e.Files = &astichartjs.Chart{
// 			Data: &astichartjs.Data{Datasets: []astichartjs.Dataset{{
// 				BackgroundColor: []string{
// 					astichartjs.ChartBackgroundColorYellow,
// 					astichartjs.ChartBackgroundColorGreen,
// 					astichartjs.ChartBackgroundColorRed,
// 					astichartjs.ChartBackgroundColorBlue,
// 					astichartjs.ChartBackgroundColorPurple,
// 				},
// 				BorderColor: []string{
// 					astichartjs.ChartBorderColorYellow,
// 					astichartjs.ChartBorderColorGreen,
// 					astichartjs.ChartBorderColorRed,
// 					astichartjs.ChartBorderColorBlue,
// 					astichartjs.ChartBorderColorPurple,
// 				},
// 			}}},
// 			Type: astichartjs.ChartTypePie,
// 		}
// 		var sizeOther int
// 		for i := len(sizes) - 1; i >= 0; i-- {
// 			for _, l := range sizesMap[sizes[i]] {
// 				if len(e.Files.Data.Labels) < 4 {
// 					e.Files.Data.Datasets[0].Data = append(e.Files.Data.Datasets[0].Data, sizes[i])
// 					e.Files.Data.Labels = append(e.Files.Data.Labels, l)
// 				} else {
// 					sizeOther += sizes[i]
// 				}
// 			}
// 		}
// 		if sizeOther > 0 {
// 			e.Files.Data.Datasets[0].Data = append(e.Files.Data.Datasets[0].Data, sizeOther)
// 			e.Files.Data.Labels = append(e.Files.Data.Labels, "other")
// 		}
// 	}
// 	return
// }
