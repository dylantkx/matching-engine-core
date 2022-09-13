package orderbook

import "github.com/google/btree"

type limitTree = btree.BTreeG[limitTreeNode]
