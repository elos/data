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

func (s *RelationshipMap) linkType(this Model, other Model) (LinkKind, error) {
	_, err := possibleLink(s, this, other)
	if err != nil {
		return "", err
	}

	return (*s)[this.Kind()][other.Kind()], nil
}

func linkWith(lk LinkKind, this Model, that Model) error {
	switch lk {
	case MulLink:
		this.LinkMul(that)
	case OneLink:
		this.LinkOne(that)
	default:
		return ErrUndefinedLinkKind
	}

	return nil
}

func unlinkWith(ln LinkKind, this Model, that Model) error {
	switch ln {
	case MulLink:
		this.UnlinkMul(that)
	case OneLink:
		this.UnlinkOne(that)
	default:
		return ErrUndefinedLinkKind
	}

	return nil
}

func (s *RelationshipMap) Link(this Model, that Model) error {
	if !Compatible(this, that) {
		return ErrIncompatibleModels
	}

	thisLinkType, err := s.linkType(this, that)

	if err != nil {
		return err
	} else {
		if err = linkWith(thisLinkType, this, that); err != nil {
			return err
		}
	}

	thatLinkType, err := s.linkType(that, this)

	if err != nil {
		return err
	} else {
		if err = linkWith(thatLinkType, that, this); err != nil {
			return err
		}
	}

	return nil
}

func (s *RelationshipMap) Unlink(this Model, that Model) error {
	if !Compatible(this, that) {
		return ErrIncompatibleModels
	}

	thisLinkType, err := s.linkType(this, that)

	if err != nil {
		return err
	} else {
		if err = unlinkWith(thisLinkType, this, that); err != nil {
			return err
		}
	}

	thatLinkType, err := s.linkType(that, this)

	if err != nil {
		return err
	} else {
		if err = unlinkWith(thatLinkType, that, this); err != nil {
			return err
		}
	}

	return nil
}
