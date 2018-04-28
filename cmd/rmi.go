// Copyright © 2017 Bharat Patel
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

	"github.com/bharat-p/goutils/cli"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"strings"
)

// rmiCmd represents the rmi command
var rmiCmd = &cobra.Command{
	Use:   "rmi",
	Short: "Remove docker images by specifying repository name or tag name",
	Long: `Remove docker images by specifying repository name or tag name, For example
:
To remove all images with repository ubuntu
docker-utils rmi -r "ubuntu"

`,
	Run: func(cmd *cobra.Command, args []string) {

		repository, _ := cmd.Flags().GetString("repository")
		tag, _ := cmd.Flags().GetString("tag")
		if repository == "" && tag == "" {
			jww.ERROR.Printf("Error: %s\n", "Please specify either --repository or --tag ")
			cmd.Usage()
			os.Exit(1)
		}
		grepCmd := "egrep '"
		if repository != "" {
			grepCmd += "^" + repository
		} else {
			grepCmd += ".*"
		}
		grepCmd += ":"
		if tag != "" {
			grepCmd += tag + ""

		} else {
			grepCmd += ".*"
		}
		grepCmd += ",'"

		oldDryRun := cli.DryRun
		cli.DryRun = false
		myCmd := fmt.Sprintf("docker images --format {{.Repository}}:{{.Tag}},{{.ID}} | %s | awk -F  ',' '{print $2}' ", grepCmd)
		out, _, _ := cli.GetOutputOfCommand("/bin/bash", "-c", myCmd)
		cli.DryRun = oldDryRun
		images := strings.Split(out, "\n")
		baseCmd :=  []string{"rmi"}
		if flagVal, _ := cmd.Flags().GetBool("force"); flagVal {
			baseCmd = append(baseCmd, "-f")
		}
		for _, imageName := range images {
			if imageName != "" {
				jww.INFO.Printf("Attempting to remove image: %s", imageName)
				exitCode, err:= cli.RunCommand(DockerBinary, append(baseCmd, imageName)...)
				if exitCode == 0 {
					jww.INFO.Printf("Removed image: %s", imageName)
				} else {
					jww.WARN.Print(err)
				}

			}

		}

	},
}

func init() {

	RootCmd.AddCommand(rmiCmd)

	rmiCmd.Flags().StringP("repository", "r", "", "Repository name to remove images for")
	rmiCmd.Flags().StringP("tag", "t", "", "Remove all image matching tag")
	rmiCmd.Flags().BoolP("force", "f", false, "Force removal of image(s)")
}
