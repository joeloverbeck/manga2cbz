# manga2cbz

Convert manga chapter folders into CBZ archives for VR reading applications.

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white)

## Overview

**manga2cbz** is a command-line tool that converts directories containing manga chapters into CBZ (Comic Book ZIP) archives. It's designed for preparing manga collections for VR readers and other comic book applications.

**Key features:**

- **Automatic WebP conversion** - Converts WebP images to PNG for VR reader compatibility
- **Natural sort ordering** - Pages sort correctly: `page2` comes before `page10`
- **Flat and recursive discovery** - Process simple or nested chapter structures
- **Memory efficient** - Streams files directly without loading entire images into memory
- **Optimized CBZ output** - Uses ZIP Store method (no compression overhead for already-compressed images)

## Installation

### Prerequisites

- Go 1.21 or later

### Build from Source

```bash
git clone https://github.com/yourusername/manga2cbz.git
cd manga2cbz
go build -o manga2cbz ./cmd/manga2cbz
```

### Production Build (Optimized)

```bash
go build -trimpath -ldflags "-s -w" -o manga2cbz ./cmd/manga2cbz
```

### Cross-Compile for Windows

```bash
GOOS=windows GOARCH=amd64 go build -o manga2cbz.exe ./cmd/manga2cbz
```

## Quick Start

```bash
# Convert all chapters in a manga folder
./manga2cbz /path/to/manga

# Process nested volumes/chapters with verbose output
./manga2cbz -r -v /path/to/manga

# Output to a different directory
./manga2cbz -o /output/folder /path/to/manga
```

## Usage

```
manga2cbz [OPTIONS] <INPUT_DIR>

Arguments:
  <INPUT_DIR>    Directory containing chapter folders

Options:
  -o, --out <DIR>     Output directory (default: input directory)
  -e, --ext <EXTS>    Comma-separated image extensions (default: jpg,jpeg,png,gif,bmp,webp)
  -f, --force         Overwrite existing CBZ files
  -r, --recursive     Process nested directory structures
  -n, --no-convert    Disable WebP to PNG conversion (WebP converted by default)
  -v, --verbose       Show detailed progress
  -q, --quiet         Suppress non-error output (-q takes precedence over -v)
  -h, --help          Show help message
      --version       Show version
```

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | All chapters processed successfully |
| 1 | Partial failure (some chapters failed) |
| 2 | Total failure (invalid arguments or no chapters found) |

## Examples

### Flat Mode (Default)

Process direct subdirectories as chapters:

```bash
./manga2cbz /manga/OnePiece
```

**Input structure:**
```
/manga/OnePiece/
  Chapter001/
    page01.jpg
    page02.jpg
  Chapter002/
    page01.jpg
    page02.jpg
```

**Output:**
```
/manga/OnePiece/
  Chapter001.cbz
  Chapter002.cbz
```

### Recursive Mode

Process nested directory structures (volumes containing chapters):

```bash
./manga2cbz -r /manga/OnePiece
```

**Input structure:**
```
/manga/OnePiece/
  Volume01/
    Chapter001/
      page01.jpg
    Chapter002/
      page01.jpg
  Volume02/
    Chapter003/
      page01.jpg
```

**Output:**
```
/manga/OnePiece/
  Volume01_Chapter001.cbz
  Volume01_Chapter002.cbz
  Volume02_Chapter003.cbz
```

### Custom Output Directory

```bash
./manga2cbz -o /output/cbz /manga/OnePiece
```

### Custom Image Extensions

Process only PNG and WebP files:

```bash
./manga2cbz -e png,webp /manga/OnePiece
```

### Force Overwrite

Regenerate existing CBZ files:

```bash
./manga2cbz -f /manga/OnePiece
```

### Verbose Progress

See detailed processing information:

```bash
./manga2cbz -v /manga/OnePiece
```

**Output:**
```
Found 3 chapter(s)
Processing: Chapter001
Created: /manga/OnePiece/Chapter001.cbz (24 images)
Processing: Chapter002
Created: /manga/OnePiece/Chapter002.cbz (22 images)
Processing: Chapter003
Created: /manga/OnePiece/Chapter003.cbz (20 images)
```

## Directory Structure

### Input Requirements

- Input must be a directory containing chapter subdirectories
- Each chapter directory should contain image files
- Images are identified by file extension (case-insensitive)

### Output Format

- One CBZ file per chapter
- CBZ filename matches chapter directory name
- In recursive mode, nested paths use underscores: `Volume01/Chapter001` becomes `Volume01_Chapter001.cbz`
- Images stored at the archive root (no nested folders inside CBZ)

## WebP Compatibility

Many VR reader applications don't support WebP images in CBZ files, resulting in black pages. To ensure maximum compatibility, manga2cbz **automatically converts WebP images to PNG format** during CBZ creation.

### Why PNG?

| Format | Pros | Cons |
|--------|------|------|
| **PNG** | Lossless, preserves quality, supports transparency | Larger file size |
| **JPEG** | Smaller file size | Lossy compression, no transparency |

PNG is chosen because manga artwork benefits from lossless compression, and transparency support is important for some manga pages.

### Disabling Conversion

If you want to keep WebP files as-is (e.g., your reader supports WebP), use the `--no-convert` or `-n` flag:

```bash
# Disable WebP conversion
./manga2cbz --no-convert /path/to/manga

# Short form
./manga2cbz -n /path/to/manga
```

## Technical Details

### Natural Sort

Standard alphabetical sorting produces incorrect page order:

```
page1.jpg, page10.jpg, page2.jpg  (WRONG)
```

manga2cbz uses natural sort for correct ordering:

```
page1.jpg, page2.jpg, page10.jpg  (CORRECT)
```

### CBZ Format

CBZ files are ZIP archives with:

- **Store method** (no compression) - Images are already compressed, so additional compression wastes CPU without reducing size
- **Flat structure** - All images at archive root for maximum compatibility

### Memory Efficiency

Files are streamed directly from disk to the archive using `io.Copy`, avoiding loading entire images into memory. This allows processing chapters with many large images without excessive memory usage.

## License

MIT License - See LICENSE file for details.
