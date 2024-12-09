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
