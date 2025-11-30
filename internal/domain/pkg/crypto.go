package pkg

type SecretSauce interface {
	MakeSauce(ingredient string) ([]byte, error)
	SauceReferee(sauce []byte, ingredient string) error
}
