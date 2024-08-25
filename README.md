# Checkiner

This is a checkin script for some useful web sites.

**You MUST change config file including the filename, email and password.**

**Please DON'T expose the config file to others, or your account may be stolen.**

```bash
# config/example
example@example.com
ThisIsAPassword
```

## Build

```bash
mkdir bin
# `-s` 表示从可执行文件中剥离符号信息
# `-w` 表示禁止编译器产生警告信息
go build -ldflags "-s -w" -o bin/checkiner src/checkiner.go src/main.go src/utils.go 
```

## [AutoStart](https://wiki.archlinuxcn.org/wiki/KDE#%E8%87%AA%E5%90%AF%E5%8A%A8)

```bash
# /home/tianen/.config/autostart/checkiner.desktop
YOUR_PROJECT_DIR = /home/tianen/go/src/Checkiner # NOTE YOU MUST CHANGE THIS DIR
[Desktop Entry]
Exec=$YOUR_PROJECT_DIR/bin/checkiner -w THY@CUTECLOUD -p $YOUR_PROJECT_DIR/conf/THY@$YOUR_PROJECT_DIR/conf/CUTECLOUD -i 60 -l $YOUR_PROJECT_DIR/log/checkiner.log
Icon=
Name=checkiner
Path=
Terminal=False
Type=Application
```

```bash
nohup ./bin/checkiner -w THT0@THY1@CUTECLOUD -p /Users/kyden/gitProj/Checkiner/conf/THY0@/Users/kyden/gitProj/Checkiner/conf/THY1@/Users/kyden/gitProj/Checkiner/conf/CUTECLOUD -i 60 -l /Users/kyden/gitProj/Checkiner/log/checkiner.log >> /Users/kyden/gitProj/Checkiner/log/checkiner.log &
```

## TODO 引入 config 配置

## Reference

- [mapstructure](github.com/mitchellh/mapstructure)

- [assert](github.com/stretchr/testify/assert)
