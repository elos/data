package data

// Compatible checks whether one model has the same
// DBType as another model
func Compatible(this Model, that Model) bool {
	return this.DBType() == that.DBType()
}

func possibleLink(s *RelationshipMap, this Model, n LinkName) (bool, error) {
	thisKind := this.Kind()

	links, ok := (*s)[thisKind]

	if !ok {
		return false, NewUndefinedKindError(thisKind)
	}

	_, linkPossible := links[n]

	if !linkPossible {
		return false, ErrUndefinedLink
	}

	return true, nil
}

func (s *RelationshipMap) linkFor(this Model, n LinkName) (Link, error) {
	_, err := possibleLink(s, this, n)
	if err != nil {
		return *new(Link), err
	}

	return (*s)[this.Kind()][n], nil
}

func (s *RelationshipMap) Link(this Model, other Model, n LinkName) error {
	if !Compatible(this, other) {
		return ErrIncompatibleModels
	}

	thisLink, err := s.linkFor(this, n)

	if err != nil {
		return err
	} else {
		if err = this.Link(other, thisLink); err != nil {
			return err
		}
	}

	if thisLink.Inverse == "" {
		return nil
	}

	otherLink, err := s.linkFor(other, thisLink.Inverse)

	if err != nil {
		return err
	} else {
		if err = other.Link(this, otherLink); err != nil {
			return err
		}
	}

	return nil
}

func (s *RelationshipMap) Unlink(this Model, other Model, n LinkName) error {
	if !Compatible(this, other) {
		return ErrIncompatibleModels
	}

	thisLink, err := s.linkFor(this, n)

	if err != nil {
		return err
	} else {
		if err = this.Unlink(other, thisLink); err != nil {
			return err
		}
	}

	if thisLink.Inverse == "" {
		return nil
	}

	otherLink, err := s.linkFor(other, thisLink.Inverse)

	if err != nil {
		return err
	} else {
		if err = other.Unlink(this, otherLink); err != nil {
			return err
		}
	}

	return nil
}
