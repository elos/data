package data

// Compatible checks whether one model has the same
// DBType as another model
func Compatible(this Model, that Model) bool {
	return this.DBType() == that.DBType()
}

func possibleLink(s *RelationshipMap, this Model, other Model) (bool, error) {
	thisKind := this.Kind()

	links, ok := (*s)[thisKind]

	if !ok {
		return false, ErrUndefinedKind
	}

	otherKind := other.Kind()

	_, linkPossible := links[otherKind]

	if !linkPossible {
		return false, ErrUndefinedLink
	}

	return true, nil
}

func (s *RelationshipMap) linkFor(this Model, other Model) (Link, error) {
	_, err := possibleLink(s, this, other)
	if err != nil {
		return *new(Link), err
	}

	return (*s)[this.Kind()][other.Kind()], nil
}

func (s *RelationshipMap) Link(this Model, that Model) error {
	if !Compatible(this, that) {
		return ErrIncompatibleModels
	}

	thisLink, err := s.linkFor(this, that)

	if err != nil {
		return err
	} else {
		if err = this.Link(that, thisLink); err != nil {
			return err
		}
	}

	thatLink, err := s.linkFor(that, this)

	if err != nil {
		return err
	} else {
		if err = that.Link(this, thatLink); err != nil {
			return err
		}
	}

	return nil
}

func (s *RelationshipMap) Unlink(this Model, that Model) error {
	if !Compatible(this, that) {
		return ErrIncompatibleModels
	}

	thisLink, err := s.linkFor(this, that)

	if err != nil {
		return err
	} else {
		if err = this.Unlink(that, thisLink); err != nil {
			return err
		}
	}

	thatLink, err := s.linkFor(that, this)

	if err != nil {
		return err
	} else {
		if err = that.Unlink(this, thatLink); err != nil {
			return err
		}
	}

	return nil
}
