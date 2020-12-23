package auth

// login constants
const ActionLogin = 0x10
const ActionLoginSucceeded = 0x01
const ActionLoginFailed = 0x02

// create account constants
const ActionCreateAccount = 0x20
const ActionCreateAccountSucceeded = 0x03
const ActionCreateAccountFailed = 0x09
const ActionCreateAccountTaken = 0x04
const ActionCreateAccountDisabled = 0x08

// password change constants
const ActionPasswordChange = 0x30
const ActionPasswordChangeSucceeded = 0x06
const ActionPasswordChangeFailed = 0x07

// request new password constants
const ActionRequestNewPassword = 0x05
