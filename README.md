# goddns

Dynamic DNS Updater for Google Domains

## Install

```sh
go get -u github.com/delphinus/goddns
```

For macOS, you can use [Homebrew][] tap.

[Homebrew]: https://brew.sh

```sh
brew tap delphinus/goddns
brew install goddns
```

## Usage

Create config in `/usr/local/etc/goddns.toml`. Examples below.

TODO: Enable to specify the config filename.

```toml
[[domains]]
username = 'hogehoge'
password = 'fugafuga'
hostname = 'hoge.example.com'

[[domains]]
username = 'hogefuga'
password = 'fugahoge'
hostname = 'fuga.example.com'
```

And run.

```sh
goddns
```

That's all! `goddns` dumps log in syslog.

```
Mar 10 20:07:17 some-host goddns[62197]: loading /usr/local/etc/goddns.toml
Mar 10 20:07:17 some-host goddns[62197]: starting: hoge.example.com
Mar 10 20:07:18 some-host goddns[62197]: detected: 123.123.234.234
Mar 10 20:07:22 some-host goddns[62197]: result: Successful! code: good, ip: 123.123.234.234
Mar 10 20:07:22 some-host goddns[62197]: starting: hoge.example.com
Mar 10 20:07:23 some-host goddns[62197]: detected: 234.234.123.123
Mar 10 20:07:27 some-host goddns[62197]: result: Successful! code: good, ip: 234.234.123.123
```

`goddns` will check IP Address and update if needed. The process is automatically repeated every 60 seconds.

TODO: Enable to specify intervals.

## Errors

### From `goddns`

`goddns` sends no text to STDOUT/STDERR. All info are stored in syslog.

```
Mar 10 20:03:01 some-host goddns[61735]: start
Mar 10 20:03:01 some-host goddns[61735]: loading /usr/local/etc/goddns.toml
Mar 10 20:03:01 some-host goddns[61735]: error occurred. trying again later: main.LoadConfig
                /Users/delphinus/.go/src/github.com/delphinus/goddns/config.go:29
          - Key: 'Configs.Domains[0].Hostname' Error:Field validation for 'Hostname' failed on the 'fqdn' tag
```

If error occurs, syslog adds log like above. When you fix config, `goddns` does again.

### From Google Domains

[Google Domains API][] returns some info in the body.

[Google Domains API]: https://support.google.com/domains/answer/6147083

| body            | IsSuccessful | isCritical |
|-----------------|--------------|------------|
| `good 1.2.3.4`  | **true**     | false      |
| `nochg 1.2.3.4` | **true**     | false      |
| `nohost`        | false        | **true**   |
| `badauth`       | false        | **true**   |
| `notfqdn`       | false        | **true**   |
| `badagent`      | false        | **true**   |
| `abuse`         | false        | **true**   |
| `911`           | false        | false      |
