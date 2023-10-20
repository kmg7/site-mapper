package cmd

import (
	mapper "github.com/kmg7/site-mapper/pkg/mapper"
	"github.com/spf13/cobra"
)

// mapCmd represents the map command
var mapCmd = &cobra.Command{
	Use:   "map",
	Short: "Gives a sitemap of any website",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		mapper.MapUrl(args[0], Depth)
	},
}

var Depth int

func init() {
	rootCmd.AddCommand(mapCmd)
	mapCmd.Flags().IntVarP(&Depth, "depth", "d", 3, "Search depth for sitemap")

}
