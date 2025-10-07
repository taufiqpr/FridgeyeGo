package db

var (
	QueryCheckEmailExists = `
        SELECT COUNT(*) 
        FROM users 
        WHERE email=$1`

	QueryRegister = `
        INSERT INTO users (name, email, password, created_at) 
        VALUES ($1, $2, $3, NOW()) 
        RETURNING id, created_at`

	QueryLogin = `
        SELECT id, name, email, password 
        FROM users 
        WHERE email=$1 AND deleted_at IS NULL`

	QueryInsertLoginHistory = `
        INSERT INTO login_history (user_email, ip_address, user_agent, timestamp)
        VALUES ($1, $2, $3, NOW())`

        QueryLoginHistory = `
	SELECT id, user_email, COALESCE(ip_address,''), 
        COALESCE(user_agent,''), timestamp
	FROM login_history
	WHERE user_email=$1
	ORDER BY timestamp DESC`

        QueryGetUserByEmail = `
	SELECT id, name, email 
	FROM users 
	WHERE email=$1 AND deleted_at IS NULL`
)
