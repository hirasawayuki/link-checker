# link-checker CLI
link-cheker is a tool to check for broken links.
By specifying the URL of the page you want to check, you can display a list of broken URLs.

## usage

```shell
$ go run main.go -u={check page URL}
```

By default, links with an HTTP status of 400 or higher will be displayed.
If you want to display HTTP status 200 links as well, add `-a` option.

```shell
$ go run main.go -u='{check page URL}' -a
```

`-t` option can be set to set the request interval.
Default value is 100(ms).
For example, if you want to make two requests in one second, you can do the following:

```shell
$ go run main.go -u='{check page URL}' -t 500
```
