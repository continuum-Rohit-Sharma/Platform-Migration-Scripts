package pdh

const (
	//ErrorSuccess is used to check the return values from the API (ERROR_SUCCESS used in MSDN)
	ErrorSuccess = 0
	//FmtDouble Return data as a double precision floating point real.
	FmtDouble = 0x00000200
	//FmtLarge Return data as large.
	FmtLarge = 0x00000400
	//MoreData returns when there is more data
	MoreData = 0x800007D2
)
