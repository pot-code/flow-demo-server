package pk

type ID string

func ParseFromString(id string) ID {
	return ID(id)
}
