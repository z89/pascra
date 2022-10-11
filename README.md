# pascra

a simple pastebin web scraper

**WARNING:** it has only been tested on arch linux, and it may not work on other systems

## Install

```
go get -u github.com/z89/pascra/
```

## Usage

```golang
// import the package
import "github.com/z89/pascra"

func main() {
    // create a new scraper instance
    instance := pascra.New()

    // start the instance
    instance.Start()
}
```

## Main Command

```
Usage:
  pascra [command]

Available Commands:
  fetch       Fetch command to download pastes
  help        Help about any command

Flags:
  -h, --help      help for pascra
  -v, --version   version for pascra
```

## Fetch Command

```
Usage:
  pascra fetch [flags]

Flags:
  -h, --help            help for fetch
      --user string     user to download pastes from (required)
      --dir string      directory to store downloaded pastes
      --pages strings   download specifc pages from user
      --delay int       time delay between downloads to prevent too many requests (milliseconds) (default 250)
      --synchronous     download pastes concurrently for performance

```

an example of a basic command:

![basic command](https://i.imgur.com/T3UR8w1.gif)

an example of an advanced command:

![advanced command](https://i.imgur.com/pwCXizA.gif)

an example of a synchronous command (disables concurrency):

![advanced command](https://i.imgur.com/nhYhV6H.gif)

## Contributing

pull requests are welcome to:

- revise the docs
- fix bugs or bad logic
- suggest/add features or improvements

## License

MIT
