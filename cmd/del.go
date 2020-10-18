package cmd

import (
	"github.com/ashrithr/s3-restore/internal/restore"
	"github.com/spf13/cobra"
)

// delMarkerCmd represents the delMarker command
var delMarkerCmd = &cobra.Command{
	Use:   "del",
	Short: "Restores the deleted objects by removing the DELETE markers.",
	Run: func(cmd *cobra.Command, args []string) {
		bucket, _ = cmd.Flags().GetString("bucket")
		prefix, _ = cmd.Flags().GetString("prefix")
		region, _ = cmd.Flags().GetString("region")
		dryrun, _ = cmd.Flags().GetBool("dryrun")

		restore.ObjsUsingDel(bucket, prefix, region, dryrun)
	},
}

func init() {
	rootCmd.AddCommand(delMarkerCmd)
}
