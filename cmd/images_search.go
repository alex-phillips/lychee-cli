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
	"strings"

	"github.com/alex-phillips/lychee/lib/log"
	"github.com/alex-phillips/lychee/lib/lychee"
	"github.com/spf13/cobra"
)

// downloadThreadCmd represents the freeleech command
var imagesSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for images",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Error.Println("You must include a search query")
			os.Exit(1)
		}

		results := api.Search(args[0])

		var tagFilter []string
		if cmd.Flag("tags").Value.String() != "" {
			tagFilter = strings.Split(cmd.Flag("tags").Value.String(), ",")
		}

		var filtered []lychee.Photo
		for _, photo := range results.Photos {
			photoTags := strings.Split(photo.Tags, ",")

			include := true

			if len(tagFilter) > 0 {
				for _, filterTag := range tagFilter {
					tagFound := false
					for _, tag := range photoTags {
						if tag == filterTag {
							tagFound = true
							break
						}
					}

					if tagFound == false {
						include = false
						break
					}
				}
			}

			if include == true {
				filtered = append(filtered, photo)
			}
		}

		var ids []string
		for _, photo := range filtered {
			log.Info.Println(photo.Title)
			ids = append(ids, photo.GetID())
		}

		if cmd.Flag("set-album").Value.String() != "" {
			album := api.GetAlbumByName(cmd.Flag("set-album").Value.String())
			api.MoveImages(ids, album.GetID())
		}

		if cmd.Flag("add-tags").Value.String() != "" {
			addTags := strings.Split(cmd.Flag("add-tags").Value.String(), ",")
			for _, photo := range filtered {
				photoTags := strings.Split(photo.Tags, ",")

				for _, addTag := range addTags {
					add := true
					for _, photoTag := range photoTags {
						if addTag == photoTag {
							add = false
							break
						}
					}

					if add == true {
						photoTags = append(photoTags, addTag)
					}
				}

				api.SetTags([]string{photo.GetID()}, photoTags)
			}
		}

		if cmd.Flag("replace-tags").Value.String() != "" {
			replaceTags := strings.Split(cmd.Flag("replace-tags").Value.String(), ",")
			api.SetTags(ids, replaceTags)
		}
	},
}

func init() {
	imagesSearchCmd.PersistentFlags().StringP("tags", "t", "", "Filter by tags")
	imagesSearchCmd.PersistentFlags().String("set-album", "", "Move filtered images to album")
	imagesSearchCmd.PersistentFlags().String("add-tags", "", "Add tags to filtered images")
	imagesSearchCmd.PersistentFlags().String("replace-tags", "", "Replace tags of filtered images")
}
