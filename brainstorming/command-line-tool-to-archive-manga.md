# Command-Line Tool to Archive Manga

Reading manga in VR requires packaging image files into a CBZ (Comic Book ZIP) or PDF, since many VR reading apps accept only those formats. A CBZ file is essentially just a ZIP archive of images (pages) with a .cbz extension.

Goal: Develop a command-line tool (for Linux/WSL2, but portable to other OSes) that takes a parent directory containing manga chapters (each chapter in its own folder of images) and produces a CBZ file for each chapter. This will automate the process so that you can read the chapters in VR without manually zipping each one. I don't know yet what programming language this repository will use to handle this task, but we know that the program should be a simple executable file.

Basic Operation: The tool will accept one argument: the path to a directory (let's call it MangaDir) which contains multiple subfolders, each representing a chapter. By default, the tool will iterate over each immediate subfolder in MangaDir and create a CBZ archive for that subfolder’s contents. For example, if MangaDir has subfolders Chapter 1, Chapter 2, ... each containing .jpg/.png files, the tool will output Chapter 1.cbz, Chapter 2.cbz, etc., in the specified output location.

## Key Features and Behavior:

### Input structure
The program assumes the input directory is organized with one level of subdirectories for chapters (e.g. MangaDir/Chapter 1/, MangaDir/Chapter 2/, ...). It will treat each subfolder as a separate chapter to package. (If needed, a recursive mode could be an option for nested structures, but by default one level is expected as per your use case.)
### Output location and naming
By default, the CBZ files could be created in the parent directory (MangaDir) or in a specified output directory. Each CBZ will be named after its source folder. For instance, Chapter 35-1/ becomes Chapter 35-1.cbz. The tool should handle spaces or special characters in folder names by either properly quoting/escaping or replacing them (so that the resulting file is valid). The output naming convention is an invariant: it should be consistent and predictable (exact folder name plus “.cbz” extension, unless a naming template is provided).
### Image file handling
The tool will include all image files in each chapter folder in the CBZ. Common image formats (PNG, JPG, JPEG, possibly BMP or GIF) should be recognized and included. Non-image files (like thumbnails, .txt notes, etc.) can be ignored or optionally excluded by default to keep the CBZ clean. The program can have a filter for file extensions (e.g., include .png/.jpg by default, maybe configurable).
### Ordering of pages
It must preserve the reading order of images. Typically, manga scanlation files are named with zero-padded numbers or in alphabetical order (e.g. 01.jpg, 02.jpg, ... 10.jpg etc.) so that a simple alphabetical sort yields the correct page sequence. The tool should sort image filenames in natural/alphanumeric order before zipping to ensure page 10 doesn’t come before page 2, for example. In practice, if files are already zero-padded (which is common), a normal lexicographical sort is sufficient. But as an invariant, the program should guarantee that within each CBZ, the images are ordered as they were in the folder (assuming the folder’s naming reflects the intended order). This may involve implementing a numeric-aware sort (treating "9" < "10" properly) if filenames are not uniformly padded.
### CBZ creation
The actual creation of the CBZ can be done by compressing the images into a ZIP archive. The program can use system zip under the hood or a zip library (like Python’s zipfile module or a Go/Rust library depending on implementation language). The crucial part is to compress the image files (preferably without any additional directory path inside the archive). The images in the archive should appear at the root of the CBZ, not nested in an extra folder, so that comic readers can display them without issues. (Most CBZ readers can handle a root folder inside the zip, but keeping images at root is a cleaner standard.)
### Performance considerations
CBZ creation can involve large image files, so the tool should stream or efficiently handle files rather than loading all images into memory at once. Since you’re on WSL2 (Ubuntu), the tool should be optimized for that environment (e.g. avoid unnecessary Windows path issues). If implementing in a scripting language (Python, Bash), ensure it can handle large directories. If using a compiled language (Go/C++), memory management should be considered for large images. Multi-threading is not critical unless dealing with many large chapters, but it could zip chapters one by one (or even in parallel if needed) to speed up a big batch.
### Error handling
The tool should handle edge cases gracefully. For example, if a chapter folder is empty or missing images, it could skip creating a CBZ for it (and possibly warn the user). If a CBZ file with the same name already exists in the output location, the tool could either overwrite it by default or require a --force flag to overwrite (to avoid unintentional overwrites). It should report any problems (like if no image files are found in a supposed chapter folder, or if it encounters a read/write error) rather than silently failing.
### User interface
Being a command-line tool, usage might look like:

manga2cbz /path/to/MangaDir --out /path/to/output --recursive --format jpg,png

Where --out is optional (defaults to the input directory if not provided), --recursive could tell it to scan all subfolders at any depth for image folders, and --format could allow specifying which image extensions to include. In the simplest usage, just:

manga2cbz /path/to/MangaDir

would create CBZs for each subfolder in MangaDir and place them alongside those folders.

(The actual command name and options can be decided as needed; the above is just an example of how it might operate.)

With these specifications in mind, we can define the key invariants and then outline tests to ensure the tool works correctly.

## Key Invariants for the Tool

The following invariants are conditions that must always hold true for the program to be considered correct. These are essentially the “must-pass” criteria for the tool:

### Complete Coverage
Every image file in each chapter folder must be included exactly once in the corresponding CBZ. No images should be missing from the archive, and no extra files should be added. The number of image files in a CBZ should match the number of image files in the source folder (this can be verified after creation).
### Correct Ordering
Images in the CBZ archive must be in the intended reading order (usually the alphabetical/numerical order of filenames). The tool should enforce this by sorting file names appropriately before archiving. An invariant check could be that for any two images A and B in a folder, if A comes before B in the natural sorted order, then A’s entry in the CBZ comes before B’s entry.
### Proper Naming
The output CBZ file name should correctly reflect the chapter folder name (to avoid confusion). For example, a folder named “Chapter 5” yields “Chapter 5.cbz”. This invariant ensures one-to-one mapping between folder and archive. If any sanitization of the name is needed (e.g., removing illegal filename characters on Windows), the mapping should remain obvious (perhaps documented or logged).
### Valid Archive Format
The .cbz files produced must be valid ZIP archives readable by comic reader apps. This means the archive isn’t corrupt and follows the zip format standards. A technical invariant is that after creation, a test like unzip -t <file>.cbz should report no errors (indicating a valid zip). Similarly, the VR app or any CBZ reader should be able to open it without issues. Essentially, if an archive is created, it must be a well-formed CBZ/ZIP (zero-length or corrupt files are not acceptable).
### No Internal Folder Structure
Unless intentionally designed otherwise, the CBZ should not contain extraneous directory hierarchy. The images should be at the root of the archive. Invariant-wise, when listing archive contents, the file paths should typically be just page1.jpg, page2.jpg, ... rather than e.g. Chapter5/page1.jpg. (This avoids potential issues with readers expecting files at root. The tool can enforce this by how it adds files to the zip – e.g., by changing into the directory before zipping or using a flattening option.)
### Idempotence (Optional):
Running the tool twice on the same input should produce the same result (assuming no changes in input in between). This means if CBZ files already exist and no content changed, either the tool will skip re-creating them or create identical archives. This is more of a quality guarantee: if nothing changed in source folders, the output archives should not mysteriously differ. (If the tool includes timestamps inside ZIPs by default, archives might differ binary-wise, but functionally they should be equivalent. For testing, one might disable timestamp to compare outputs exactly.)
### Robustness
The tool should handle unusual but possible scenarios. For instance, if a chapter folder contains a mix of uppercase .JPG and lowercase .jpg extensions, or has files with spaces, the invariant is that all valid image files still get included and ordering is correct. Non-image files (if any) should either be excluded or handled in a defined way (e.g., an invariant could be that only files with approved extensions end up in the CBZ, ensuring no stray files like .txt or .db are inside).
These invariants ensure that the CBZ creation tool reliably produces correct and usable archives for each manga chapter, without data loss or ordering mistakes.

## Testing Plan
To validate the tool against the specifications and invariants, we should design a series of tests. Each test will provide a certain input scenario and then verify the outcomes (the CBZ files) against expected results. Below are several important test cases:
### Basic Multiple Chapters Test
Prepare a TestManga directory with two subfolders: ChapterA and ChapterB. Put a few image files in each (e.g., 01.jpg, 02.jpg, 03.jpg in ChapterA, and 01.png, 02.png in ChapterB). Run the tool on TestManga. It should produce ChapterA.cbz and ChapterB.cbz. Verify that:

Both CBZ files are created in the output location.

Opening ChapterA.cbz reveals exactly the three images from ChapterA in the correct order (01.jpg, 02.jpg, 03.jpg).

Similarly, ChapterB.cbz contains the two .png files in order.

Check with a ZIP tester or by opening in a reader that the archives are not corrupted (e.g., unzip -t ChapterA.cbz returns OK).
### Ordering/Sorting Test
Create a folder ChapterC with images named non-uniformly, e.g. 1.png, 2.png, 10.png, 11.png, 3.png. Run the tool on a parent directory containing ChapterC. The resulting ChapterC.cbz should have the images ordered as 1, 2, 3, 10, 11 (assuming natural numeric order) rather than the pure alphabetical (which would sort "10.png" after "1.png" but before "2.png" if not handled). If the tool is intended to handle numeric sequences, verify that it does so. If the tool relies purely on alphabetical and expects zero-padded naming, then adjust the test to use names like 01, 02, ... 10 and ensure that ordering is alphabetical which coincides with numeric. The main point is to ensure the output ordering matches expected reading order.
### Mixed File Types and Filtering Test
In a test chapter folder ChapterD, include various file types: e.g., 01.jpg, 02.png, 03.gif as images, and also include a non-image file like notes.txt or a hidden file. Run the tool. The resulting ChapterD.cbz should contain only the image files (jpg, png, gif) and not the notes.txt (unless the tool is designed to include all files by default). Verify the archive contents exclude the unwanted file. This test confirms that the file type filtering works and that the tool handles multiple image formats together. All three images should be present and in correct order (respecting how they were named).
### Empty or Missing Chapter Test
Create an empty folder ChapterEmpty (no image files inside) under the test directory. Run the tool. It should detect that there are no images. Depending on design, it might skip creating a CBZ altogether for an empty folder, or create an empty CBZ (less useful). The preferable behavior would be to skip and possibly warn the user. The test passes if the program either outputs a warning like "No images found in ChapterEmpty, skipping" and does not create a CBZ, or it creates an empty CBZ but also warns (and the empty CBZ should at least be a valid zip file with no entries). We verify that it didn’t crash on this edge case and handled it as expected.
### Output Overwrite Test
Place a few images in ChapterX and run the tool to create ChapterX.cbz. Then run the tool again on the same input directory without changing anything. Since ChapterX.cbz already exists from the first run, the tool should handle this gracefully. If the specification is to overwrite by default, then after the second run ChapterX.cbz still exists (possibly updated with the same content). If the specification is to avoid overwriting unless a flag is given, then the tool might skip creating ChapterX.cbz the second time and log a message that it already exists. Test both scenarios: by default it might overwrite or skip (check which one is implemented), and with an explicit overwrite flag (if available) ensure that it does replace the file. The invariant here is that no unintended duplicate or corrupt file results from re-processing the same input.
### Large Batch/Performance Test (optional)
For a more intensive test, use a directory with, say, 50 chapter subfolders, each containing 100+ images (you can use copies of a few images to simulate). Run the tool on the entire set. Verify that all 50 CBZ files are created correctly. This test checks that the tool can scale to a real-world manga volume or series. Important things to observe: it doesn’t run out of memory or crash, it completes in a reasonable time, and all outputs pass integrity checks. After creation, pick a couple of archives (e.g., first, middle, last) and verify their contents manually or with a script (count files, compare filenames to source). This gives confidence that the tool can handle a large library conversion in bulk.
Each of these tests targets specific requirements and invariants of the program. By executing them, you can catch issues early (like ordering bugs, skipping files, or format problems) and ensure that the manga-to-CBZ converter works reliably. Once all tests pass, you’ll have a robust tool to batch-convert your manga chapters into CBZ files, ready for enjoyment in VR.
