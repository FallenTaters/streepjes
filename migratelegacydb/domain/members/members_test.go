package members

type testRepo struct {
	_getAll       func() ([]Member, error)
	_get          func(id int) (Member, error)
	_updateMember func(member Member) error
	_addMember    func(member Member) error
	_deleteMember func(id int) error
}

func (tr testRepo) getAll() ([]Member, error) {
	return tr._getAll()
}

func (tr testRepo) get(id int) (Member, error) {
	return tr._get(id)
}

func (tr testRepo) updateMember(member Member) error {
	return tr._updateMember(member)
}

func (tr testRepo) addMember(member Member) error {
	return tr._addMember(member)
}

func (tr testRepo) deleteMember(id int) error {
	return tr._deleteMember(id)
}
