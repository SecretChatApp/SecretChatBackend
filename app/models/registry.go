package models

type Model struct {
	Model interface{}
}

func RegisterModel() []Model {
	return []Model{
		{Model: User{}},
		{Model: ChatRoom{}},
		{Model: Message{}},
	}
}
