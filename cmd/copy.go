package cmd

import (
	"github.com/ashrithr/s3-restore/internal/restore"
	"github.com/spf13/cobra"
)

var (
	dstBucket string
	dstPrefix string
)

// copyObjCmd represents the copyObj command
var copyObjCmd = &cobra.Command{
	Use:   "copy",
	Short: "Restores the deleted objects by copying them from their previous versions.",
	Long:  `Restores the deleted objects by copying them from their previous versions.`,
	Run: func(cmd *cobra.Command, args []string) {
		bucket, _ = cmd.Flags().GetString("bucket")
		prefix, _ = cmd.Flags().GetString("prefix")
		region, _ = cmd.Flags().GetString("region")
		dstBucket, _ = cmd.Flags().GetString("dstBucket")
		dstPrefix, _ = cmd.Flags().GetString("dstPrefix")
		dryrun, _ = cmd.Flags().GetBool("dryrun")

		restore.ObjsUingCopy(bucket, prefix, region, dstBucket, dstPrefix, dryrun)
	},
}

func init() {
	rootCmd.AddCommand(copyObjCmd)

	copyObjCmd.PersistentFlags().String("dstBucket", "", "Optionally specify another destination bucket to copy the objects to")
	copyObjCmd.PersistentFlags().String("dstPrefix", "", "Optionally specify destination prefix")
}
