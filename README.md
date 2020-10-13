<br>
<!-- logo --> <div align="center"><a href="https://github.com/z89/"><img width="600px" height="150px" src="/logo.png" alt="Logo"></a></div>
<br>
<div align="center">
<img src="https://img.shields.io/badge/go-v1.15.2-blue?style=for-the-badge&logo=go"></img>
<img src="https://img.shields.io/badge/wget-v1.20.3-green?style=for-the-badge"></img>
<img src="https://img.shields.io/badge/COLLY-v2.1.0-orange?style=for-the-badge"></img>
<br>
<img src="https://img.shields.io/badge/CODE%20QUALITY-D---blueviolet?style=for-the-badge&logo=codacy"></img>
<img src="https://img.shields.io/badge/linux-supported-blue?style=for-the-badge"></img>
<img src="https://img.shields.io/badge/windows-not%20supported-red?style=for-the-badge"></img>
</div>
<br>
<p align="center">A pastebin.com web crawler to download <strong>ANY</strong> user pastes written in Go using <a href="https://github.com/gocolly/colly">colly</a> and <a href="https://www.gnu.org/software/wget/">wget</a></p>

<br>
<strong><p  align="center">z89 (Author): This is a beginner Go project for me to get introduced to  learning Go. Lots of optimizations are still yet to be added, this is an early release for this production build! However the base functionality seems to be working fine. The code is horrible, but understand it was written in under 48 hours by someone who had never touched Go before! Enjoy :) </p></strong>

<br>

## Installation

<h3><a href="https://github.com/z89/pascra/releases">Binaries</a></h3>
A precompiled binary exists already in the repo, click the binaries link above to get the lastest version. The only dependency is the <strong>wget</strong> GNU package. Below is an example of how to download <strong>wget</strong> on Arch Linux & Ubuntu:

#### Arch Linux
```sh
sudo pacman -S wget
```

#### Ubuntu
```sh
sudo apt install wget
```


You can also manually download this repo compile the Go source code. This requires Go to be installed to be able to build the package however.


```sh
git clone https://github.com/z89/pascra
```
```sh
cd pascra/
```
```sh
go build -o pascra
```

## Usage
```
 Pascra v1.0-alpha written by z89 

    -h, --help              show this help file

    -v, --version           do not display directory overwrite warning

    -d, --dir               specify the directory for the downloaded pastes

    -u, --user              select the user from pastebin.com to download from

    -q, --quiet             do not display directory overwrite warning and show limited features
```



## License

Fuck licenses, use it how you wish! Open source ftw #HTP
