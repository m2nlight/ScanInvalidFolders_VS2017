# Scan invalid folders in visual studio 2017 install layout directory.

## build

```cmd
go build
```

## run

### scan and outout file

```cmd
C:\> scaninvalidfolders.exe -d E:\vs2017community -o 1.txt
2018/01/31 10:00:56 Loading E:\vs2017community\Catalog.json ...
2018/01/31 10:00:56 Loaded success.
2018/01/31 10:00:56 Parsing valid folder names...
2018/01/31 10:00:56 Has 5334 valid folder names.
2018/01/31 10:00:56 Comparing folders of E:\vs2017community...
2018/01/31 10:00:57 2 folders is invalid.
New folder
New folder (2)
2018/01/31 10:00:57 Writting output file: 1.txt ...
2018/01/31 10:00:57 1.txt is written.
2018/01/31 10:00:57 Completed.
C:\> type 1.txt
New folder
New folder (2)

C:\>
```

And delete invalid folder from `1.txt` command:

```cmd
C:\> for /f "tokens=*" %a in (1.txt) do rd /s /q "E:\vs2017community\%a"
```

### only show invalid folder names

```cmd
C:\> scaninvalidfolders.exe -d E:\vs2017community -q
New folder
New folder (2)

C:\>
```

Delete invalid folder at cmd:

```cmd
C:\> for /f "tokens=*" %a in ('scaninvalidfolders.exe -d E:\vs2017community -q') do rd /s /q "E:\vs2017community\%a"
```

### show detail

```cmd
C:\> scaninvalidfolders.exe -d E:\vs2017community -v
```

### other

```cmd
scaninvalidfolders.exe -version
scaninvalidfolders.exe -help
```
