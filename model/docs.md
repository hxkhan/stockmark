# Model API Reference

```go
// Authorization 
func Register(fname, lname, email, password string) (Permit, error)
func Login(email, password string) (Permit, error)

// Permit
func (Permit) Logout() error
func (Permit) PresentUser() (UserPresentation, error)
func (Permit) Deposit(amount int) error

// Others
func Inquire(ticker string) (StockPresentationBasic, error)

// Errors
var ErrUserAlreadyExists = errors.New("user already exists")
var ErrUserNotFound = errors.New("user not found")
var ErrIncorrectPassword = errors.New("incorrect password")
var ErrNoPermitProvided = errors.New("no permit provided")
var ErrNotLoggedIn = errors.New("not logged in")
var ErrUnknownTicker = errors.New("unknown ticker")
var ErrInvalidArguments = errors.New("invalid arguments")
```