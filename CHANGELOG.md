# Changelog

## Rename Machine

All notable changes to this project will be documented in this file.

## [1.6.1] 2021-04-02

### Added

- The title function, the spin-button cannot be out of range if the number of separators selected is greater than the actual number of separators in the sentence (title).

### Fixed

- Correct handling of error when a file/dir passed by argument does not exist.

- Title function now concatenate the rest of the splitted sentence (title) after the selected separator encountered, (useful when separator is a space and the title contain some other spaces).

### Changed

- Character classes (posix) [[:blank:]] replaced by [[:space:]] and added to non strict mode. That means the using of strict-mode for matching patterns containing spaces is no longer needed.

- Internal changes, full code re-writing for treeviews and images display.

- Repository name was changed to [https://github/hfmrow/rename-machine](https://github.com/hfmrow/search-engine) instead of `https://github.com/hfmrow/renMachine`

---

## [1.5] 2019-07-28

### Added

- Add Cumulative drag and drop checkbox.

### Fixed

- Now, keep the options for each tab when changing them.
- Fixing resizing issue on changing tab number.
- The last window size is correctly restored at application start.
- Fixed some issues on scaling up main window about some controls that have  an unwanted resizing effect.
- Move function now properly handle duplicate and/or existing filenames.

### Changed

- Removed "Reset to original files list" button, not needed anymore since options will correctly saved on changing to other tab.
- Some codes was rewrote for better results when managing files passed by command line or by drag & drop.
- The functions of genLib (my private lib) are now included, in order to get a more lightweight final executable.
- Sort function removed to avoid confusion on titleing & renaming.