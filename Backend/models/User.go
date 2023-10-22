package models

type User struct {
	ID        string  `json:"_id"`
	Rev       string  `json:"_rev,omitempty"`
	Type      string  `json:"type"` // Admin, User, ASCFI, CREA, ASOC, ACADEMIC, etc A list of several different types of users that could result in some discounts or benefits
	Name      string  `json:"name"`
	NIF       string  `json:"nif"`
	EMAIL     string  `json:"email"`
	PSWD_HASH string  `json:"password_hash"`
	Credits   float32 `json:"credits"`
}

// func QueryAllUsers(db *kivik.DB, ctx context.Context) ([]User, error) {
// 	rows := db.Query(ctx, os.Getenv("USER_DESIGN_DOC"), os.Getenv("ALL_USERS_VIEW"))
// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}
// 	users := []User{}
// 	for rows.Next() {
// 		usr := User{}
// 		if err := rows.ScanKey(&usr); err != nil {
// 			return nil, err
// 		}
// 		fmt.Println(usr.EMAIL)
// 		users = append(users, usr)
// 	}

// 	return users, nil
// }

// func QueryAdminUsers(db *kivik.DB, ctx context.Context) ([]User, error) {
// 	rows := db.Query(ctx, os.Getenv("USER_DESIGN_DOC"), os.Getenv("ADMIN_USERS_VIEW"))
// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}
// 	users := []User{}
// 	for rows.Next() {
// 		usr := User{}
// 		if err := rows.ScanKey(&usr); err != nil {
// 			return nil, err
// 		}
// 		fmt.Println(usr.EMAIL)
// 		users = append(users, usr)
// 	}

// 	return users, nil
// }

// func QueryBasicUsers(db *kivik.DB, ctx context.Context) ([]User, error) {
// 	rows := db.Query(ctx, os.Getenv("USER_DESIGN_DOC"), os.Getenv("BASIC_USERS_VIEW"))
// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}
// 	users := []User{}
// 	for rows.Next() {
// 		usr := User{}
// 		if err := rows.ScanKey(&usr); err != nil {
// 			return nil, err
// 		}
// 		fmt.Println(usr.EMAIL)
// 		users = append(users, usr)
// 	}

// 	return users, nil
// }

// func QueryUserByEmail(db *kivik.DB, ctx context.Context) ([]User, error) {
// 	rows := db.Query(ctx, os.Getenv("USER_DESIGN_DOC"), os.Getenv("ALL_USERS_VIEW"), kivik.Options{
// 		"startkey": json.RawMessage(`{"email": "vicente.rojo@lluch.es"}`),
// 	})
// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}
// 	users := []User{}
// 	for rows.Next() {
// 		usr := User{}
// 		if err := rows.ScanKey(&usr); err != nil {
// 			return nil, err
// 		}
// 		fmt.Println(usr.EMAIL)
// 		users = append(users, usr)
// 	}

// 	return users, nil
// }
