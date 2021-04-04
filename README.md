# Rename Machine

*This program is designed to rename/cleaning filenames, adding titles (from a list) and provide specific tag insertion and moving files contained in multi directories to a single folder ... lot of options available (regex, Posix character classes, case sensitive, keep between, extract, titling).*

#### Last update 2021-04-02

Take a look [here, H.F.M repositories](https://github.com/hfmrow/) for other useful linux softwares.

- If you just want to use it, simply download the '*.deb' compiled version under the [Releases](https://github.com/hfmrow/rename-machine/releases) tab.

- If you want to play inside code, see below "How to compile" section.

## How it's made

- Programmed with go language: [golang](https://golang.org/doc/) 
- GUI provided by [Gotk3 (gtk3 v3.22)](https://github.com/gotk3/gotk3), GUI library for Go (minimum required v3.16).
- I use home-made software: "Gotk3ObjHandler" to embed images/icons, UI-information and manage/generate gtk3 objects code from [glade ui designer](https://glade.gnome.org/). and "Gotk3ObjTranslate" to generate the language files and the assignment of a tooltip on the gtk3 objects (both are not published at the moment, in fact, they need documentations and, for the moment, I have not had the time to do them).

## Functionalities

- Single or multi-files interface.
- Rename single or multiple files at once.
- Remove parts of filename (3 patterns at once).
- Replace parts of filename (3 patterns at once).
- Keep between patterns possibility (Case sensitive, posix [Character classes](https://www.regular-expressions.info/posixbrackets.html)).
- Pattern subtracting possibility (Case sensitive, posix [Character classes](https://www.regular-expressions.info/posixbrackets.html)).
- Ability to use regex for removing patterns.
- Drag and drop capacity available.
- Add at begin or at end an incremental number.
- Mask for filtering by extensions.
- Function to preserve extensions.
- Add titles list to filenames (from an existing list, useful for series, music albums)
- Move files contained in multiple folders and sub-folders to a single directory in one click.
- Each function have his tooltip for explanations.

## Some pictures and explanations

**Single entry window.**  
![](assets/readme/single-entry.jpg) 

**This is the main screen.**  
![](assets/readme/mainScr-rename-engine.jpg)  

**Remove & Replace multiples patterns.**  
![](assets/readme/ren-repl.jpg)  

**Add incremental numbers to file names.**  
![](assets/readme/inc.jpg)  

**The keep between window**  
![](assets/readme/keep-btw1.jpg)  

**Keep between.**  
*Applied using posix [character classes](https://www.regular-expressions.info/posixbrackets.html)  option. By this way, numeric values "00,01,02,03,..." will always matching for replacement or removing.*  
![](assets/readme/keep-btw2.jpg)  

**Using regular expression to replace pattern.**  
*In this case, (\i) means to be case insensitive context, between parentheses, there is the pattern to find (upper and lower cases), the dot "." means any single character, and "[[:digit:]]" for any single number "0-9" in posix [character classes](https://www.regular-expressions.info/posixbrackets.html) notation.*  
![](assets/readme/regex-1.jpg)  

 **Subtraction pattern using posix**[ Character classes](https://www.regular-expressions.info/posixbrackets.html).  
*Useful for series or audio files. The (**1** entry) will be internally transformed into ([[:upper:]][[:upper:]][[:upper:]][[:upper:]][[:space:]][[:digit:]][[:digit:]]) full (**2** character class compliant), using (**3** strict mode)*  
![Substract pattern](assets/readme/substract-1.jpg "Substract pattern")  

 **Adding Titles to filenames.**  
 *On upper left box, you can past titles list (from wikipedia for example), At bottom center you have your files, at top right box you got the result. As you can see, before each title you have “**1-**” a number and a dash, you only want to keep the line after the dash. So, you define “Separator” as **-**, the “Field” at **1** and a simple space “Before title”. I have added more entry to see some possibilities.*  
![Adding Titles](assets/readme/title-example.jpg "Adding Titles")  

## How to compile

- 

- Open terminal window and at command prompt, type: `go get github.com/hfmrow/rename-machine`

- See [Gotk3 Installation instructions](https://github.com/gotk3/gotk3/wiki#installation) for gui installation instruction.

- To change gtk3 interface you need to use a home made software, (not published actually). So don't change gtk3 interface (glade file) ...

- To change language file you need to use another home made software, (not published actually). So don't change language file ...

- To Produce a stand-alone executable, you must change inside "main.go" file:
  
  ```go
    func main() {
        devMode = true
    ...
  ```
  
  into
  
  ```go
    func main() {
        devMode = false
    ...
  ```

This operation indicate that externals data (Image/Icons) must be embedded into the executable file.

### Os informations (build with)

| Name                                                       | Version / Info / Name                          |
| ---------------------------------------------------------- | ---------------------------------------------- |
| GOLANG                                                     | V1.16.2 -> GO111MODULE="off", GOPROXY="direct" |
| DISTRIB                                                    | LinuxMint Xfce                                 |
| VERSION                                                    | 20                                             |
| CODENAME                                                   | ulyana                                         |
| RELEASE                                                    | #46-Ubuntu SMP Fri Jul 10 00:24:02 UTC 2020    |
| UBUNTU_CODENAME                                            | focal                                          |
| KERNEL                                                     | 5.8.0-48-generic                               |
| HDWPLATFORM                                                | x86_64                                         |
| GTK+ 3                                                     | 3.24.20                                        |
| GLIB 2                                                     | 2.64.3                                         |
| CAIRO                                                      | 1.16.0                                         |
| [GtkSourceView](https://github.com/hfmrow/gotk3_gtksource) | 4.6.0                                          |
| [LiteIDE](https://github.com/visualfc/liteide)             | 37.4 qt5.x                                     |
| Qt5                                                        | 5.12.8 in /usr/lib/x86_64-linux-gnu            |

- The compilation have not been tested under Windows or Mac OS, but all file access functions, line-end manipulations or charset implementation are made with OS portability in mind.

## You got an issue ?

- Give information (as above), about used platform and OS version.
- Provide a method to reproduce the problem.
