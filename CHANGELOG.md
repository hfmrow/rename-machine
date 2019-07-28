# Changelog

## Rename Machine

All notable changes to this project will be documented in this file.

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
