package repo

const (
	createSchema     = `Create Schema %s`
	createUsersTable = `Create table %s.users (
					id UUID PRIMARY KEY,
					name text,
					email text UNIQUE,
					password text,
					role text,
					is_active boolean,
					created_at timestamp default current_timestamp,
    					updated_at timestamp default current_timestamp
				);`
	createApplication = `INSERT INTO applications (id, name, email, status, is_verified)
				VALUES ($1, $2, $3, $4, $5);
	 `
	getApplication    = `SELECT id,name,email,is_verified,status from applications`
	UpdateApplication = `Update applications SET name = $1 WHERE id = $2`
)
