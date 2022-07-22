export interface RegisterArgs {
	name: string
	email: string
	password: string
}

export interface LoginArgs {
	email: string
	password: string
}

export interface UserUpdateArgs {
	name: string
	email: string
	group: string
}

export interface PasswordUpdateArgs {
	email: string
	oldPassword: string
	newPassword: string
}

export interface ScopeArgs {
	email: string
	scope: string
}
