package apimachinery

type MockClientSetInterface interface {
	ClientSetInterface
	MockDo(output interface{}) ClientSetInterface
}

func NewMockApiMachinery() MockClientSetInterface {
	return &MockClientSet{
		ClientSet: NewMockClientSet(),
	}
}

type MockClientSet struct {
	*ClientSet
}

func (mc *MockClientSet) MockDo(output interface{}) ClientSetInterface {
	mc.ClientSet.Mock.SetMockData = true
	mc.ClientSet.Mock.MockData = output
	return mc.ClientSet
}
