# torrmleaf
Clean Leaf files from torrented directories.  

When downloading torrents over already populated directories we are prone to end with leftover files that can become quite a bunch over time.This utility will list the files from a directory that aren't on a torrent file and prompt you to move them out.  
  
* Usage:  
Specifying the torrent file and his already downloaded directory:  
`torrmleaf <.torrent file> <downladed torrent dir>`  
You can also use the the torrent file only if you are in the same dir of the default torrent dir:  
`torrmleaf <.torrent file>`
* Example:  
```
# torrmleaf_freebsd-amd64 MAME\ 0.191\ Software\ List\ CHDs\ \(merged\).torrent
Reading 'MAME 0.191 Software List CHDs (merged).torrent'.
    Name: MAME 0.191 Software List CHDs (merged)
    Creation Date: 2017-10-25 11:28:18 +0200 CEST
    Created By: qBittorrent v3.3.16
    Has 7960 files.
    Example: 'psx/gol/game of life, the (usa).chd'
Checking directory 'MAME 0.191 Software List CHDs (merged)'
    Has 7960 files.
```
