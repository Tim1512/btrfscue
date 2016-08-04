/*
 * btrfscue version 0.3
 * Copyright (c)2011-2016 Christian Blichmann
 *
 * BTRFS filesystem structures
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *     * Redistributions of source code must retain the above copyright
 *       notice, this list of conditions and the following disclaimer.
 *     * Redistributions in binary form must reproduce the above copyright
 *       notice, this list of conditions and the following disclaimer in the
 *       documentation and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 */

package btrfs // import "blichmann.eu/code/btrfscue/btrfs"

import (
	"blichmann.eu/code/btrfscue/uuid"
	"time"
)

const (
	// "_BHRfS_M" in little-endian
	Magic = 0x4D5F53665248425F

	DefaultBlockSize = 1 << 12
)

// Offsets of all superblock copies
const (
	SuperInfoOffset  = 0x10000         // 64 KiB
	SuperInfoOffset2 = 0x4000000       // 64 MiB
	SuperInfoOffset3 = 0x4000000000    // 256 GiB
	SuperInfoOffset4 = 0x4000000000000 // 1 PiB
)

// Object ids
const (
	// Holds pointers to all of the tree roots
	RootTreeObjectId = 1

	// Stores information about which extents are in use, and reference
	// counts
	ExtentTreeObjectId = 2

	// The chunk tree stores translations from logical -> physical block
	// numbering the super block points to the chunk tree
	ChunkTreeObjectId = 3

	// Stores information about which areas of a given device are in use. One
	// per device. The tree of tree roots points to the device tree.
	DevTreeObjectId = 4

	// One per subvolume, storing files and directories
	FSTreeObjectId = 5

	// Directory objectid inside the root tree
	RootTreeDirObjectId = 6

	// Holds checksums of all the data extents
	CSumTreeObjectId = 7

	// Orhpan objectid for tracking unlinked/truncated files
	OrphanObjectId = ^uint64(5) + 1

	// Does write ahead logging to speed up fsyncs
	TreeLogObjectId      = ^uint64(6) + 1
	TreeLogFixupObjectId = ^uint64(7) + 1

	// For space balancing
	TreeRelocObjectId     = ^uint64(8) + 1
	DataRelocTreeObjectId = ^uint64(9) + 1

	// Extent checksums all have this objectid. This allows them to share the
	// logging tree for fsyncs.
	ExtentCSumObjectId = ^uint64(10) + 1

	// For storing free space cache */
	FreeSpaceObjectId = ^uint64(11) + 1

	// Dummy objectid represents multiple objectids
	MultipleObjectIdS = ^uint64(255) + 1

	// All files have objectids in this range
	FirstFreeObjectId      = 256
	LastFreeObjectId       = ^uint64(256) + 1
	FirstChunkTreeObjectId = 256

	// The device items go into the chunk tree. The key is in the form
	// [ 1 DevItemKey device_id ]
	DevItemsObjectId = 1

	BtreeInodeObjectId = 1

	EmptySubvolDirObjectId = 2
)

// Entity sizes
const (
	CSumSize             = 32
	LabelSize            = 256
	SystemChunkArraySize = 2048
)

// Key types
const (
	// Inode items have the data typically returned from stat and store other
	// info about object characteristics. There is one for every file and dir
	// in the FS.
	InodeItemKey = 1

	InodeRefKey    = 12
	InodeExtrefKey = 13
	XAttrItemKey   = 24
	OrphanItemKey  = 48

	// dir items are the name -> inode pointers in a directory. There is one
	// for every name in a directory.
	DirLogItemKey  = 60
	DirLogIndexKey = 72
	DirItemKey     = 84
	DirIndexKey    = 96

	// Extent data is for file data.
	ExtentDataKey = 108

	// Extent csums are stored in a separate tree and hold csums for
	// an entire extent on disk.
	ExtentCSumKey = 128

	// root items point to tree roots. They are typically in the root
	// tree used by the super block to find all the other trees.
	RootItemKey = 132

	// Root backrefs tie subvols and snapshots to the directory entries that
	// reference them.
	RootBackRefKey = 144

	// Root refs make a fast index for listing all of the snapshots and
	// subvolumes referenced by a given root. They point directly to the
	// directory item in the root that references the subvol.
	RootRefKey = 156

	// Extent items are in the extent map tree. These record which blocks
	// are used, and how many references there are to each block.
	ExtentItemKey = 168

	// The same as the BTRFS_EXTENT_ITEM_KEY, except it's metadata we already
	// know the length, so we save the level in key->offset instead of the
	// length.
	MetadataItemKey = 169

	TreeBlockRefKey   = 176
	ExtentDataRefKey  = 178
	ExtentRefV0Key    = 180
	SharedBlockRefKey = 182
	SharedDataRefKey  = 184

	// Block groups give us hints into the extent allocation trees. Which
	// blocks are free etc.
	BlockGroupItemKey = 192

	// Every block group is represented in the free space tree by a free space
	// info item, which stores some accounting information. It is keyed on
	// (block_group_start, FREE_SPACE_INFO, block_group_length).
	FreeSpaceInfoKey = 198

	// A free space extent tracks an extent of space that is free in a block
	// group. It is keyed on (start, FREE_SPACE_EXTENT, length).
	FreeSpaceExtentKey = 199

	// When a block group becomes very fragmented, we convert it to use bitmaps
	// instead of extents. A free space bitmap is keyed on
	// (start, FREE_SPACE_BITMAP, length); the corresponding item is a bitmap
	// with (length / sectorsize) bits.
	FreeSpaceBitmapKey = 200

	DevExtentKey = 204
	DevItemKey   = 216
	ChunkItemKey = 228

	// Records the overall state of the qgroups.
	// There's only one instance of this key present,
	// (0, BTRFS_QGROUP_STATUS_KEY, 0)
	QgroupStatusKey = 240

	// Records the currently used space of the qgroup.
	// One key per qgroup, (0, BTRFS_QGROUP_INFO_KEY, qgroupid).
	QgroupInfoKey = 242

	// Contains the user configured limits for the qgroup.
	// One key per qgroup, (0, BTRFS_QGROUP_LIMIT_KEY, qgroupid).
	QgroupLimitKey = 244

	// Records the child-parent relationship of qgroups. For each relation, 2
	// keys are present:
	// (childid, BTRFS_QGROUP_RELATION_KEY, parentid)
	// (parentid, BTRFS_QGROUP_RELATION_KEY, childid)
	QgroupRelationKey = 246

	BalanceItemKey = 248

	// Persistantly stores the io stats in the device tree.
	// One key for all stats, (0, BTRFS_DEV_STATS_KEY, devid).
	DevStatsKey = 249

	// Persistantly stores the device replace state in the device tree.
	// The key is built like this: (0, BTRFS_DEV_REPLACE_KEY, 0).
	DevReplaceKey = 250

	// String items are for debugging. They just store a short string of data
	// in the FS.
	StringItemKey = 253
)

type CSum [CSumSize]byte

type Header struct {
	CSum CSum
	// The following three fields must match struct SuperBlock
	// File system specific UUID
	FSID uuid.UUID
	// The start of this block relative to the begining of the backing device
	ByteNr uint64
	Flags  uint64
	// Allowed to be different from SuperBlock from here on
	ChunkTreeUUID uuid.UUID
	Generation    uint64
	Owner         uint64
	NrItems       uint32
	Level         uint8
}

func (h *Header) Parse(b *ParseBuffer) {
	copy(h.CSum[:], b.Next(CSumSize))
	copy(h.FSID[:], b.Next(uuid.UUIDSize))
	h.ByteNr = b.NextUint64()
	h.Flags = b.NextUint64()
	copy(h.ChunkTreeUUID[:], b.Next(uuid.UUIDSize))
	h.Generation = b.NextUint64()
	h.Owner = b.NextUint64()
	h.NrItems = b.NextUint32()
	h.Level = b.NextUint8()
}

func (h *Header) IsLeaf() bool {
	return h.Level == 0
}

type Key struct {
	ObjectId uint64
	Type     uint8
	Offset   uint64
}

func (k *Key) Parse(b *ParseBuffer) {
	k.ObjectId = b.NextUint64()
	k.Type = b.NextUint8()
	k.Offset = b.NextUint64()
}

type itemData interface {
	Parse(b *ParseBuffer)
}

type Item struct {
	Key
	Offset uint32
	Size   uint32
	Data   itemData
}

func (i *Item) Parse(b *ParseBuffer) {
	i.Key.Parse(b)
	i.Offset = b.NextUint32()
	i.Size = b.NextUint32()
}

func (i *Item) ParseData(b *ParseBuffer) {
	switch i.Type {
	case InodeItemKey:
		i.Data = &InodeItem{}
	case InodeRefKey:
		i.Data = &InodeRefItem{}
	case XAttrItemKey:
		fallthrough
	case DirItemKey:
		fallthrough
	case DirIndexKey:
		i.Data = &DirItem{}
	case ExtentDataKey:
		i.Data = &FileExtentItem{}
	case ExtentCSumKey:
		i.Data = &CSumItem{}
	case RootItemKey:
		i.Data = &RootItem{}
	case RootBackRefKey:
		fallthrough
	case RootRefKey:
		i.Data = &RootRef{}
	case ExtentItemKey:
		i.Data = &ExtentItem{}
	case BlockGroupItemKey:
		i.Data = &BlockGroupItem{}
	default:
		return
	}
	i.Data.Parse(b)
}

type InodeItem struct {
	// NFS style generation number
	Generation uint64
	// Transid that last touched this inode
	Transid    uint64
	Size       uint64
	Nbytes     uint64
	BlockGroup uint64
	Nlink      uint32
	Uid        uint32
	Gid        uint32
	Mode       uint32
	Rdev       uint64
	Flags      uint64

	// Modification sequence number for NFS
	Sequence uint64

	// A little future expansion, for more than this we can just grow the
	// inode item and version it.
	reserved [4]uint64
	Atime    time.Time
	Ctime    time.Time
	Mtime    time.Time
	Otime    time.Time
}

func (i *InodeItem) Parse(b *ParseBuffer) {
	i.Generation = b.NextUint64()
	i.Transid = b.NextUint64()
	i.Size = b.NextUint64()
	i.Nbytes = b.NextUint64()
	i.BlockGroup = b.NextUint64()
	i.Nlink = b.NextUint32()
	i.Uid = b.NextUint32()
	i.Gid = b.NextUint32()
	i.Mode = b.NextUint32()
	i.Rdev = b.NextUint64()
	i.Flags = b.NextUint64()
	i.Sequence = b.NextUint64()
	b.Next(4 * 8)
	i.Atime = time.Unix(int64(b.NextUint64()), int64(b.NextUint32()))
	i.Ctime = time.Unix(int64(b.NextUint64()), int64(b.NextUint32()))
	i.Mtime = time.Unix(int64(b.NextUint64()), int64(b.NextUint32()))
	i.Otime = time.Unix(int64(b.NextUint64()), int64(b.NextUint32()))
}

type InodeRefItem struct {
	Index   uint64
	NameLen uint16
	Name    string
}

func (i *InodeRefItem) Parse(b *ParseBuffer) {
	i.Index = b.NextUint64()
	i.NameLen = b.NextUint16()
	l := int(i.NameLen)
	if l > 255 {
		l = 255
	}
	i.Name = string(b.Next(l))
}

type DirItem struct {
	Location Key
	TransId  uint64
	DataLen  uint16
	NameLen  uint16
	Type     uint8
	Name     string
	Data     string
}

func (i *DirItem) Parse(b *ParseBuffer) {
	i.Location.Parse(b)
	i.TransId = b.NextUint64()
	i.DataLen = b.NextUint16()
	i.NameLen = b.NextUint16()
	i.Type = b.NextUint8()
	l := int(i.NameLen)
	if l > 255 {
		l = 255
	}
	i.Name = string(b.Next(l))
	l = int(i.DataLen)
	if l > DefaultBlockSize {
		l = DefaultBlockSize
	}
	i.Data = string(b.Next(l))
}

type BlockGroupItem struct {
	Used          uint64
	ChunkObjectId uint64
	Flags         uint64
}

func (i *BlockGroupItem) Parse(b *ParseBuffer) {
	i.Used = b.NextUint64()
	i.ChunkObjectId = b.NextUint64()
	i.Flags = b.NextUint64()
}

type FileExtentItem struct {
	// Transaction id that created this extent
	Generation uint64

	// Max number of bytes to hold this extent in ram when we split a
	// compressed extent we can't know how big each of the resulting pieces
	// will be. So, this is an upper limit on the size of the extent in ram
	// instead of an exact limit.
	RamBytes uint64

	// 32 bits for the various ways we might encode the data, including
	// compression and encryption. If any of these are set to something a
	// given disk format doesn't understand it is treated like an incompat
	// flag for reading and writing, but not for stat.
	Compression   uint8
	Encryption    uint8
	OtherEncoding uint16 // For later use

	// Are we inline data or a real extent?
	Type uint8

	// Disk space consumed by the extent, checksum blocks are included in
	// these numbers.

	// At this offset in the structure, the inline extent data start.
	DiskByteNr   uint64
	DiskNumBytes uint64

	// The logical offset in file blocks (no csums) this extent record is
	// for. This allows a file extent to point into the middle of an existing
	// extent on disk, sharing it between two snapshots (useful if some bytes
	// in the middle of the extent have changed.
	Offset uint64

	// The logical number of file blocks (no csums included). This always
	// reflects the size uncompressed and without encoding.
	NumBytes uint64
}

func (i *FileExtentItem) Parse(b *ParseBuffer) {
	i.Generation = b.NextUint64()
	i.RamBytes = b.NextUint64()
	i.Compression = b.NextUint8()
	i.Encryption = b.NextUint8()
	i.OtherEncoding = b.NextUint16()
	i.Type = b.NextUint8()
	// TODO(cblichmann): Inline extents
	i.DiskByteNr = b.NextUint64()
	i.DiskNumBytes = b.NextUint64()
	i.Offset = b.NextUint64()
	i.NumBytes = b.NextUint64()
}

type CSumItem struct {
	CSum uint8
}

func (i *CSumItem) Parse(b *ParseBuffer) {
	i.CSum = b.NextUint8()
	// TODO(cblichmann): Parse the actual checksums
}

type RootItem struct {
	Inode        InodeItem
	Generation   uint64
	RootDirId    uint64
	ByteNr       uint64
	ByteLimit    uint64
	BytesUsed    uint64
	LastSnapshot uint64
	Flags        uint64
	Refs         uint32
	DropProgress Key
	DropLevel    uint8
	Level        uint8

	// The following fields appear after subvol_uuids+subvol_times were
	// introduced.

	// This generation number is used to test if the new fields are valid
	// and up to date while reading the root item. Everytime the root item
	// is written out, the "generation" field is copied into this field. If
	// anyone ever mounted the fs with an older kernel, we will have
	// mismatching generation values here and thus must invalidate the
	// new fields.}
	GenerationV2 uint64
	UUID         uuid.UUID
	ParentUUID   uuid.UUID
	ReceivedUUID uuid.UUID
	CTransId     uint64
	OTransId     uint64
	STransId     uint64
	RTransId     uint64
	Reserved     [8]uint64
}

func (i *RootItem) Parse(b *ParseBuffer) {
	i.Inode.Parse(b)
	i.Generation = b.NextUint64()
	i.RootDirId = b.NextUint64()
	i.ByteNr = b.NextUint64()
	i.ByteLimit = b.NextUint64()
	i.LastSnapshot = b.NextUint64()
	i.Flags = b.NextUint64()
	i.Refs = b.NextUint32()
	i.DropProgress.Parse(b)
	i.DropLevel = b.NextUint8()
	i.Level = b.NextUint8()
	i.GenerationV2 = b.NextUint64()
	if i.Generation == i.GenerationV2 {
		copy(i.UUID[:], b.Next(uuid.UUIDSize))
		copy(i.ParentUUID[:], b.Next(uuid.UUIDSize))
		copy(i.ReceivedUUID[:], b.Next(uuid.UUIDSize))
		i.CTransId = b.NextUint64()
		i.OTransId = b.NextUint64()
		i.STransId = b.NextUint64()
		i.RTransId = b.NextUint64()
		for j, _ := range i.Reserved {
			i.Reserved[j] = b.NextUint64()
		}
	}
}

// This is used for both forward and backward root refs
type RootRef struct {
	DirId    uint64
	Sequence uint64
	NameLen  uint16
	Name     string
}

func (i *RootRef) Parse(b *ParseBuffer) {
	i.DirId = b.NextUint64()
	i.Sequence = b.NextUint64()
	i.NameLen = b.NextUint16()
	l := int(i.NameLen)
	if l > 255 {
		l = 255
	}
	i.Name = string(b.Next(l))
}

// Items in the extent btree are used to record the objectid of the
// owner of the block and the number of references.
type ExtentItem struct {
	Refs       uint64
	Generation uint64
	Flags      uint64
}

func (i *ExtentItem) Parse(b *ParseBuffer) {
	i.Refs = b.NextUint64()
	i.Generation = b.NextUint64()
	i.Flags = b.NextUint64()
}

type Leaf struct {
	Header
	Items []Item
}

func (l *Leaf) Parse(b *ParseBuffer) {
	if l.Header.NrItems == 0 {
		return
	}
	headerEnd := uint32(b.Offset())
	// Clamp maximum number of items to avoid running OOM in case NrItems is
	// corrupted. 0x19 is the typical item size without item data.
	maxItems := b.Unread() / 0x19
	numItems := l.Header.NrItems
	if numItems > uint32(maxItems) {
		numItems = uint32(maxItems)
	}

	l.Items = make([]Item, numItems)
	for i, _ := range l.Items {
		l.Items[i].Parse(b)
	}
	for i, _ := range l.Items {
		item := &l.Items[i]
		b.SetOffset(int(headerEnd + item.Offset))
		item.ParseData(b)
	}
}