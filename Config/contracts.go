package Config

// Config -
type Config interface {
	CreateOrReturnInstance() Config
	DefaultConfig() Config
	String() string
}
