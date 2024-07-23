package core

// Blockchain 结构体表示区块链
// Blockchain struct represents the blockchain
type Blockchain struct {
	store     Storage
	headers   []*Header
	validator Validator
}

// NewBlockchain 创建一个新的区块链
// NewBlockchain creates a new blockchain
func NewBlockchain(genesis *Block) (*Blockchain, error) {
	bc := &Blockchain{
		headers: []*Header{},
		store:   NewMemoryStore(),
	}
	bc.validator = NewBlockValidator(bc)
	err := bc.addBlockWithoutValidation(genesis)

	return bc, err
}

// SetValidator 设置区块链的验证器
// SetValidator sets the validator for the blockchain
func (bc *Blockchain) SetValidator(v Validator) {
	bc.validator = v
}

// AddBlock 添加一个区块到区块链
// AddBlock adds a block to the blockchain
func (bc *Blockchain) AddBlock(b *Block) error {
	// validate
	if err := bc.validator.ValidateBlock(b); err != nil {
		return err
	}
	return bc.addBlockWithoutValidation(b)
}

// HasBlock 检查区块链是否包含某个高度的区块
// HasBlock checks if the blockchain contains a block of a certain height
func (bc *Blockchain) HasBlock(height uint32) bool {
	return height <= bc.Height()
}

// Height 返回区块链的高度
// Height returns the height of the blockchain
func (bc *Blockchain) Height() uint32 {
	return uint32(len(bc.headers) - 1)
}

// addBlockWithoutValidation 添加一个未验证的区块到区块链
// addBlockWithoutValidation adds a block to the blockchain without validation
func (bc *Blockchain) addBlockWithoutValidation(b *Block) error {
	bc.headers = append(bc.headers, b.Header)
	return bc.store.Put(b)
}
