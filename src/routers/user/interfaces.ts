export interface RegisterArgs {
	name: string
	email: string
	password: string
}

export interface LoginArgs {
	email: string
	password: string
}

export interface PasswordUpdateArgs {
	email: string
	oldPassword: string
	newPassword: string
}
