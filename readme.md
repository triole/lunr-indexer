# Lunr Indexer

<!--- mdtoc: toc begin -->

1.	[Synopsis](#synopsis)
2.	[Index JSON Layout](#index-json-layout)
3.	[Help & Usage](#help---usage)<!--- mdtoc: toc end -->

## Synopsis

Lunr Indexer recursively finds all markdown files inside a directory and builds a [lunrjs](https://lunrjs.com/) search index from them. It is blazingly fast and can run a watcher that reindexes automatically on changes inside the folder.

## Index JSON Layout

To check the layout of the resulting index json have a look at the test data inside the repo. For example [here](triole/lunr-indexer/blob/master/testdata/set1/assert.json).

## Help & Usage

For help and more information consider running

```shell
lunr-indexer -h
```
