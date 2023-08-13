package repo

const (
	createSchema     = `Create Schema %s`
	createUsersTable = `Create table %s.users (
					id text PRIMARY KEY,
					name text,
					email text UNIQUE,
					password text,
					token text,
					role text,
					is_active boolean,
					created_at timestamp default current_timestamp,
    					updated_at timestamp default current_timestamp
				);`
	createApplicationTable = `CREATE TABLE %s.applications (
					id text PRIMARY KEY,
					secret text,
					name text,
					email text UNIQUE,
					status text,
					is_verified boolean,
					secret_viewed boolean,
					created_at timestamp default current_timestamp,
					updated_at timestamp default current_timestamp
				);`

	registerApplication = `INSERT INTO applications (id, name, email, status, is_verified,secret)
				VALUES ($1, $2, $3, $4, $5,$6);
	 `
	getApplication    = `SELECT id,secret,name,email,is_verified,status,secret_viewed from applications`
	UpdateApplication = `Update applications SET name = $1,secret_viewed= true WHERE id = $2`
)
