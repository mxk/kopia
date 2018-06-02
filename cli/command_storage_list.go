package cli

import (
	"context"
	"fmt"

	"github.com/kopia/kopia/repo"
)

var (
	storageListCommand = storageCommands.Command("list", "List storage blocks").Alias("ls")
	storageListPrefix  = storageListCommand.Flag("prefix", "Block prefix").String()
	storageListMinSize = storageListCommand.Flag("min-size", "Minimum size").Int64()
	storageListMaxSize = storageListCommand.Flag("max-size", "Maximum size").Int64()
)

func runListStorageBlocks(ctx context.Context, rep *repo.Repository) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	ch := rep.Storage.ListBlocks(ctx, *storageListPrefix)

	for b := range ch {
		if b.Error != nil {
			return b.Error
		}

		if *storageListMaxSize != 0 && b.Length > *storageListMaxSize {
			continue
		}

		if *storageListMinSize != 0 && b.Length < *storageListMinSize {
			continue
		}

		fmt.Printf("%-70v %10v %v\n", b.BlockID, b.Length, b.TimeStamp.Local().Format(timeFormat))
	}

	return nil
}

func init() {
	storageListCommand.Action(repositoryAction(runListStorageBlocks))
}
