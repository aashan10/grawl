## GRAWL

Grawl is a sitemap crawler tool written in golang. The primary purpose of grawl is to crawl your sitemaps and check the health of the pages enlisted in the sitemaps 


### Installation

You'll need golang to install this project. You can get golang from the [official downloads page](https://go.dev/dl/)

With golang installed and set up in your You can install grawl with `go install` with the following command 
```bash 
go install github.com/aashan10/grawl@latest
```

Once this command runs, you'll have `grawl` binary in your `$GOPATH/bin` directory (usually `$HOME/go/bin`)

If you want this binary globally, make sure to include `export PATH="$PATH:$HOME/go/bin"` in your shell configuration file
> You can find your `GOPATH` by running `go env` command. You can follow the [official docs](https://go.dev/wiki/SettingGOPATH) to set up the proper environment according to your operating system.


### Usage 

Use `grawl help` command to list out all the available options.

### Crawling Sitemap 

Use `grawl scan <your-sitemap-root-url>` to start scanning the sitemaps. 

> If you are using https locally and do not have a valid certificate, add `-k` flag with `grawl scan` to skip TLS certificate check. For more info, run `grawl scan -h`
