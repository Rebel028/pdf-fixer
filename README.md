## Simple tool to fix all the files made with Quartz PDFContext on MacOS

Based on https://github.com/pdfcpu/pdfcpu

Recursively walks all files in current directory, rewriting them in PDF/A-compliant format.

## How to

1. Download `pdf-fixer.exe` from the [latest release](https://github.com/Rebel028/pdf-fixer/releases/latest)
2. Make sure you backed up your pdf files (just in case)
3. Put the executable in directory that contains files you need to fix (files may be located in separate folders in this dir)
4. Launch `pdf-fixer.exe`

Tool will enumerate all .pdf files in current directory and its subfolders. If any of these files were created with Quartz PDFContext it will fix and overwrite them.

## Build from source

**Prerequisites:** Go 1.21

```shell
git clone https://github.com/Rebel028/pdf-fixer.git
cd pdf-fixer
go build -o pdf-fixer.exe
```
