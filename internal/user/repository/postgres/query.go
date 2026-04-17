package postgres

const (
	getUserByIDQuery = `
SELECT id,username,firstname,lastname,middlename,email,created_at,updated_at,is_active,locked_at
FROM users WHERE id=@id`

	countUsersQueryBase    = `SELECT count(id) FROM users`
	getPagedUsersQueryBase = `SELECT id, username, is_active, firstname, lastname, middlename FROM users`
	getPagedUsersQueryTail = `ORDER BY username LIMIT @page_size OFFSET @offset`

	searchFilter = `
		(username like '%' || @search || '%'
		or lower(firstname || ' ' || middlename || ' ' || lastname) like '%' || @search || '%')
	`
	createUserCommand = `
insert into users
(id,username,email,firstname,lastname,middlename,created_at,updated_at)
values (@id, @username, @email, @firstname, @lastname, @middlename, @created_at, @updated_at)`

	updateUserCommand = `
update users
SET
username=@username, 
email=@email,
updated_at=@updated_at,
firstname=@firstname,
lastname=@lastname,
middlename=@middlename
where id=@id`

	deleteUserCommand = `DELETE FROM users WHERE id=@id`
)
