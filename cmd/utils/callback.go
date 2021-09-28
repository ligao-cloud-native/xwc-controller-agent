package utils

type Callback struct {
	Step int
	Second int
	CallbackUrl string
}

func NewCallback() Callback {
	return Callback{
		Step: 5,
		Second: 5,
		CallbackUrl: Env.CallbackUrl,
	}
}