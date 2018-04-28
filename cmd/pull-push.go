// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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
	"fmt"
	"os"
	"strings"

	"github.com/bharat-p/goutils/cli"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

// pushCmd represents the push command
var pullPushCmd = &cobra.Command{
	Use:   "pull-push",
	Short: "Pull and push images between repositories",
	Long:  `Pull images from one docker registry and push it to another.`,
	Run: func(cmd *cobra.Command, args []string) {
		var pull, push, tagLocal bool
		if Verbose {
			jww.SetStdoutThreshold(jww.LevelInfo)
		}

		// TODO: Work your own magic here
		pullFrom, _ := cmd.Flags().GetString("from")
		pushTo, _ := cmd.Flags().GetString("to")

		pullFrom = strings.TrimSuffix(strings.TrimSpace(pullFrom), "/")
		pushTo = strings.TrimSuffix(strings.TrimSpace(pushTo), "/")
		if pullFrom == "" && pushTo == "" {
			jww.ERROR.Printf("Error: %s\n", "Please specify either --from or --to ")
			cmd.Usage()
			os.Exit(1)
		}
		tagLocal, _ = cmd.Flags().GetBool("local")

		if pushTo != "" {
			push = true
			pushTo += "/"
		}

		if pullFrom != "" {
			pull = true
			pullFrom += "/"
		} else if tagLocal {
			jww.ERROR.Printf("Error: %s\n", "--local ( -l) can only be used with --from flag ")
			cmd.Usage()
			os.Exit(1)
		}

		varRemoveImage, _ := cmd.Flags().GetBool("remove")

		imageNames := []string{}
		images, _ := cmd.Flags().GetStringArray("images")
		for _, imageOption := range images {
			for _, imageVal := range strings.Split(imageOption, ",") {
				imageVal = strings.TrimSpace(imageVal)
				if imageVal != "" {
					imageNames = append(imageNames, imageVal)
				}
			}

		}
		if len(imageNames) <= 0 {
			jww.ERROR.Printf("Error: %s\n", "Please specify atlease one image name.")
			cmd.Usage()
			os.Exit(1)
		}
		for _, image := range imageNames {

			pullImageName := fmt.Sprintf("%s%s", pullFrom, image)
			if pull {
				err := pullImage(pullImageName)
				if err != nil {
					jww.ERROR.Printf("Failed to pull image: %s\n", pullImageName)
					jww.ERROR.Printf("%s\n", err)
					continue
				}
				if tagLocal {
					tagImage(pullImageName, image)
				}

			}
			if push {
				pushImageName := fmt.Sprintf("%s%s", pushTo, image)
				tagImage(pullImageName, pushImageName)
				err := pushImage(pushImageName)
				if err != nil {
					fmt.Printf("Error:%s\n", err)
				} else if varRemoveImage {
					removeImage(pushImageName)
				}
			}
			if varRemoveImage && pull {
				removeImage(pullImageName)

			}

		}

	},
}

func init() {

	RootCmd.AddCommand(pullPushCmd)
	pullPushCmd.Flags().StringP("from", "f", "", "Registry from where to pull images")
	pullPushCmd.Flags().StringP("to", "t", "", "Registry to push images")
	pullPushCmd.Flags().BoolP("local", "l", false, "Tag as local (applicable with --from only if --to is not used)")
	pullPushCmd.Flags().BoolP("remove", "r", false, "Remove images after pull/push is done.")
	pullPushCmd.Flags().StringArrayP("images", "i", []string{}, "Image(s) to pull/push")
}

func pullImage(imageName string) (err error) {
	jww.INFO.Printf("Pulling image %s", imageName)
	_, err = cli.RunCommand(DockerBinary, "pull", imageName)

	return err
}

func removeImage(imageName string) (err error) {
	jww.INFO.Printf("Removing image %s", imageName)
	_, err = cli.RunCommand(DockerBinary, "rmi", imageName)

	return err
}

func tagImage(imageName string, tagAs string) (err error) {
	jww.INFO.Printf("Tagging image %s as %s", imageName, tagAs)
	_, err = cli.RunCommand(DockerBinary, "tag", imageName, tagAs)
	return err
}

func pushImage(imageName string) (err error) {
	jww.INFO.Printf("Pushing image %s", imageName)
	_, err = cli.RunCommand(DockerBinary, "push", imageName)
	return err
}
