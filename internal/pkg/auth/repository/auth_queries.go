package repository

const (
	GetUserByUsernameQuery = `
		SELECT id, username, password_hash
		FROM public.user
		WHERE username = $1;
	`

	CreateUserQuery = `
		INSERT INTO public.user (username, password_hash)
		VALUES ($1, $2)
		RETURNING id, username, password_hash;
	`
)
