package customerrors

import "errors"

var (
	ErrDBNotConnected     = errors.New("OOPS !!!! sorry to say but you found some issue while data base connection ")
	ErrDBNotPinnged       = errors.New("not able to connect the ping the connected database ")
	ErrServerNOtConnected = errors.New("something went wrong, not able to reach the defined server or some other app is running on this ")

	// user
	ErrUserIDNotValid      = errors.New("user id is not valid ")
	ErrEmptyUserId         = errors.New("user id can't be empty")
	ErrUserNotFound        = errors.New("user not found")
	ErrUserExits           = errors.New("user already exists with this username ")
	ErrUserEmailExsits     = errors.New("user already exists with this email")
	ErrUserPhoneNoExists   = errors.New("user already exists with this phone number")
	ErrSomethingWentWrong  = errors.New("OOPS !!! something went wrong, please try again sometime later")
	ErrUpdateUsernameExits = errors.New("can't update user already exists with this username ")
	ErrUpdateEmailExits    = errors.New("can't update user already exists with this Email ")
	ErrUpdatePhoneNoeExits = errors.New("can't update user already exists with this Phone number ")
	ErrCanNotUpdateUser    = errors.New("could not able to update user")

	// common
	ErrListParsing = errors.New("error while parsing the list of data ")
)
