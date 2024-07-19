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
)

var (
	resize       bool
	targetWidth  uint
	targetHeight uint

	targetEncoding string
)

// transformCmd represents the transform command
var transformCmd = &cobra.Command{
	Use:   "transform [image|directory]",
	Args:  cobra.ExactArgs(1),
	Short: "Do multiple operations",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("transform called with args: %v\n", args)
		path, err := filepath.Abs(args[0])
		if err != nil {
			fmt.Println("invalid path provided")
		}
		img, format, err := image.ReadImage(path)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if resize {
			img = image.ResizeIfBigger(img, targetWidth, targetHeight)
		}

		enc := format
		if targetEncoding != "" {
			enc, err = image.EncodingFromFormat(targetEncoding)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		outPath := "out." + enc.Extension()
		outPath, err = filepath.Abs(outPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = image.SaveImage(img, enc, outPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(format)
	},
}

func init() {
	rootCmd.AddCommand(transformCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:

	transformCmd.Flags().BoolVarP(&resize, "resize", "r", false, "Resize image (used with width and height flags)")
	transformCmd.Flags().UintVarP(&targetWidth, "target-width", "w", 0, "Target max width for image if resize option is used (0) for no change")
	transformCmd.Flags().UintVarP(&targetHeight, "target-height", "v", 0, "Target max height for image if resize option is used (0) for no change")

	transformCmd.Flags().StringVarP(&targetEncoding, "encoding", "e", "", "Save output image in this format. Valid values are 'jpeg', 'png', and 'webp'")
}
