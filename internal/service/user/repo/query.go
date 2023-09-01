package repo

const (
	createUser string = `Insert into %s.users (id,name,email,isActive,created_at,updated_at) values ($1,$2,$3,$4,$5,$6)`
	updateUser string = `Update %s.users SET name= $1 where id=$2`
	deleteUser string = `Delete From %s.users where id=$2`
	getUser    string = `Select id,name,email,password from %s.users`
)
