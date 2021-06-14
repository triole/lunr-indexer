# Lunr Indexer

<!--- mdtoc: toc begin -->

1.	[Synopsis](#synopsis)
2.	[Index JSON Layout](#index-json-layout)
3.	[Help & Usage](#help---usage)
4.	[Disclaimer](#disclaimer)<!--- mdtoc: toc end -->

## Synopsis

Lunr Indexer recursively finds all markdown files inside a directory and builds a [lunrjs](https://lunrjs.com/) search index from them. It is blazingly fast and can run a watcher that reindexes automatically on changes inside the folder.

## Index JSON Layout

To check the layout of the resulting index json have a look at the test data inside the repo. For example [here](triole/lunr-indexer/blob/master/testdata/set1/assert.json).

## Help & Usage

For help and more information consider running

```shell
# simple usage
$ lunr-indexer /path/to/md_files

# custom output file
$ lunr-indexer /path/to/md_files -o /path/to/output.json

# running the watcher
$ lunr-indexer /path/to/md_files -w

# help and more information
$ lunr-indexer -h
```

## Disclaimer

Warning. Use this software at your own risk. I may not be hold responsible for any data loss, starving your kittens or losing the bling bling powerpoint presentation you made to impress human resources with the efficiency of your employee's performance.
