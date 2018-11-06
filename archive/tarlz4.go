package archive

import (
	"fmt"
	"io"
	"strings"

	"github.com/pierrec/lz4"
)

// TarLz4 facilitates lz4 compression
// (https://github.com/lz4/lz4/tree/master/doc)
// of tarball archives.
type TarLz4 struct {
	*Tar

	// The compression level to use when writing.
	// Minimum 0 (fast compression), maximum 12
	// (most space savings).
	CompressionLevel int
}

// Archive creates a compressed tar file at destination
// containing the files listed in sources. The destination
// must end with ".tar.lz4" or ".tlz4". File paths can be
// those of regular files or directories; directories will
// be recursively added.
func (tlz4 *TarLz4) Archive(sources []string, destination string) error {
	if !strings.HasSuffix(destination, ".tar.lz4") &&
		!strings.HasSuffix(destination, ".tlz4") {
		return fmt.Errorf("output filename must have .tar.lz4 or .tlz4 extension")
	}
	tlz4.wrapWriter()
	return tlz4.Tar.Archive(sources, destination)
}

// Unarchive unpacks the compressed tarball at
// source to destination. Destination will be
// treated as a folder name.
func (tlz4 *TarLz4) Unarchive(source, destination string) error {
	tlz4.wrapReader()
	return tlz4.Tar.Unarchive(source, destination)
}

// Walk calls walkFn for each visited item in archive.
func (tlz4 *TarLz4) Walk(archive string, walkFn WalkFunc) error {
	tlz4.wrapReader()
	return tlz4.Tar.Walk(archive, walkFn)
}

// Create opens tlz4 for writing a compressed
// tar archive to out.
func (tlz4 *TarLz4) Create(out io.Writer) error {
	tlz4.wrapWriter()
	return tlz4.Create(out)
}

// Open opens t for reading a compressed archive from
// in. The size parameter is not used.
func (tlz4 *TarLz4) Open(in io.Reader, size int64) error {
	tlz4.wrapReader()
	return tlz4.Tar.Open(in, size)
}

// Extract extracts a single file from the tar archive.
// If the target is a directory, the entire folder will
// be extracted into destination.
func (tlz4 *TarLz4) Extract(source, target, destination string) error {
	tlz4.wrapReader()
	return tlz4.Tar.Extract(source, target, destination)
}

func (tlz4 *TarLz4) wrapWriter() {
	var lz4w *lz4.Writer
	tlz4.Tar.writerWrapFn = func(w io.Writer) (io.Writer, error) {
		lz4w = lz4.NewWriter(w)
		lz4w.Header.CompressionLevel = tlz4.CompressionLevel
		return lz4w, nil
	}
	tlz4.Tar.cleanupWrapFn = func() {
		lz4w.Close()
	}
}

func (tlz4 *TarLz4) wrapReader() {
	tlz4.Tar.readerWrapFn = func(r io.Reader) (io.Reader, error) {
		return lz4.NewReader(r), nil
	}
}

func (tlz4 *TarLz4) String() string { return "tar.lz4" }

// Compile-time checks to ensure type implements desired interfaces.
var (
	_ = Reader(new(TarLz4))
	_ = Writer(new(TarLz4))
	_ = Archiver(new(TarLz4))
	_ = Unarchiver(new(TarLz4))
	_ = Walker(new(TarLz4))
	_ = Extractor(new(TarLz4))
)

// DefaultTarLz4 is a convenient archiver ready to use.
var DefaultTarLz4 = &TarLz4{
	CompressionLevel: 9, // https://github.com/lz4/lz4/blob/1b819bfd633ae285df2dfe1b0589e1ec064f2873/lib/lz4hc.h#L48
	Tar:              DefaultTar,
}
