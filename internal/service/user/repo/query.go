package repo

const (
	createUser string = `Insert into %s.users (id,name,email,is_active,created_at,updated_at,password) values ($1,$2,$3,$4,$5,$6,$7)`
	updateUser string = `Update %s.users SET name= $1 where id=$2`
	deleteUser string = `Delete From %s.users where id=$2`
	getUser    string = `Select id,name,email,password from %s.users`
)
