# Model API Reference

```go
// Authorization 
func Register(fname, lname, email, password string) (Permit, error)
func Login(email, password string) (Permit, error)

// Permit
func (Permit) Logout() error
func (Permit) RenderPortfolio() (PortfolioPage, error)
func (Permit) RenderTrade() (TradePage, error)
func (Permit) Deposit(amount int) error

// Others
func Inquire(ticker string) (StockPresentationBasic, error)

// Errors
var ErrAccountExists = errors.New("account exists")
var ErrAccountNotExists = errors.New("account not exists")
var ErrIncorrectPassword = errors.New("incorrect password")
var ErrNoPermitProvided = errors.New("no permit provided")
var ErrNotLoggedIn = errors.New("not logged in")
var ErrUnknownTicker = errors.New("unknown ticker")
var ErrInvalidArguments = errors.New("invalid arguments")
```