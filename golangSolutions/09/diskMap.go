package main

import "log/slog"

type FileInformation struct {
	startBlockIndex int
	numBlocks       int
	fileID          int
}

func (file FileInformation) computeChecksumContribution() int {
	// slog.Debug("checksum computation", "interval", interval, "start disk index", startDiskIndex, "end disk index", endDiskIndex)
	endBlockIndex := file.startBlockIndex + file.numBlocks
	return file.fileID * ((endBlockIndex * (endBlockIndex - 1) / 2) - (file.startBlockIndex * (file.startBlockIndex - 1) / 2))
}

type DiskMap struct {
	fileList        *LinkedList
	totalDiskLength int
}

func ParseLineToDiskMap(line string) *DiskMap {
	diskMap := DiskMap{
		fileList: NewDiskMapLinkedList(),
	}

	currentBlockIndex := 0
	for runeIndex, runeValue := range line {
		fileNumBlocks := (int(runeValue) - '0')
		if runeIndex%2 == 0 {
			file := &FileInformation{
				startBlockIndex: currentBlockIndex,
				numBlocks:       fileNumBlocks,
				fileID:          runeIndex / 2,
			}
			diskMap.fileList.Add(file)
			slog.Debug("parsed next interval", "new file", file)
		}
		currentBlockIndex += fileNumBlocks
	}

	finalMapInterval := diskMap.fileList.tail
	diskMap.totalDiskLength = finalMapInterval.fileInfo.startBlockIndex + finalMapInterval.fileInfo.numBlocks

	return &diskMap
}

func (diskMap *DiskMap) ComputeChecksum() int {
	checksum := 0
	for currentFileNode := diskMap.fileList.head; currentFileNode != nil; currentFileNode = currentFileNode.next {
		checksumContribution := currentFileNode.fileInfo.computeChecksumContribution()
		checksum += checksumContribution
		slog.Debug(
			"computing checksum",
			"file information", currentFileNode.fileInfo,
			"checksum contribution", checksumContribution,
			"updated checksum", checksum,
		)
	}
	return checksum
}

func (diskMap *DiskMap) DefragmentMoveBlocks() {
	currentForwardFile := diskMap.fileList.head
	currentReverseFile := diskMap.fileList.tail
	currentBlockGap := 0

	for currentForwardFile != currentReverseFile && currentForwardFile.next != nil {
		currentBlockGap = currentForwardFile.next.fileInfo.startBlockIndex - (currentForwardFile.fileInfo.startBlockIndex + currentForwardFile.fileInfo.numBlocks)
		if currentBlockGap == 0 {
			currentForwardFile = currentForwardFile.next
			continue
		}

		// We have a block gap, fill with blocks from reverse files
		// Three cases:
		// 	- we can fit the reverse file into this space
		// 	- we cannot fit the reverse file into this space

		if currentBlockGap >= currentReverseFile.fileInfo.numBlocks {
			// Case one:
			// Splice the current reverse file out of its current position
			// Splice the entire reverse file into this position
			// Update the current reverse file to the old previous
			// Update the current forward file to the new next

			nextReverseFile := currentReverseFile.prev
			diskMap.fileList.SpliceOut(currentReverseFile)
			currentReverseFile.fileInfo.startBlockIndex = currentForwardFile.fileInfo.startBlockIndex + currentForwardFile.fileInfo.numBlocks
			diskMap.fileList.SpliceIn(currentForwardFile, currentReverseFile)
			currentReverseFile = nextReverseFile
			currentForwardFile = currentForwardFile.next
		} else {
			// Case two:
			// Decrement the size of the current reverse file by a corresponding amount (shift that many blocks into the gap)
			// Make a new file with an amount of blocks equal to the gap, splicing it in
			// Shift the current forward file forward to the new file

			currentReverseFile.fileInfo.numBlocks -= currentBlockGap
			newFile := &FileInformation{
				startBlockIndex: currentForwardFile.fileInfo.startBlockIndex + currentForwardFile.fileInfo.numBlocks,
				numBlocks:       currentBlockGap,
				fileID:          currentReverseFile.fileInfo.fileID,
			}
			newFileNode := &LinkedListNode{fileInfo: newFile}
			diskMap.fileList.SpliceIn(currentForwardFile, newFileNode)
			currentForwardFile = currentForwardFile.next
		}
	}
}

func (diskMap *DiskMap) DefragmentMoveFiles() {
	currentReverseFile := diskMap.fileList.tail
reverseFileLoop:
	for currentReverseFile != nil {
		// Look forward (starting from the head and ending at the current reverse file)
		// looking for the first gap that could house the current reverse file, splicing in there

		slog.Debug("reverse file loop", "current reverse file info", currentReverseFile.fileInfo)

		for currentForwardFile := diskMap.fileList.head; currentForwardFile != currentReverseFile; currentForwardFile = currentForwardFile.next {
			currentBlockGap := currentForwardFile.next.fileInfo.startBlockIndex - (currentForwardFile.fileInfo.startBlockIndex + currentForwardFile.fileInfo.numBlocks)
			slog.Debug("forward file loop", "current reverse file info", currentReverseFile.fileInfo, "current forward file info", currentForwardFile.fileInfo, "current block gap", currentBlockGap)

			if currentBlockGap >= currentReverseFile.fileInfo.numBlocks {
				// Splice the current reverse file out of its current position
				// Splice the entire reverse file into this position
				// Update the current reverse file to the old previous
				// Update the current forward file to the new next

				nextReverseFile := currentReverseFile.prev
				diskMap.fileList.SpliceOut(currentReverseFile)
				currentReverseFile.fileInfo.startBlockIndex = currentForwardFile.fileInfo.startBlockIndex + currentForwardFile.fileInfo.numBlocks
				diskMap.fileList.SpliceIn(currentForwardFile, currentReverseFile)
				currentReverseFile = nextReverseFile

				slog.Debug("found gap for current reverse file", "new reverse file info", currentReverseFile.fileInfo)
				continue reverseFileLoop
			}
		}

		// If we made it here, we did not find a home for the poor reverse file, so we move to the next reverse file (i.e. currentReverseFile.prev)
		currentReverseFile = currentReverseFile.prev
	}
}
