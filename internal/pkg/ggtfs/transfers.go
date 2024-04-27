package ggtfs

type Transfer struct {
	FromStopId      string
	ToStopId        string
	TransferType    uint
	MinTransferTime uint
	LineNumber      int
}
