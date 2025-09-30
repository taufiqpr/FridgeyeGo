package db

var (
	QueryGetUserByEmail = `
        SELECT id, name, email 
        FROM users 
        WHERE email=$1 AND deleted_at IS NULL`

	QueryGetUserByID = `
        SELECT id, name, email 
        FROM users 
        WHERE id=$1 AND deleted_at IS NULL`

	QueryUpdateUser = `
        UPDATE users 
        SET name=$1, updated_at=NOW() 
        WHERE id=$2`

	QuerySoftDeleteUser = `
        UPDATE users 
        SET deleted_at=NOW() 
        WHERE id=$1`

	QueryDeleteLoginHistoryByEmail = `
        DELETE FROM login_history 
        WHERE user_email=$1`
)
