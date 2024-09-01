package validator

func anyToPtrString(value any) *string {
	if value == nil {
		return nil
	}

	return value.(*string)
}

func anyToPtrBool(value any) *bool {
	if value == nil {
		return nil
	}

	return value.(*bool)
}
