package Helpers

// Is -
func Is(err1 error, err2 error) bool {
	if err1 == err2 {
		return true
	}

	if err1 != nil && err2 != nil {
		return err1.Error() == err2.Error()
	}

	return false
}
