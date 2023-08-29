package repo

const (
	createApplicationIdentity = `INSERT INTO application_identities (id, application_id, secret, secret_viewed)
				VALUES ($1, $2, $3, $4);
	`
	getApplicationIdentity = `Select application_id,secret,secret_viewed
                                        from application_identities
                                where application_id = $1`
	deleteApplicationIdentity = `Delete from application_identities where application_id = $1`
	updateApplicationIdentity = `Update application_identities Set secret_viewed=true where application_id = $1`
)
