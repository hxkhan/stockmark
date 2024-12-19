package model

func (p Permit) Deposit(amount int) error {
	user, err := p.exchange()
	if err != nil {
		return err
	}

	return user.Deposit(amount)
}
