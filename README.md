<br>
<!-- logo --> <div align="center">
<a href="https://github.com/z89/"><img width="600px" height="150px" src="https://i.imgur.com/r6uqOq5.png" alt="Logo"></a></div>

<div align="center">

<img src="https://img.shields.io/badge/CODE%20QUALITY-D---blueviolet?style=for-the-badge&logo=codacy"></img>
<img src="https://img.shields.io/badge/archlinux-supported-blue?style=for-the-badge"></img>
<img src="https://img.shields.io/badge/windows-not%20supported-red?style=for-the-badge"></img>
</div>
<br>
<p align="center">A pastebin.com web crawler to download <strong>ANY</strong> user pastes written in Go using <a href="https://github.com/gocolly/colly">colly</a> and <a href="https://www.gnu.org/software/wget/">wget</a></p>

<br>
<strong><p  align="center">z89 (Author): This is a beginner Go project for me to get introduced to learning Go. However the base functionality seems to be working fine. The code is horrible, but understand it was written in under 48 hours by someone who had never touched Go before! Enjoy :)</p></strong>
<br>

## Installation

<h3><a href="https://github.com/z89/pascra/releases">Binaries</a></h3>
A precompiled binary exists already in the repo, click the binaries link above to get the lastest version. The only dependency is the <strong>wget</strong> GNU package. Below is an example of how to download <strong>wget</strong> on Arch Linux & Ubuntu from their repo's:

#### Arch Linux
```sh
sudo pacman -S wget
```

#### Ubuntu
```sh
sudo apt install wget
```


You can also download this repo to compile the Go source code manually. This requires Go to be installed locally of course so make sure you already have that working.


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
 Pascra v1.1.0-release written by z89 

    -h, --help              show these help instructions

    -v, --version           display the program version

    -d, --dir              specify the directory for the downloaded pastes

    -u, --user              select the user from pastebin.com to download from
```



## License

Fuck licenses, Open source ftw #HTP :)
