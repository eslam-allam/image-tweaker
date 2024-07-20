/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/eslam-allam/image-tweaker/internal/image"
	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag/v2"
)

var (
	resize       bool
	targetWidth  uint
	targetHeight uint

	targetFormat image.ImgFormat = image.UNSUPPORTED

	outPath string

	threadCount uint
)

// transformCmd represents the transform command
var transformCmd = &cobra.Command{
	Use:   "transform [image|directory]",
	Args:  cobra.ExactArgs(1),
	Short: "Do multiple operations",
	Long: `Transform image by doing multiple operations such as resizing and
changing the format. 

Example:
	image-tweaker transform --resize -w 500 --format webp

this will resize the image to be 500px wide and convert the format to webp`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("transform called with args: %v\n", args)
		path, err := filepath.Abs(args[0])
		if err != nil {
			fmt.Println("invalid path provided")
			os.Exit(1)
		}
		err = image.Transform(path, outPath, resize, targetWidth, targetHeight, targetFormat, threadCount)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(transformCmd)

	transformCmd.Flags().BoolVarP(&resize, "resize", "r", false, "Resize image (used with width and height flags)")
	transformCmd.Flags().UintVarP(&targetWidth, "target-width", "w", 0, "Target max width for image if resize option is used (0) for no change")
	transformCmd.Flags().UintVarP(&targetHeight, "target-height", "v", 0, "Target max height for image if resize option is used (0) for no change")

	te := enumflag.NewWithoutDefault(
		&targetFormat,
		"format",
		image.GetFormatNames(),
		enumflag.EnumCaseInsensitive,
	)
	te.RegisterCompletion(transformCmd, "format",
		enumflag.Help[image.ImgFormat]{
			image.JPEG: "jpeg image format",
			image.PNG:  "Png image format",
			image.WEBP: "Webp image format",
		})

	transformCmd.Flags().VarP(te, "format", "f", "Save output image in this format. Valid values are 'jpeg', 'png', and 'webp'")

	transformCmd.Flags().StringVarP(&outPath, "output", "o", "output", "output directory or file")
	transformCmd.Flags().UintVarP(&threadCount, "threads", "t", 0, "[DANGER!] Number of images proccessed at a time (CPU count by default)")
}
