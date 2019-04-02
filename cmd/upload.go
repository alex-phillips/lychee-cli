// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/alex-phillips/lychee/lib/log"
	"github.com/spf13/cobra"
)

// downloadThreadCmd represents the freeleech command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload images",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var files []string
		var ids []string
		for _, fpath := range args {
			filepathInfo, err := os.Stat(fpath)
			if err != nil {
				log.Error.Println(err)
				continue
			}

			switch mode := filepathInfo.Mode(); {
			case mode.IsDir():
				err := filepath.Walk(fpath, func(path string, info os.FileInfo, err error) error {
					if info.Mode().IsRegular() {
						files = append(files, path)
					}

					return nil
				})

				if err != nil {
					log.Error.Println(err)
				}

				break
			case mode.IsRegular():
				files = append(files, fpath)
				break
			}
		}

		parentID := "0"
		if cmd.Flag("parent").Value.String() != "" {
			parentID = api.GetAlbumByName(cmd.Flag("parent").Value.String()).GetID()
		}

		for _, file := range files {
			log.Info.Println("Uploading " + file)
			id, err := api.Upload(file, parentID)
			if err != nil {
				log.Error.Println("Error uploading " + file + ": " + err.Error())
				continue
			}

			ids = append(ids, id)
		}

		tags := cmd.Flag("tags").Value.String()
		if tags != "" {
			api.SetTags(ids, strings.Split(tags, ","))
		}
	},
}

func init() {
	uploadCmd.PersistentFlags().StringP("tags", "t", "", "Set tags of new photo(s)")
	uploadCmd.PersistentFlags().String("parent", "", "Put photos in specified album")
	rootCmd.AddCommand(uploadCmd)
}
