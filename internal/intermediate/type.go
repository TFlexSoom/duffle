package intermediate

type TypeId uint32

const (
	TYPEID_NO_TYPE TypeId = iota
	TYPEID_BOOLEAN        // 1 bit (when possible)
	TYPEID_BYTE           // 1 byte
	TYPEID_CHAR           // 1-4 bytes
	TYPEID_INTEGER        // 1-4 bytes
	TYPEID_DECIMAL        // 2 or 4 bytes
	TYPEID_TEXT           // 0+ bytes
	TYPEID_LIST           // 0+ bytes For Configuration
	TYPEID_STRUCT         // 0+ bytes For Configuration
)
