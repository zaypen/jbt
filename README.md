# jbt

![](https://travis-ci.org/zaypen/jbt.svg?branch=master)

A cli tool to manage your JetBrains products.

## Features

- List all products
- Check for updates
- Download updates

_currently only works on macOS_

## Installation

1. Check latest version of jbt from [GitHub Release](https://github.com/zaypen/jbt/releases)
1. Download right binary for your system, please check out darwin for macOS
1. Unzip and copy it to `/usr/local/bin` or anywhere you like but in PATH
1. Run `chmod +x /usr/local/bin/jbt` to make it executable
1. Try `jbt version`

## Usage

### List installations of all products

```bash
> jbt list                                                                                            
Code |Product          |Installed |Version
---- |-------          |--------- |-------
AC   |AppCode          |no        |
CL   |CLion            |yes       |2018.1.5
DG   |DataGrip         |no        |
IIU  |IntelliJ IDEA    |yes       |2018.1.5
IIC  |IntelliJ IDEA CE |no        |
PS   |PhpStorm         |no        |
PCP  |PyCharm          |yes       |2018.1.4
PCC  |PyCharm CE       |yes       |2018.1.4
RM   |RubyMine         |no        |
WS   |WebStorm         |yes       |2018.1.5
```

### Check updates of all installed products

```bash
> jbt check
Code |Product          |Installed |Version  |Update |Latest
---- |-------          |--------- |-------  |------ |------
AC   |AppCode          |no        |         |no     |
CL   |CLion            |yes       |2018.1.4 |yes    |2018.1.4
DG   |DataGrip         |no        |         |no     |
IIU  |IntelliJ IDEA    |yes       |2018.1.5 |no     |
IIC  |IntelliJ IDEA CE |no        |         |no     |
PS   |PhpStorm         |no        |         |no     |
PCP  |PyCharm          |yes       |2018.1.4 |no     |
PCC  |PyCharm CE       |yes       |2018.1.4 |no     |
RM   |RubyMine         |no        |         |no     |
WS   |WebStorm         |yes       |2018.1.5 |no     |
```

### Update(download) all updates, or a specified one
```bash
> jbt update # download all
> jbt update CL  # download CLion
```


## TODO

- [ ] Install a new product
- [ ] Linux support
