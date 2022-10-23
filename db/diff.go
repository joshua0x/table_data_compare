package db

type TableLevelDiff struct {
	TableName string
	IdOnlyInSrc []interface{}
	IdOnlyInDst []interface{}
	RowGotDiff []*RowLevelDiff
}

type RowLevelDiff struct {
	PkId interface{}
	ColDiffs []*ColDiff
}

type ColDiff struct {
	ColName string
	SrcVal interface{}
	DstVal interface{}
}

func setUpdiff(diff *TableLevelDiff,splited *Splited) {
	diff.IdOnlyInSrc = splited.OnlyInSrcIdList
	diff.IdOnlyInDst = splited.OnlyInDstIdList
}