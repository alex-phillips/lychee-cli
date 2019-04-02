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
	"github.com/spf13/cobra"
)

// downloadThreadCmd represents the freeleech command
var albumCmd = &cobra.Command{
	Use:   "albums",
	Short: "Download all media from a subreddit",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// home, _ := homedir.Dir()
		// thread := args[0]

		// reddit.SaveDir = cmd.Flag("save-dir").Value.String()
		// if reddit.SaveDir == "" {
		// 	reddit.SaveDir = home + "/go"
		// }
		// reddit.SaveDir = strings.TrimRight(reddit.SaveDir, "/") + "/" + thread

		// log.Debug.Println("Saving files to " + reddit.SaveDir)

		// if _, err := os.Stat(reddit.SaveDir); os.IsNotExist(err) {
		// 	os.MkdirAll(reddit.SaveDir, 0700)
		// }

		// if reddit.DownloadComments == true {
		// 	log.Debug.Println("Downloading comments")
		// }

		// reddit.DownloadPost("/"+args[0]+"/", "", "", 0)
	},
}

func init() {
	albumCmd.AddCommand(albumCreateCmd)
	rootCmd.AddCommand(albumCmd)
}
