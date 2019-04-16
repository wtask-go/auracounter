package counter

type repositoryMock struct {
	num, delta, max int
}

func NewRepositoryMock() Repository {
	return &repositoryMock{
		num: 0,
		delta: 1,
		max: 	
	}
}