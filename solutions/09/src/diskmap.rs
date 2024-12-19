use std::{error::Error, fmt::Display};

#[allow(unused_imports)]
use tracing::{debug, error, info, trace};

#[allow(dead_code)]
#[derive(Debug, Clone, Copy)]
pub enum FileType {
    Empty,
    File(i64),
}

#[derive(Debug)]
pub struct FileBlock {
    filetype: FileType,
    start: usize,
    length: usize,
}

#[derive(Debug)]
pub struct DiskMap {
    file_blocks: Vec<FileBlock>,
}

pub fn parse_input_to_diskmap(diskmap_specification: Vec<u8>) -> DiskMap {
    let mut current_fileid = 0;
    let mut current_block_index = 0;

    let mut file_blocks = Vec::new();
    for (specification_index, specification_character) in diskmap_specification.iter().enumerate() {
        let block_length = (*specification_character - b'0') as usize;
        if specification_index % 2 == 1 {
            // file_blocks.push(FileBlock {
            //     filetype: FileType::Empty,
            //     start: current_block_index,
            //     length: block_length,
            // });
            // debug! {"next FileBlock"=?file_blocks.last().unwrap(), "parsed block"}
        } else {
            file_blocks.push(FileBlock {
                filetype: FileType::File(current_fileid),
                start: current_block_index,
                length: block_length,
            });
            debug! {"next FileBlock"=?file_blocks.last().unwrap(), "parsed block"}
            current_fileid += 1;
        }
        current_block_index += block_length;
    }

    DiskMap { file_blocks }
}

impl Display for DiskMap {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        for (index, file) in self.file_blocks.iter().enumerate() {
            writeln!(f, "{index}: {:?}", file).unwrap();
        }

        Ok(())
    }
}

impl DiskMap {
    pub fn compute_checksum(&self) -> i64 {
        self.file_blocks
            .iter()
            .map(|file_block| {
                if let FileType::File(fileid) = file_block.filetype {
                    let start_index = file_block.start as i64;
                    let end_index = (file_block.start + file_block.length) as i64;
                    let index_sum =
                        (end_index * (end_index - 1) / 2) - (start_index * (start_index - 1) / 2);
                    trace!(?file_block, ?end_index, ?index_sum, "checksum_calculation");

                    fileid * index_sum as i64
                } else {
                    0
                }
            })
            .sum()
    }

    pub fn defragment_block_wise(&mut self) -> Result<(), Box<dyn Error>> {
        let mut current_forward_file_index = 0;
        let Some(mut current_reverse_file_index) = self
            .file_blocks
            .iter()
            .enumerate()
            .rev()
            .find(|(_, item)| matches!(item.filetype, FileType::File(_)))
            .map(|(index, _)| index)
        else {
            return Err("cannot defragment an empty disk".into());
        };

        trace!("prior disk map"=%format!("\n{self}"), ?current_forward_file_index, ?current_reverse_file_index, "before updating");

        while current_forward_file_index < current_reverse_file_index {
            let current_forward_file = &self.file_blocks[current_forward_file_index];
            let next_forward_file = &self.file_blocks[current_forward_file_index + 1];
            let current_forward_file_end_index =
                current_forward_file.start + current_forward_file.length;
            let current_block_gap = next_forward_file.start - current_forward_file_end_index;
            if current_block_gap == 0 {
                current_forward_file_index += 1;
                trace!(
                    ?current_forward_file_index,
                    ?current_reverse_file_index,
                    "block gap zero"
                );
                continue;
            }

            let current_reverse_file = &self.file_blocks[current_reverse_file_index];

            if current_block_gap >= current_reverse_file.length {
                let mut current_reverse_file = self.file_blocks.remove(current_reverse_file_index);
                current_reverse_file.start = current_forward_file_end_index;
                self.file_blocks
                    .insert(current_forward_file_index + 1, current_reverse_file);
                current_forward_file_index += 1;
                trace!("updated disk map"=%format!("\n{self}"), ?current_forward_file_index, ?current_reverse_file_index, "current block gap larger than or equal to reverse file")
            } else {
                let current_reverse_file = &mut self.file_blocks[current_reverse_file_index];
                current_reverse_file.length -= current_block_gap;

                let new_file = FileBlock {
                    filetype: current_reverse_file.filetype,
                    start: current_forward_file_end_index,
                    length: current_block_gap,
                };
                self.file_blocks
                    .insert(current_forward_file_index + 1, new_file);
                current_forward_file_index += 1;
                current_reverse_file_index += 1;
                trace!("updated disk map"=%format!("\n{self}"), ?current_forward_file_index, ?current_reverse_file_index, "current block gap smaller than reverse file")
            }
        }

        Ok(())
    }

    pub fn defragment_file_wise(&mut self) -> Result<(), Box<dyn Error>> {
        if self.file_blocks.len() == 0 {
            return Err("cannot defragment an empty disk".into());
        }
        trace!("prior disk map"=%format!("\n{self}"), "before updating");

        let mut current_reverse_file_index = self.file_blocks.len()-1;
        'reverse_file_loop: while current_reverse_file_index > 0 {
            let current_reverse_file_length = self.file_blocks[current_reverse_file_index].length;
            for current_forward_file_index in 0..current_reverse_file_index {
                let current_forward_file = &self.file_blocks[current_forward_file_index];
                let next_forward_file = &self.file_blocks[current_forward_file_index + 1];
                let current_forward_file_end_index =
                current_forward_file.start + current_forward_file.length;
                let current_block_gap = next_forward_file.start - current_forward_file_end_index;
                // trace!(?current_forward_file_index, ?current_reverse_file_index, ?current_reverse_file_length, ?current_block_gap);
                if current_block_gap >= current_reverse_file_length {
                    let mut current_reverse_file =
                        self.file_blocks.remove(current_reverse_file_index);
                    current_reverse_file.start = current_forward_file_end_index;
                    self.file_blocks
                        .insert(current_forward_file_index + 1, current_reverse_file);
                    trace!("updated disk map"=%format!("\n{self}"), "inserted index"=?current_forward_file_index+1, "found gap large enough for reverse file");
                    continue 'reverse_file_loop;
                }
            }
            // trace!("updated disk map"=%format!("\n{self}"), ?current_reverse_file_index, "no gap found for reverse file");
            current_reverse_file_index-=1;
        }

        Ok(())
    }
}
