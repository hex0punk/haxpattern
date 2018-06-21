# haxpattern
Go implementation of pattern_create and offset_pattern utilities for exploit development. The code is mostly based on the implementation of the same tools by mona.py, a plugin for Immunity Debugger 

## Usage
```
λ haxpattern.exe -h
Usage of haxpattern.exe:
  -c    Create pattern
  -e string
        Egg to hunt for
  -o    Find offset
  -s int
        Size of pattern (default 20280)
 ```
 
 ## Create a pattern 
 We can create a new pattern by using the `-c` flag and with an optional size argument using the `-s` flag
 
 ```
 λ haxpattern.exe -c -s 2000

=====================================================
[+] Creating pattern of 2000 bytes

=====================================================
Pattern displayed here...
```

## Find offset 
You can locate the offset byte location by using the `-o` flag and entering a required egg (in ascii or hex)

```
λ haxpattern.exe -o -e 39694438

=====================================================
[+] Creating pattern of 20280 bytes
[+] Looking for egg 9iD8 in pattern of 20280 bytes
[+] Egg pattern 8Di9 (9iD8 reversed) found in cyclic pattern at position 2606

=====================================================
```
