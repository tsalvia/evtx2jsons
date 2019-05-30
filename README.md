# evtxdump

Parse the EVTX file and output it in JSON format.

## Build

```
go get -u github.com/0xrawsec/golang-evtx/evtx
go build evtxdump.go
```

## Usage

```
$ evtxdump.exe -i Security.evtx
```

## Options

```
-d string
      This option is a short version of "--directory" option.
-directory string
      Specifies the destination directory for the converted files.
       (default "output")
-i string
      This option is a short version of "--input" option.
-ids string
      Specifies the event ID you want to output JOSN files.
      Use "," to separate multiple IDs.
      (default All Event IDs)
-input string
      This option is required.
      Specifies the EVTX file you want to convert to JSON file.
```

## Examples

1. Basic Usage
    ```
    $ evtxdump.exe -i Security.evtx
    ```

2. Specify the event IDs you want to output.
    ```
    $ evtxdump.exe -i Security.evtx -ids 4624,4625,1102
    ```

3. Specify the destination directory.
    ```
    $ evtxdump.exe -i Security.evtx -d output/jsons
    ```