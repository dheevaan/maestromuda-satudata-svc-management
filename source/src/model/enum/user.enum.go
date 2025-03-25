package enum

type UserStatus int

const (
	UserStatus_active UserStatus = iota
	UserStatus_blocked
	UserStatus_inactive
)

func (index UserStatus) String() string {
	return []string{
		"active",
		"blocked",
		"inactive",
	}[index]
}
