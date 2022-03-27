# link-checker CLI
link-cheker is a tool to check for broken links.
By specifying the URL of the page you want to check, you can display a list of broken URLs.

<img width="720" alt="スクリーンショット 2022-02-06 22 14 59" src="https://user-images.githubusercontent.com/48427044/152682723-63413a89-664f-481a-878e-ca035d9e88a4.png">



## usage

```shell
// install
$ go install github.com/hirasawayuki/link-checker/cmd/link-checker@v1.0.0
```

```shell
$ link-checker -u='{check page URL}'
```

By default, links with an HTTP status of 400 or higher will be displayed.
If you want to display HTTP status 200 links as well, add `-a` option.

```shell
$ link-checker -u='{check page URL}' -a
```

`-t` option can be set to set the request interval.
Default value is 100(ms).
For example, if you want to make two requests in one second, you can do the following:

```shell
$ link-checker -u='{check page URL}' -t 500
```
